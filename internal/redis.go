package internal

import (
	"time"

	"github.com/go-redis/redis"
)

func SetupRedis(address string) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: address,

		ReadTimeout:  time.Second * 5,
		WriteTimeout: time.Second * 2,
		DialTimeout:  time.Second * 5,
	})
}
