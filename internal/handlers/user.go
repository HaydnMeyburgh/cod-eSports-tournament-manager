package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/haydnmeyburgh/cod-eSports-tournament-manager/internal/models"
)

type UserHandler struct{}

// Handles user registration
func (h *UserHandler) RegisterUser(c *gin.Context) {
	var newUser models.User

	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := models.RegisterUser(c, &newUser); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User registered Successfully"})
}

// Handles user login
func (h *UserHandler) LoginUser(c *gin.Context) {
	var loginUser struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"Password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&loginUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := models.LoginUser(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

// Handles user logout
func (h *UserHandler) LogoutUser(c *gin.Context) {

}

// Handles user profile updates
func (h *UserHandler) UpdateUser(c *gin.Context) {

}

func NewUserHandler() *UserHandler {
	return &UserHandler{}
}
