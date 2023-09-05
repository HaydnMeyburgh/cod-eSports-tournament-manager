package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/haydnmeyburgh/cod-eSports-tournament-manager/internal/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MatchHandler struct {}

// Handlers creation of a new match
func (h *MatchHandler) CreateMatch(c *gin.Context) {
	userIDStr, err := models.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	var newMatch models.Match
	if err := c.ShouldBindJSON(&newMatch); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}

	newMatch.OrganizerID = userID
	createdMatch, err := models.CreateMatch(c, &newMatch)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Add the created match to the associated tournament
	tournamentID := newMatch.TournamentID
		err = models.AddMatchToTournament(c, createdMatch.ID, tournamentID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

	c.JSON(http.StatusCreated, createdMatch)
}

// handles the retrieval of a match
func (h *MatchHandler) GetMatchByID(c *gin.Context) {
	matchID := c.Param("id")

	id, err := primitive.ObjectIDFromHex(matchID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Match ID format"})
		return
	}

	match, err := models.GetMatchByID(c, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, match)
}

// Handles the updating of an existing match.
func (h *MatchHandler) UpdateMatch(c *gin.Context) {
	matchID := c.Param("id")

	id, err := primitive.ObjectIDFromHex(matchID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Match ID format"})
		return
	}

	var updatedMatch models.Match

	if err := c.ShouldBindJSON(&updatedMatch); err != nil {
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

	if userID != updatedMatch.OrganizerID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not the organiser of this match"})
	}

	err = models.UpdateMatch(c, id, &updatedMatch)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Match updated successfully"})
}

// Handles the deletion of a match by its ID.
func (h *MatchHandler) DeleteMatch(c *gin.Context) {
	matchID := c.Param("id")

	id, err := primitive.ObjectIDFromHex(matchID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Match ID format"})
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

	match, err := models.GetMatchByID(c, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	if userID != match.OrganizerID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not the organizer of this tournament"})
		return
	}

	err = models.DeleteMatch(c, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Match deleted successfully"})
}

func NewMatchHandler() *MatchHandler {
	return &MatchHandler{}
}