package handlers

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/haydnmeyburgh/cod-eSports-tournament-manager/internal/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MatchResultHandler struct {}

// Handles the creation of a new match result.
func (h *MatchResultHandler) CreateMatchResult(c *gin.Context) {
	userIDStr, err := models.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	var newMatchResult models.MatchResult

	if err := c.ShouldBindJSON(&newMatchResult); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}

	newMatchResult.OrganizerID = userID

	createdMatchResult, err := models.CreateMatchResult(c, &newMatchResult)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, createdMatchResult)
}

// GetMatchResultHandler handles the retrieval of a match result by its ID.
func (h *MatchResultHandler) GetMatchResultHandler(c *gin.Context) {
	matchResultID := c.Param("id")

	id, err := primitive.ObjectIDFromHex(matchResultID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Match Result ID format"})
		return
	}

	matchResult, err := models.GetMatchResultByID(c, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, matchResult)
}