package helpers

import (
	"github.com/go-redis/redis/v8"
	"os"
)

func ConnectToRedis() (client *redis.Client) {
	return redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_HOST"),
		Password: os.Getenv("REDIS_PASS"),
		DB:       0,
	})
}
