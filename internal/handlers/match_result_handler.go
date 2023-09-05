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

// Handles the retrieval of a match result by its ID.
func (h *MatchResultHandler) GetMatchResultById(c *gin.Context) {
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

// Handles the updating of an existing match result.
func (h *MatchResultHandler) UpdateMatchResult(c *gin.Context) {
	matchResultID := c.Param("id")

	id, err := primitive.ObjectIDFromHex(matchResultID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Match Result ID format"})
		return
	}

	var updatedMatchResult models.MatchResult

	if err := c.ShouldBindJSON(&updatedMatchResult); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userIDStr, err := models.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userID, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}

	if userID != updatedMatchResult.OrganizerID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not the organizer of this team"})
		return
	}

	err = models.UpdateMatchResult(c, id, &updatedMatchResult)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Match Result updated successfully"})
}

// Handles the deletion of a match result by its ID.
func (h *MatchResultHandler) DeleteMatchResult(c *gin.Context) {
	matchResultID := c.Param("id")

	id, err := primitive.ObjectIDFromHex(matchResultID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Match Result ID format"})
		return
	}

	userIDStr, err := models.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userID, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}

	matchResult, err := models.GetMatchResultByID(c, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	if userID != matchResult.OrganizerID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not the organizer of this tournament"})
		return
	}

	err = models.DeleteMatchResult(c, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Match Result deleted successfully"})
}

func NewMatchResultHandler() *MatchResultHandler {
	return &MatchResultHandler{}
}