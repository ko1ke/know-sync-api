package redis_db

import (
	"net/url"
	"os"

	"github.com/go-redis/redis/v7"
)

var (
	Client *redis.Client
)

func Start() {
	addr, password := getRedisAddr()
	Client = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password})

	_, err := Client.Ping().Result()
	if err != nil {
		panic(err)
	}
}

func getRedisAddr() (addr string, password string) {
	redisURL := os.Getenv("REDIS_URL")

	redisInfo, err := url.Parse(redisURL)
	if err != nil {
		return
	}

	if redisInfo.Host != "" {
		// production
		addr = redisInfo.Host
	} else {
		// local dev
		addr = redisURL
	}

	if redisInfo.User != nil {
		password, _ = redisInfo.User.Password()
	}
	return
}
