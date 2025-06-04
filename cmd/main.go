package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type falseData struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
}

var testData = []falseData{
	{Id: 123, Username: "Bob"},
	{Id: 456, Username: "Alex"},
	{Id: 789, Username: "John"},
}

// GET '/' Route
func homeResponse(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, testData)
}

// GET '/hello' Route
func helloResponse(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, "Hello")
}

func main() {
	router := gin.Default()

	// Define Routes
	router.GET("/", homeResponse)
	router.GET("/hello", helloResponse)

	router.Run()
}
