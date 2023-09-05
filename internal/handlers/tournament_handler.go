package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/haydnmeyburgh/cod-eSports-tournament-manager/internal/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TournamentHandler struct {}

// Handles creationg of a new tournament
func (h *TournamentHandler) CreateTournament(c * gin.Context) {
	userIDStr, err := models.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	var newTournament models.Tournament
	if err := c.ShouldBindJSON(&newTournament); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}

	newTournament.OrganizerID = userID
	createdTournament, err := models.CreateTournament(c, &newTournament)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, createdTournament)

}

// handles the retrieval of a tournament by id
func (h *TournamentHandler) GetTournamentHandler(c *gin.Context) {
	tournamentID := c.Param("id")

	id, err := primitive.ObjectIDFromHex(tournamentID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Tournament ID format"})
		return
	}

	tournament, err := models.GetTournamentByID(c, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tournament)
}