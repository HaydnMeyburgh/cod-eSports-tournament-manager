package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/haydnmeyburgh/cod-eSports-tournament-manager/internal/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TeamHandler struct{}

func (h *TeamHandler) CreateTeam(c *gin.Context) {
	var newTeam models.Team

	if err := c.ShouldBindJSON(&newTeam); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	createdTeam, err := models.CreateTeam(c, &newTeam)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, createdTeam)
}

