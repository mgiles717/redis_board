package leaderboard

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
)

func InitRedis() *redis.Client {

	RedisClient := redis.NewClient(&redis.Options{
		Addr: "localhost:6379", // Default Redis Addr + Port
	})

	_, err := RedisClient.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("Failed to connect to the Redis instance: %v", err)
	}
	log.Println("Connected to the Redis instance!")

	return RedisClient
}

func RedisHello(rdb *redis.Client) {
	log.Println("Hello!")
}

func SetUserScore(rdb *redis.Client, username string, score int64) error {
	ctx := context.Background()

	return rdb.ZAdd(ctx, "leaderboard", redis.Z{
		Member: username,
		// ZAdd only accepts f64
		Score: float64(score),
	}).Err()
}

// func GetLeaderboard(rdb *redis.Client) error {

// }
