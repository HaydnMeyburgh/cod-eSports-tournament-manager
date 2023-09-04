package models

import (
	"testing"
	"time"

	"github.com/stretchr/testing/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestRegisterUser(t *testing.T) {
	newUser := &User{
		Username: "NewUser",
		Email: "newuser@example.com",
		Password: "newpassword",
	}

	err := RegisterUser(nil, newUser)
	assert.NoError(t, err)

	registeredUser, err := GetUserByEmail("newuser@example.com")
	assert.NoError(t, err)
	assert.Equal(t, "NewUser", registeredUser.Username)
	assert.Equal(t, "newuser@example.com", registeredUser.Email)
}