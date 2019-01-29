package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/tinklabs/golibs/utils"
)

type CmdFlag struct {
	Debug             bool
	Env               string
	ServerName        string
	ServerAddress     string
	ServerPort        int
	ConsulAddress     string
	ConsulPort        string
	ConsulAccessToken string
}

var cmdFlag *CmdFlag

func Init() {
	var port int
	var debug, dontCheck bool

	serverName := GetEnvPanic("SERVER_NAME")
	profileEnv := GetEnvWithDefault("PROFILE_ENV", "dev")
	consulAddress := GetEnvWithDefault("CONSUL_ADDRESS", "http://127.0.0.1")
	consulPort := GetEnvWithDefault("CONSUL_PORT", "8500")
	consulAccessToken := GetEnvWithDefault("CONSUL_ACCESS_TOKEN", "")

	if GetEnvWithDefault("DONT_CHECK_ETH_NAME", "false") == "false" {
		dontCheck = false
	} else {
		dontCheck = true
	}

	serverAddress := GetEnvWithDefault("SERVER_ADDRESS", utils.GetIntranetIp(dontCheck))

	if GetEnvWithDefault("DEBUG", "true") == "false" {
		debug = false
	} else {
		debug = true
	}

	port, err := strconv.Atoi(GetEnvWithDefault("SERVER_PORT", "8080"))
	if err != nil {
		panic(fmt.Sprintf("server port:%v", err))
	}

	if GetEnvWithDefault("RANDOM_PORT", "false") == "true" {
		port = utils.GetPort()
	}

	if profileEnv == "production" {
		consulAddress = utils.GetConsulAddressFromMetadata()
	}

	cmdFlag = &CmdFlag{
		Debug:             debug,
		ServerName:        serverName,
		ServerAddress:     serverAddress,
		ServerPort:        port,
		ConsulAddress:     consulAddress,
		ConsulPort:        consulPort,
		ConsulAccessToken: consulAccessToken,
	}
}

func GetCmdFlag() *CmdFlag {
	return cmdFlag
}

func IsDebug() bool {
	return cmdFlag.Debug
}

func GetEnvWithDefault(env, option string) string {
	rv := os.Getenv(env)
	if len(rv) < 1 {
		return option
	}

	return rv
}

func GetEnvPanic(env string) string {
	rv := os.Getenv(env)
	if len(rv) < 1 {
		panic(fmt.Sprintf("%s not provided", env))
	}

	return rv
}
