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
