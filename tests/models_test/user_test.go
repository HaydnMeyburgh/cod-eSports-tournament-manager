package models_test

import (
	"context"
	"log"
	"os"
	"testing"

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
		Email:    "test@example.com",
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

func TestGetUserByEmail(t *testing.T) {
	user := models.User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "testpassword",
	}

	collection := database.GetMongoClient().Database("test_database").Collection("users")
	_, err := collection.InsertOne(nil, user)
	if err != nil {
		t.Fatalf("Error inserting user: %v", err)
	}

	retrievedUser, err := models.GetUserByEmail("test@example.com")
	if err != nil {
		t.Fatalf("Error fetching user by email: %v", err)
	}

	assert.Equal(t, user.Username, retrievedUser.Username)
	assert.Equal(t, user.Email, retrievedUser.Email)
	assert.Equal(t, user.Password, retrievedUser.Password)
}

func TestUpdateUser(t *testing.T) {
	user := models.User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "testpassword",
	}

	collection := database.GetMongoClient().Database("test_database").Collection("users")
	res, err := collection.InsertOne(nil, user)
	if err != nil {
		t.Fatalf("Error inserting user: %v", err)
	}
	userID := res.InsertedID.(primitive.ObjectID)

	updateUser := struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}{
		Username: "",
		Password: "updatedpassword",
	}

	err = models.UpdateUser(nil, userID.Hex(), &updateUser)
	if err != nil {
		t.Fatalf("Error updating user: %v", err)
	}

	retrievedUser, err := models.GetUserByID(userID)
	if err != nil {
		t.Fatalf("Error fetching updated user: %v", err)
	}

	assert.Equal(t, updateUser.Username, retrievedUser.Username)
	// password should be different after update
	assert.NotEqual(t, user.Password, retrievedUser.Password)
}

func TestRegistration(t *testing.T) {
	newUser := &models.User{
		Username: "newuser",
		Email:    "newuser@example.com",
		Password: "newpassword",
	}

	err := models.RegisterUser(nil, newUser)
	if err != nil {
		t.Fatalf("Error registering new user: %v", err)
	}

	retrievedUser, err := models.GetUserByEmail("newuser@example.com")
	if err != nil {
		t.Fatalf("Error fetching registered user: %v", err)
	}

	assert.Equal(t, newUser.Username, retrievedUser.Username)
	assert.Equal(t, newUser.Email, retrievedUser.Email)
	// Password should come back hasehd
	assert.NotEqual(t, newUser.Password, retrievedUser.Password)
}

func TestLoginUser(t *testing.T) {
	user := models.User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "testpassword",
	}

	collection := database.GetMongoClient().Database("test_database").Collection("users")
	_, err := collection.InsertOne(nil, user)
	if err != nil {
		t.Fatalf("Error inserting user: %v", err)
	}

	loginUser := struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}{
		Email:    "test@example.com",
		Password: "testpassword",
	}

	token, err := models.LoginUser(nil, &loginUser)
	if err != nil {
		t.Fatalf("Error logging in user: %v", err)
	}

	assert.NotEmpty(t, token)
}
