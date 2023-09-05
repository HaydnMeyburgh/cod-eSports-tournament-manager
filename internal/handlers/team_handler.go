package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/haydnmeyburgh/cod-eSports-tournament-manager/internal/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TeamHandler struct{}

// Handler for create team
func (h *TeamHandler) CreateTeam(c *gin.Context) {
	var newTeam models.Team

	if err := c.ShouldBindJSON(&newTeam); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tournamentID := newTeam.TournamentID

	createdTeam, err := models.CreateTeam(c, &newTeam)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := models.AddTeamToTournament(c, tournamentID, createdTeam.ID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, createdTeam)
}

// Handler to for getting a team by ID
func (h *TeamHandler) GetTeamByID(c *gin.Context) {
	teamID := c.Param("id")

	objectID, err := primitive.ObjectIDFromHex(teamID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid team ID format"})
		return
	}

	team, err := models.GetTeamByID(c, objectID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, team)
}

// Handler for updating a team
func (h *TeamHandler) UpdateTeam(c *gin.Context) {
	teamID := c.Param("id")

	objectID, err := primitive.ObjectIDFromHex(teamID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid team ID format"})
		return
	}

	var updatedTeam models.Team

	if err :=c.ShouldBindJSON(&updatedTeam); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := models.UpdateTeam(c, objectID, &updatedTeam); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Team successfully updated"})
}

// Handler to delete team
func (h *TeamHandler) DeleteTeam(c *gin.Context) {
	teamID := c.Param("id")

	objectID, err := primitive.ObjectIDFromHex(teamID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid team ID format"})
		return
	}

	if err := models.DeleteTeam(c, objectID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return 
	}

	c.JSON(http.StatusOK, gin.H{"message": "Team deleted successfully"})
}

func NewTeamHandler() *TeamHandler {
	return &TeamHandler{}
}