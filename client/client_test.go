package client

import (
	"github.com/tinklabs/golibs/cmd"
	"github.com/tinklabs/golibs/consul"
	"github.com/tinklabs/golibs/log"
)

func init() {
	cmd.Init()
	log.Init()
	consul.Init()
}
