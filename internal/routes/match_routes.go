package routes

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/haydnmeyburgh/cod-eSports-tournament-manager/internal/auth"
	"github.com/haydnmeyburgh/cod-eSports-tournament-manager/internal/handlers"
)

// Setup tournament routes
func SetupMatchRoutes(r *gin.Engine, matchHandler *handlers.MatchHandler) {
	jwtSecret := os.Getenv("SECRET_KEY")
	if jwtSecret == "" {
		log.Fatalf("SECRET_KEY environment variable is not set")
	}

	matchRoutes := r.Group("/match")
	{
		matchRoutes.GET("/:id", matchHandler.GetMatchByID)

		matchRoutes.Use(auth.AuthMiddleware(jwtSecret))
		
		matchRoutes.POST("/", matchHandler.CreateMatch)
		matchRoutes.PUT("/:id", matchHandler.UpdateMatch)
		matchRoutes.DELETE("/:id", matchHandler.DeleteMatch)
	}
}
