package redis

import (
	"context"

	"log"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client

func InitRedis() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr: "localhost:6379", // Default Redis Addr + Port
	})

	_, err := RedisClient.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("Failed to connect to the Redis instance: %v", err)
	}
	log.Println("Connected to the Redis instance!")
}
