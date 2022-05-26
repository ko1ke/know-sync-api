package redis_db

import (
	"github.com/go-redis/redis/v7"
)

var (
	Client *redis.Client
)

func init() {
	Client = redis.NewClient(&redis.Options{
		Addr: "redis:6379",
		DB:   0})
	_, err := Client.Ping().Result()
	if err != nil {
		panic(err)
	}
}
