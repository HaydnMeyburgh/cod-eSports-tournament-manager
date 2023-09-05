package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/haydnmeyburgh/cod-eSports-tournament-manager/internal/handlers"
)

// Setup user routes
func SetupTeamRoutes(r *gin.Engine, teamHandler *handlers.TeamHandler) {
	teamRoutes := r.Group("/teams")
	{
		teamRoutes.POST("/", teamHandler.CreateTeam)
		teamRoutes.GET("/:id", teamHandler.GetTeamByID)
		teamRoutes.PUT("/:id", teamHandler.UpdateTeam)
		teamRoutes.DELETE("/:id", teamHandler.DeleteTeam)
	}
}