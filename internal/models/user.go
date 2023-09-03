package models

import (
	"errors"
	"log"
	"net/http"
	"os"
	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/haydnmeyburgh/cod-eSports-tournament-manager/internal/auth"
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

// Logs in a user and returns a JWT token on successful login
func loginUser( c *gin.Context) {
	var loginUser struct {
		Email string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&loginUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if user exists by email (Need to create Get User By Email function)
	user, err := GetUserByEmail(loginUser.Email)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
	}

	// Compare provided password with the hashed password for the database
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginUser.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Credentials"})
		return
	}

	// Generate a JWT token for the user
	jwtSecret := os.Getenv("SECRET_KEY")
	if jwtSecret == "" {
		log.Fatal("MONGODB_URI environment variable is not set")
	}
	token, err := auth.GenerateJWT(user.ID.Hex(), jwtSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})

}

// LogoutUser logs out a user by invalidating their token
func LogoutUser(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}