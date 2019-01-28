package client

import (
	"fmt"
	"testing"

	"github.com/tinklabs/golibs/cmd"
	"github.com/tinklabs/golibs/consul"
	"github.com/tinklabs/golibs/log"
)

func init() {
	cmd.Init()
	log.Init()
	consul.Init()
	Init()
}

func TestGetService(t *testing.T) {
	s, ok := GetService("b2c-gateway")
	fmt.Printf("s = %+v\n", s)
	fmt.Printf("ok = %+v\n", ok)
}
