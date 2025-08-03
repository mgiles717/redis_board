package leaderboard

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
)

// Creates a Redis Client that binds to localhost:6379
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

// Retrieves the whole ZSET with the key "leaderboard"
func GetWholeLeaderboard(rdb *redis.Client) []redis.Z {
	ctx := context.Background()

	membersWithScores, err := rdb.ZRangeWithScores(ctx, "leaderboard", 0, -1).Result()
	if err != nil {
		panic(err)
	}

	return membersWithScores
}

// Sets a username and score in the Redis ZSET with the key "leaderboard"
func SetUserScore(rdb *redis.Client, username string, score int64) error {
	ctx := context.Background()

	return rdb.ZAdd(ctx, "leaderboard", redis.Z{
		Member: username,
		// ZAdd only accepts f64
		Score: float64(score),
	}).Err()
}
