package models_test

import (
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/haydnmeyburgh/cod-eSports-tournament-manager/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestRegistration(t *testing.T) {
	c, _ := gin.CreateTestContext(nil)

	user := models.User{
		Username: "TestUser",
		Email:    "test1@example.com",
		Password: "testpassword",
	}

	err := models.RegisterUser(c, &user)
	assert.NoError(t, err)

	retrievedUser, err := models.GetUserByEmail("test1@example.com")
	assert.NoError(t, err)
	assert.NotNil(t, retrievedUser)
	assert.NotEmpty(t, retrievedUser.ID)
	assert.Equal(t, user.Username, retrievedUser.Username)
	assert.Equal(t, user.Email, retrievedUser.Email)
}

func TestGetUserByID(t *testing.T) {
	c, _ := gin.CreateTestContext(nil)

	user := models.User{
		Username: "TestUser",
		Email:    "test2@example.com",
		Password: "testpassword",
	}

	err := models.RegisterUser(c, &user)
	assert.NoError(t, err)

	retrievedUser, err := models.GetUserByEmail("test2@example.com")
	assert.NoError(t, err)
	assert.NotNil(t, retrievedUser)

	userID := retrievedUser.ID
	fmt.Println("UserID: ", userID)

	retrievedUserByID, err := models.GetUserByID(userID)
	assert.NoError(t, err)
	assert.NotNil(t, retrievedUserByID)
	assert.Equal(t, user.Username, retrievedUserByID.Username)
	assert.Equal(t, user.Email, retrievedUserByID.Email)
	assert.Equal(t, user.Password, retrievedUserByID.Password)
}

func TestGetUserByEmail(t *testing.T) {
	c, _ := gin.CreateTestContext(nil)

	user := models.User{
		Username: "TestUser",
		Email:    "test3@example.com",
		Password: "testpassword",
	}

	err := models.RegisterUser(c, &user)
	assert.NoError(t, err)

	retrievedUser, err := models.GetUserByEmail("test3@example.com")
	assert.NoError(t, err)
	assert.NotNil(t, retrievedUser)
	assert.Equal(t, user.Username, retrievedUser.Username)
	assert.Equal(t, user.Email, retrievedUser.Email)
	assert.Equal(t, user.Password, retrievedUser.Password)
}

func TestUpdateUser(t *testing.T) {
	c, _ := gin.CreateTestContext(nil)

	user := models.User{
		Username: "testuser",
		Email:    "test4@example.com",
		Password: "testpassword",
	}

	err := models.RegisterUser(c, &user)
	assert.NoError(t, err)

	retrievedUser, err := models.GetUserByEmail("test4@example.com")
	assert.NoError(t, err)
	assert.NotNil(t, retrievedUser)

	userID := retrievedUser.ID

	updateUser := struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}{
		Username: "updateduser",
		Password: "updatedpassword",
	}

	err = models.UpdateUser(c, userID.Hex(), &updateUser)
	assert.NoError(t, err)

	retrievedUserByID, err := models.GetUserByID(userID)
	assert.NoError(t, err)
	assert.NotNil(t, retrievedUserByID)
	assert.Equal(t, updateUser.Username, retrievedUserByID.Username)
	// password should be different after update
	assert.NotEqual(t, user.Password, retrievedUserByID.Password)
}

func TestLoginUser(t *testing.T) {
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	user := models.User{
		Username: "testuser",
		Email:    "test5@example.com",
		Password: "testpassword",
	}

	err := models.RegisterUser(c, &user)
	assert.NoError(t, err)

	loginUser := struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}{
		Email:    "test5@example.com",
		Password: "testpassword",
	}

	token, err := models.LoginUser(c, &loginUser)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}
