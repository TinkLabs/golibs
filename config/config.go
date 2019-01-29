package config

import (
	"encoding/json"
	"fmt"

	"github.com/tinklabs/golibs/cmd"
	"github.com/tinklabs/golibs/consul"
)

var Data map[string]interface{}

func Init() {
	cf := cmd.GetCmdFlag()
	cc := consul.GetConsulClient()

	// Lookup the pair
	configPath := fmt.Sprintf("b2c/%s/config", cf.ServerName)
	pair, _, err := cc.KV.Get(configPath, nil)
	if err != nil {
		panic(fmt.Sprintf("get configuration from consul:%v", err))
	}

	err = json.Unmarshal(pair.Value, &Data)
	if err != nil {
		panic(fmt.Sprintf("config is not json:%v", err))
	}
}

func TakeDbUrl() string {
	if v, isExist := Data["db_url"]; isExist {
		if v, ok := v.(string); ok {
			return v
		} else {
			panic("db url is not string")
		}
	} else {
		panic("db url not exist")
	}
}
