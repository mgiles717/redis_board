package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"

	"github.com/mgiles717/redis_board/pkg/leaderboard"
)

// Home page response for '/' route (index)
func HomeResponse(c *gin.Context) {
	c.HTML(http.StatusOK, "index.tmpl", gin.H{
		"title": "Homepage",
	})
}

// SendLeaderboardResponse returns a Gin Handler Function,
// fetching all leaderboard data from the Redis ZSET with the key "leaderboard".
// This data is then converted to an array of UserRecords, which is a struct consisting
// of the user name and score.
//
// Finally, this is passed into the template.
func SendLeaderboardResponse(rdb *redis.Client) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		membersWithScores := leaderboard.GetWholeLeaderboard(rdb)

		type UserRecord struct {
			Member string
			Score  uint32
		}

		items := make([]UserRecord, len(membersWithScores))
		for i, z := range membersWithScores {
			items[i] = UserRecord{
				Member: z.Member.(string),
				Score:  uint32(z.Score),
			}
		}

		ctx.HTML(http.StatusOK, "leaderboard.tmpl", gin.H{
			"LeaderboardItems": items,
			"Title":            "Leaderboard",
		})
	}
}

// SendScoreResponse returns a Gin Handler Function
// that processes a JSON post request and the username parameter
// set in the route e.g (/users/abcd) to set the username and score
// in the Redis ZSET with the key "leaderboard"
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
	// router.LoadHTMLGlob("templates/*")
	// LoadHTMLGlob seems to load in an unnamed 3rd template??
	router.LoadHTMLFiles("templates/index.tmpl", "templates/leaderboard.tmpl")

	// Define Routes
	// GET
	router.GET("/", HomeResponse)
	router.GET("/leaderboard", SendLeaderboardResponse(redisClient))
	// PUT
	router.PUT("/users/:username", SendScoreResponse(redisClient))

	router.Run()
}
