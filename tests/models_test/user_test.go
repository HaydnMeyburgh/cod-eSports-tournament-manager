package models_test

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/haydnmeyburgh/cod-eSports-tournament-manager/internal/database"
	"github.com/haydnmeyburgh/cod-eSports-tournament-manager/internal/models"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

func cleanupTestData() error {
	collection := database.GetMongoClient().Database("esports-tournament-manager").Collection("users")

	filter := bson.M{}

	_, err := collection.DeleteMany(context.TODO(), filter)
	if err != nil {
		return err
	}

	return nil
}

func TestMain(m *testing.M) {
	// Load environment variables from the .env file
	if err := godotenv.Load("../../.env"); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	err := database.ConnectToMongoDB()
	if err != nil {
		log.Fatalf("Error initialising MongoDB: %v", err)
	}

	exitCode := m.Run()

	if err = cleanupTestData(); err != nil {
		log.Fatalf("Error cleaning up test data: %v", err)
	}

	database.GetMongoClient().Disconnect(context.TODO())

	os.Exit(exitCode)
}

func TestGetUserByID(t *testing.T) {
	user := models.User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "testpassword",
	}
	collection := database.GetMongoClient().Database("esports-tournament-manager").Collection("users")
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

	collection := database.GetMongoClient().Database("esports-tournament-manager").Collection("users")
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

	collection := database.GetMongoClient().Database("esports-tournament-manager").Collection("users")
	res, err := collection.InsertOne(nil, user)
	if err != nil {
		t.Fatalf("Error inserting user: %v", err)
	}
	userID := res.InsertedID.(primitive.ObjectID)

	updateUser := struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}{
		Username: "updateduser",
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
	newUser := models.User{
		Username: "newuser",
		Email:    "newuser@example.com",
		Password: "newpassword",
	}

	c, _ := gin.CreateTestContext(nil)

	err := models.RegisterUser(c, &newUser)
	if err != nil {
		t.Fatalf("Error registering new user: %v", err)
	}

	retrievedUser, err := models.GetUserByEmail("newuser@example.com")
	if err != nil {
		t.Fatalf("Error fetching registered user: %v", err)
	}

	assert.Equal(t, newUser.Username, retrievedUser.Username)
	assert.Equal(t, newUser.Email, retrievedUser.Email)
	assert.Equal(t, newUser.Password, retrievedUser.Password)
}

func TestLoginUser(t *testing.T) {

	user := models.User{
		Username: "testuser2",
		Email:    "test2@example.com",
		Password: "testpassword2",
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	user.Password = string(hashedPassword)

	fmt.Println("User before database isertion:", user)
	collection := database.GetMongoClient().Database("esports-tournament-manager").Collection("users")
	res, err := collection.InsertOne(nil, user)

	fmt.Println("result after insertion:", res)
	if err != nil {
		t.Fatalf("Error inserting user: %v", err)
	}

	loginUser := struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}{
		Email:    "test2@example.com",
		Password: "testpassword2",
	}
	fmt.Println("hashed Password:", user.Password)

	token, err := models.LoginUser(nil, &loginUser)
	if err != nil {
		t.Fatalf("Error logging in user: %v", err)
	}

	assert.NotEmpty(t, token)
}
