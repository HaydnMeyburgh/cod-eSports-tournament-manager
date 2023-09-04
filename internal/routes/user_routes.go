package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/haydnmeyburgh/cod-eSports-tournament-manager/internal/handlers"
)

// Setup user routes
func SetupUserRoutes(r *gin.Engine, userHandler *handlers.UserHandler) {
	userRoutes := r.Group("/users")
	{
		userRoutes.POST("/register", userHandler.RegisterUser)
		userRoutes.POST("/login", userHandler.LoginUser)
		userRoutes.POST("/logout", userHandler.LogoutUser)
		userRoutes.PUT("/update", userHandler.UpdateUser)
	}
}