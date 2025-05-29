package config

import (
	"context"
	"os"
	"sync"

	"github.com/redis/go-redis/v9"
)

var (
	Client *redis.Client
	once   sync.Once
)

func InitRedis() *redis.Client {
	once.Do(func() {
		redisAddr := os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT")
		Client = redis.NewClient(&redis.Options{
			Addr:     redisAddr,
			Password: os.Getenv("REDIS_PASSWORD"),
			DB:       0,
		})

		if err := Client.Ping(context.Background()).Err(); err != nil {
			panic("Failed to connect to Redis: " + err.Error())
		}
	})
	return Client
}
