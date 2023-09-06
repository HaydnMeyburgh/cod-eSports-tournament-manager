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

	newTournament.OrganiserID = userID
	createdTournament, err := models.CreateTournament(c, &newTournament)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, createdTournament)

}

// Handles getting all tournaments by organiser ID
func (h *TournamentHandler) GetTournamentsByOrganiserID(c *gin.Context) {
	userID, err := models.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	organiserID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invaliad User ID format"})
		return
	}

	tournaments, err := models.GetTournamentsByOrganiserID(c, organiserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tournaments)
}

// handles the retrieval of a tournament by id
func (h *TournamentHandler) GetTournamentByID(c *gin.Context) {
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

// Handles updating of a tournament
func (h *TournamentHandler) UpdateTournament(c *gin.Context) {
	tournamentID := c.Param("id")

	id, err := primitive.ObjectIDFromHex(tournamentID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Tournament ID format"})
		return
	}

	var updatedTournament models.Tournament

	if err := c.ShouldBindJSON(&updatedTournament); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Ensure that the user making the request is the tournament organizer
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

	if userID != updatedTournament.OrganiserID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not the organizer of this tournament"})
		return
	}

	err = models.UpdateTournament(c, id, &updatedTournament)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Tournament updated successfully"})
}

// Handles the deletion of a tournament by its ID
func (h *TournamentHandler) DeleteTournament(c *gin.Context) {
	tournamentID := c.Param("id")

	id, err := primitive.ObjectIDFromHex(tournamentID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Tournament ID format"})
		return
	}

	// Ensure that the user making the request is the tournament organizer
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

	tournament, err := models.GetTournamentByID(c, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	if userID != tournament.OrganiserID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not the organizer of this tournament"})
		return
	}

	err = models.DeleteTournament(c, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Tournament deleted successfully"})
}

func NewTournamentHandler() *TournamentHandler {
	return &TournamentHandler{}
}