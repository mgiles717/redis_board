package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"

	"github.com/mgiles717/redis_board/pkg/redishandler"
)

// GET '/' Route
func homeResponse(c *gin.Context) {
	c.HTML(http.StatusOK, "index.tmpl", gin.H{
		"title": "Homepage",
	})
}

// GET '/hello' Route
//
// Handler function to allow access to redis client while satisfying
// Gin handler signature.
func leaderboardResponse(rdb *redis.Client) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		redishandler.RedisHello(rdb)
	}

}

func main() {
	redisClient := redishandler.InitRedis()
	router := gin.Default()

	// Allow access to read template directory outside of cmd
	router.LoadHTMLGlob("templates/*")

	// Define Routes
	router.GET("/", homeResponse)
	router.GET("/leaderboard", leaderboardResponse(redisClient))

	router.Run()
}
