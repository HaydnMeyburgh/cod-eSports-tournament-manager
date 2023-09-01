package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port
	}

	router := gin.Default()
	router.GET("/", func(ctx *gin.Context) {
		ctx.IndentedJSON(http.StatusOK, gin.H{
			"message": "Hello there, server is running",
		})
	})

	// error handling for routes that aren't defined
	router.NoRoute(func(ctx *gin.Context) {
		ctx.IndentedJSON(http.StatusNotFound, gin.H{
			"error": "Route not found",
		})
	})

	// Log error if server can't start
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
