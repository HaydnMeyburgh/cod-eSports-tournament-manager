package models_test

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/haydnmeyburgh/cod-eSports-tournament-manager/internal/auth"
	"github.com/haydnmeyburgh/cod-eSports-tournament-manager/internal/database"
	"github.com/haydnmeyburgh/cod-eSports-tournament-manager/internal/models"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CloseTestMongoDB() {
	// Close the mongoDB client connection
	if database.GetMongoClient() != nil {
		err := database.GetMongoClient().Disconnect(context.Background())
		if err != nil {
			log.Printf("Error closing test MongoDB client: %v", err)
		}
	}
}

func TestMain(m *testing.M) {
	err := database.ConnectToTestMongoDB()
	if err != nil {
		log.Fatalf("Error initialising MongoDB: %v", err)
	}

	exitCode := m.Run()

	CloseTestMongoDB()

	os.Exit(exitCode)
}

func TestGetUserByID(t *testing.T) {
	user := models.User{
		Username: "testuser",
		Email: "test@example.com",
		Password: "testpassword",
	}
	collection := database.GetMongoClient().Database("test_database").Collection("users")
	res, err := collection.InsertOne(nil, user)
	if err != nil {
		t.Fatalf("Error inserting user: %v", err)
	}
	userID := res.InsertedID.(primitive.ObjectID)

	retrievedUser, err := models.GetUserByID(userID)
	if err != nil {
		t.Fatalf("Error fetching user by ID: %v", err)
	}

	assert.Equal(t, user.Username, retrievedUser.Username)
	assert.Equal(t, user.Email, retrievedUser.Email)
	assert.Equal(t, user.Password, retrievedUser.Password)
}