package helpers

import "github.com/go-redis/redis/v8"

func ConnectToRedis() (client *redis.Client) {
	return redis.NewClient(&redis.Options{
		Addr:     "localhost:30001",
		Password: "",
		DB:       0,
	})
}
