package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/haydnmeyburgh/cod-eSports-tournament-manager/internal/models"
)

type MatchHandler struct {}

// Handlers creation of a new match
func (h *MatchHandler) CreateMatchHandler(c *gin.Context) {
	var newMatch models.Match

	if err := c.ShouldBindJSON(&newMatch); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

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