// Exectuable package
package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/haydnmeyburgh/cod-eSports-tournament-manager/internal/database"
	"github.com/haydnmeyburgh/cod-eSports-tournament-manager/internal/handlers"
	"github.com/haydnmeyburgh/cod-eSports-tournament-manager/internal/routes"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from the .env file
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Connect to MongoDB
	if err := database.ConnectToMongoDB(); err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer database.GetMongoClient().Disconnect(context.TODO())

	// Get server port from env variable or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port
	}

	// Gin router creation
	router := gin.Default()

	// SImple root route
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

	// Initialise handlers
	userHandler := handlers.NewUserHandler()

	// Setup routes
	routes.SetupUserRoutes(router, userHandler)

	// Start server, or log error if problem with server starting
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
