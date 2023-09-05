package models_test

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/haydnmeyburgh/cod-eSports-tournament-manager/internal/database"
	"github.com/haydnmeyburgh/cod-eSports-tournament-manager/internal/models"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
)

func cleanUpTestData() error {
	collection := database.GetMongoClient().Database("esports-tournament-manager").Collection("teams")

	filter := bson.M{}

	_, err := collection.DeleteMany(context.TODO(), filter)
	if err != nil {
		return err
	}

	return nil
}

func TestMain(m *testing.M) {
	if err := godotenv.Load("../../.env"); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	if err := database.ConnectToMongoDB(); err != nil {
		log.Fatalf("Error initializing Mongodb: %v", err)
	}

	exitCode := m.Run()

	if err := cleanUpTestData(); err != nil {
		log.Fatalf("Error cleaning up test data: %v", err)
	}

	database.GetMongoClient().Disconnect(context.TODO())

	os.Exit(exitCode)
}

func TestCreateTeam(t *testing.T) {
	team := models.Team{
		Name:    "TestTeam",
		Players: []string{"player1", "player2", "player3", "player4", "player5"},
	}

	createdTeam, err := models.CreateTeam(nil, &team)
	assert.NoError(t, err)
	assert.NotNil(t, createdTeam)
	assert.NotEmpty(t, createdTeam.ID)

}

func TestGetTeamByID(t *testing.T) {
	team := models.Team{
		Name: "TestTeam",
		Players: []string{"Player1", "Player2", "Player3", "Player4", "Player5"},
	}

	createdTeam, err := models.CreateTeam(nil, &team)
	assert.NoError(t, err)
	assert.NotNil(t, createdTeam)
	assert.NotEmpty(t, createdTeam.ID)

	teamID := createdTeam.ID

	retrievedTeam, err := models.GetTeamByID(nil, teamID)
	assert.NoError(t, err)
	assert.NotNil(t, retrievedTeam)
	assert.Equal(t, team.Name, retrievedTeam.Name)
	assert.Equal(t, team.Players, retrievedTeam.Players)
}

func TestUpdateTeam(t *testing.T) {
	team := models.Team{
		Name: "TestTeam",
		Players: []string{"Player1", "Player2", "Player3", "Player4", "Player5"},
	}

	createdTeam, err := models.CreateTeam(nil, &team)
	assert.NoError(t, err)
	assert.NotNil(t, createdTeam)
	assert.NotEmpty(t, createdTeam.ID)

	teamID := createdTeam.ID

	updatedTeam := models.Team{
		Name: "UpdatedTeam",
		Players: []string{"UpdatedPlayer1", "Player2", "Player3", "Player4", "Player5"},
	}

	err = models.UpdateTeam(nil, teamID, &updatedTeam)
	assert.NoError(t, err)

	retrievedTeam, err := models.GetTeamByID(nil, teamID)
	assert.NoError(t, err)
	assert.NotNil(t, retrievedTeam)
	assert.Equal(t, updatedTeam.Name, retrievedTeam.Name)
	assert.Equal(t, updatedTeam.Players, retrievedTeam.Players)
}

