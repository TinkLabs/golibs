package consul

import (
	"fmt"
	"time"

	consul "github.com/hashicorp/consul/api"

	"github.com/tinklabs/golibs/cmd"
	"github.com/tinklabs/golibs/log"
	"github.com/tinklabs/golibs/utils"
)

var cc *ConsulClient

type ConsulClient struct {
	ServerID      string
	ServerName    string
	ServerAddress string
	ServerPort    int
	TTL           time.Duration
	Agent         *consul.Agent
	Health        *consul.Health
	KV            *consul.KV
}

func GetConsulClient() *ConsulClient {
	return cc
}

func Init() {
	cf := cmd.GetCmdFlag()

	dc := consul.DefaultConfig()
	dc.Token = cf.ConsulAccessToken
	dc.Address = fmt.Sprintf("%s:%s", cf.ConsulAddress, cf.ConsulPort)
	c, err := consul.NewClient(dc)
	if err != nil {
		panic(fmt.Sprintf("new consul client:%v", err))
	}

	agent := c.Agent()
	health := c.Health()
	kv := c.KV()

	uuid, err := utils.UUID()
	if err != nil {
		panic(fmt.Sprintf("generate uuid:%v", err))
	}

	id := fmt.Sprintf("%s-%s", cf.ServerName, uuid)

	cc = &ConsulClient{
		ServerID:      id,
		ServerName:    cf.ServerName,
		ServerAddress: cf.ServerAddress,
		ServerPort:    cf.ServerPort,
		TTL:           time.Second * 30,
		Agent:         agent,
		Health:        health,
		KV:            kv,
	}
}

func (c *ConsulClient) Register() {
	def := &consul.AgentServiceRegistration{
		ID:      c.ServerID,
		Name:    c.ServerName,
		Address: c.ServerAddress,
		Port:    c.ServerPort,
		Check: &consul.AgentServiceCheck{
			DeregisterCriticalServiceAfter: "60m",
			TTL:                            c.TTL.String(),
		},
	}

	if err := cc.Agent.ServiceRegister(def); err != nil {
		panic(fmt.Sprintf("register:%v", err))
	}
	log.Info("Register service:" + c.ServerID)

	go c.updateTTL()
}

func (c *ConsulClient) updateTTL() {
	ticker := time.NewTicker(c.TTL / 2)
	for range ticker.C {
		if err := c.Agent.UpdateTTL("service:"+c.ServerID, "I'm alive", "pass"); err != nil {
			log.Error(err)
		}
	}
}

func (c *ConsulClient) Deregister() error {
	log.Info("Deregister service:" + c.ServerID)
	return c.Agent.ServiceDeregister(c.ServerID)
}
