package cache

import (
	"fmt"

	"github.com/go-redis/redis"
	"github.com/tinklabs/golibs/config"
)

var (
	Client *redis.Client
)

func Init() {
	addr, pw := config.TakeCacheAddressAndPassword()
	Client = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: pw, // no password set
		DB:       0,  // use default DB
	})

	_, err := Client.Ping().Result()
	if err != nil {
		panic(fmt.Sprintf("init cache:%v", err))
	}
}
