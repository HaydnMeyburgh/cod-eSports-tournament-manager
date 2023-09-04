package routes

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/haydnmeyburgh/cod-eSports-tournament-manager/internal/auth"
	"github.com/haydnmeyburgh/cod-eSports-tournament-manager/internal/handlers"
)

// Setup user routes
func SetupUserRoutes(r *gin.Engine, userHandler *handlers.UserHandler) {
	jwtSecret := os.Getenv("SECRET_KEY")
	if jwtSecret == "" {
		log.Fatalf("SECRET_KEY environment variable is not set")
	}

	userRoutes := r.Group("/users")
	{
		userRoutes.POST("/register", userHandler.RegisterUser)
		userRoutes.POST("/login", userHandler.LoginUser)
		
		userRoutes.Use(auth.AuthMiddleware(jwtSecret))
		
		userRoutes.POST("/logout", userHandler.LogoutUser)
		userRoutes.PUT("/update", userHandler.UpdateUser)
	}
}