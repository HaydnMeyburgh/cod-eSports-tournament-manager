package routes

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/haydnmeyburgh/cod-eSports-tournament-manager/internal/handlers"
	"github.com/haydnmeyburgh/cod-eSports-tournament-manager/internal/auth"
)

// Setup user routes
func SetupMatchResultRoutes(r *gin.Engine, matchResultHandler *handlers.MatchResultHandler) {
	jwtSecret := os.Getenv("SECRET_KEY")
	if jwtSecret == "" {
		log.Fatalf("SECRET_KEY environment variable is not set")
	}

	matchResultRoutes := r.Group("/match-results")
	{
		matchResultRoutes.GET("/", matchResultHandler.GetMatchResultsByOrganiserID)
		matchResultRoutes.GET("/:id", matchResultHandler.GetMatchResultById)

		matchResultRoutes.Use(auth.AuthMiddleware(jwtSecret))

		matchResultRoutes.POST("/", matchResultHandler.CreateMatchResult)
		matchResultRoutes.PUT("/:id", matchResultHandler.UpdateMatchResult)
		matchResultRoutes.DELETE("/:id", matchResultHandler.DeleteMatchResult)
	}
}