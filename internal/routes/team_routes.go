package routes

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/haydnmeyburgh/cod-eSports-tournament-manager/internal/handlers"
	"github.com/haydnmeyburgh/cod-eSports-tournament-manager/internal/auth"
)

// Setup user routes
func SetupTeamRoutes(r *gin.Engine, teamHandler *handlers.TeamHandler) {
	jwtSecret := os.Getenv("SECRET_KEY")
	if jwtSecret == "" {
		log.Fatalf("SECRET_KEY environment variable is not set")
	}

	teamRoutes := r.Group("/teams")
	{
		teamRoutes.GET("/", teamHandler.GetTeamsByOrganiserID)
		teamRoutes.GET("/:id", teamHandler.GetTeamByID)

		teamRoutes.Use(auth.AuthMiddleware(jwtSecret))

		teamRoutes.POST("/", teamHandler.CreateTeam)
		teamRoutes.PUT("/:id", teamHandler.UpdateTeam)
		teamRoutes.DELETE("/:id", teamHandler.DeleteTeam)
	}
}