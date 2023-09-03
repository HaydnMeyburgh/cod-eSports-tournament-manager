package models

import (
	"errors"
	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/haydnmeyburgh/cod-eSports-tournament-manager/internal/database"
	"github.com/haydnmeyburgh/cod-eSports-tournament-manager/internal/validation"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Username string             `bson:"username"`
	Email    string             `bson:"email"`
	Password string             `bson:"password"`
}

var (
	emailRegex        = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	minPasswordLength = 8
)

// Register a new user and store them in the database.
func RegisterUser(c *gin.Context) error {
	var newUser User
	if err := c.ShouldBindJSON(&newUser); err != nil {
		return err
	}

	// Validate email
	if !emailRegex.MatchString(newUser.Email) {
		return errors.New("invalid email format")
	}

	// Check if the email already exists in the db
	emailExists, err := validation.EmailExistsInDB(newUser.Email)
	if err != nil {
		return err
	}
	if emailExists {
		return errors.New("Email already exists")
	}

	// Validate password length - 9 characters
	if len(newUser.Password) < minPasswordLength {
		return errors.New("Password must be at least 8 characters")
	}

	// Use Bcrypt to hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	newUser.Password = string(hashedPassword)

	// Get a handle to the "users" collection
	collection := database.GetMongoClient().Database("esports-tournament-manager").Collection("users")

	// Insert new users document into collection
	_, err = collection.InsertOne(c, newUser)
	if err != nil {
		return err
	}

	return nil
}

func GetUserByID(id primitive.ObjectID) (*User, error) {
	collection := database.GetMongoClient().Database("esports-tournament-managemer").Collection("users")
	var user User
	err := collection.FindOne(nil, bson.M{"_id": id}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
