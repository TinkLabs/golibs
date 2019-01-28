package client

import (
	"fmt"
	"sync"
	"time"

	api "github.com/hashicorp/consul/api"

	"github.com/tinklabs/golibs/consul"
	terr "github.com/tinklabs/golibs/error"
	"github.com/tinklabs/golibs/log"
)

var (
	services sync.Map
)

func Init() {
	services = getServicesFromAgent()
	go startSync()
}

func startSync() {
	ticker := time.NewTicker(60 * time.Second)
	for range ticker.C {
		services = getServicesFromAgent()
		fmt.Printf("services = %+v\n", services)
	}
}

func getServicesFromAgent() (sm sync.Map) {
	var err error
	var m map[string]*api.AgentService

	cc := consul.GetConsulClient()
	if m, err = cc.Health.Service(); err != nil {
		log.Warn(terr.ErrConsul.AddExtra(err.Error()))
		return
	}

	for _, sa := range m {
		if l, loaded := sm.LoadOrStore(sa.Service, &[]*api.AgentService{sa}); loaded {
			l := append(*(l.(*[]*api.AgentService)), sa)
			sm.Store(sa.Service, &l)
		}
	}

	return sm
}

func GetService(name string) (string, bool) {
	if v, ok := services.Load(name); ok {
		v := v.(*[]*api.AgentService) // TODO random choose
		s := (*v)[0]
		return fmt.Sprintf("%s:%d", s.Address, s.Port), true
	}

	return "", false
}
