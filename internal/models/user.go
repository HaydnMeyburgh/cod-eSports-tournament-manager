package models

import (
	"errors"
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

// Queries DB to get a user by ID
func GetUserByID(id primitive.ObjectID) (*User, error) {
	collection := database.GetMongoClient().Database("esports-tournament-managemer").Collection("users")
	var user User
	err := collection.FindOne(nil, bson.M{"_id": id}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Retrieves a user by their email fromthe db
func GetUserByEmail(email string) (*User, error) {
	collection := database.GetMongoClient().Database("esports-tournament-manager").Collection("users")
	var user User
	err := collection.FindOne(nil, bson.M{"email": email}).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// Updates a user document in the db
func UpdateUserInDB(user *User) error {
	collection := database.GetMongoClient().Database("esports-tournament-manager").Collection("users")
	_, err := collection.ReplaceOne(nil, bson.M{"_id": user.ID}, user)
	if err != nil {
		return err
	}
	return nil
}

// Register a new user and store them in the database.
func RegisterUser(c *gin.Context, newUser *User) error {
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

// Logs in a user and returns a JWT token on successful login
func LoginUser(c *gin.Context) (string, error) {
	var loginUser struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&loginUser); err != nil {
		return "", err
	}

	// Check if user exists by email
	user, err := GetUserByEmail(loginUser.Email)
	if err != nil {
		return "", err
	}

	// Compare provided password with the hashed password for the database
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginUser.Password))
	if err != nil {
		return "", err
	}

	// Generate a JWT token for the user
	jwtSecret := os.Getenv("SECRET_KEY")
	if jwtSecret == "" {
		return "", errors.New("SECRET_KEY environment variable is not set")
	}
	token, err := auth.GenerateJWT(user.ID, jwtSecret)
	if err != nil {
		return "", err
	}

	// Set the token ass an HTTP cookie
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "jwtToken",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
	})

	return token, nil
}

// Logs out a user and and clears the JWT token cookie
func LogoutUser(c *gin.Context) {
	// Clear the JWT token cookie
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "jwtToken",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		MaxAge:   -1,
		SameSite: http.SameSiteNoneMode,
	})

	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}

// Updates the username and/or password of a user
func UpdateUser(c *gin.Context) {
	userID := c.MustGet("user_id").(string)

	// Convert userID string to primitive.ObjectID
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}

	var updateUser struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&updateUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if user exists by ID
	user, err := GetUserByID(objectID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Update username if provided
	if updateUser.Username != "" {
		user.Username = updateUser.Username
	}

	// Update password if provided
	if updateUser.Password != "" {
		// Hash new password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(updateUser.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash new password"})
			return
		}
		user.Password = string(hashedPassword)
	}

	// Update the user in the database
	err = UpdateUserInDB(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
}
