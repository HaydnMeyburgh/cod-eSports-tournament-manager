package routes

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/haydnmeyburgh/cod-eSports-tournament-manager/internal/auth"
	"github.com/haydnmeyburgh/cod-eSports-tournament-manager/internal/handlers"
)

// Setup tournament routes
func SetupTournamentRoutes(r *gin.Engine, tournamentHandler *handlers.TournamentHandler) {
	jwtSecret := os.Getenv("SECRET_KEY")
	if jwtSecret == "" {
		log.Fatalf("SECRET_KEY environment variable is not set")
	}

	tournamentRoutes := r.Group("/tournament")
	{
		tournamentRoutes.GET("/:id", tournamentHandler.GetTournamentByID)

		tournamentRoutes.Use(auth.AuthMiddleware(jwtSecret))
		
		tournamentRoutes.POST("/", tournamentHandler.CreateTournament)
		tournamentRoutes.PUT("/:id", tournamentHandler.UpdateTournament)
		tournamentRoutes.DELETE("/:id", tournamentHandler.DeleteTournament)
	}
}