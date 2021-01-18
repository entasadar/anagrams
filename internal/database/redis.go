package database

import (
	"context"
	"github.com/go-redis/redis/v8"
)

var (
	Ctx         = context.Background()
	RedisClient *redis.Client
)

func InitRedis(redisHost string) error {
	RedisClient = redis.NewClient(&redis.Options{
		Addr: redisHost,
	})

	if _, err := RedisClient.Ping(Ctx).Result(); err != nil {
		return err
	}
	return nil
}
