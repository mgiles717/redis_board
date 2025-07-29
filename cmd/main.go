package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"

	"github.com/mgiles717/redis_board/pkg/leaderboard"
)

// GET '/' Route
func HomeResponse(c *gin.Context) {
	c.HTML(http.StatusOK, "index.tmpl", gin.H{
		"title": "Homepage",
	})
}

// GET '/hello' Route
//
// Handler function to allow access to redis client while satisfying
// Gin handler signature.
// func LeaderboardResponse(rdb *redis.Client) gin.HandlerFunc {
// 	return func(ctx *gin.Context) {
// 		leaderboard.GetLeaderboard(rdb)
// 	}

// }

func SendScoreResponse(rdb *redis.Client) gin.HandlerFunc {
	type UserPayload struct {
		Score int64 `json:"score"`
	}

	return func(ctx *gin.Context) {
		var payload UserPayload
		var username string = ctx.Param("username")
		// Store in Redis
		if err := ctx.ShouldBindJSON(&payload); err != nil {
			ctx.JSON(400, gin.H{"error": "Invalid Request Body"})
			return
		}
		err := leaderboard.SetUserScore(rdb, username, payload.Score)

		// Response
		if err != nil {
			ctx.JSON(503, gin.H{"status": err})
			return
		}
		ctx.JSON(200, gin.H{"status": "OK"})
	}
}

func main() {
	redisClient := leaderboard.InitRedis()
	router := gin.Default()

	// Allow access to read template directory outside of cmd
	router.LoadHTMLGlob("templates/*")

	// Define Routes
	router.GET("/", HomeResponse)
	router.PUT("/users/:username", SendScoreResponse(redisClient))
	// router.GET("/leaderboard", LeaderboardResponse(redisClient))

	router.Run()
}
