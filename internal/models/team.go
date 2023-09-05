package models

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/haydnmeyburgh/cod-eSports-tournament-manager/internal/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-drivers/bson/primitive"
)

type Team struct {
	ID      primitive.ObjectID `bson:"_id,omitempty"`
	Name    string             `bson:"name" binding:"required"`
	Players []string           `bson:"players"`
}

// Creates a new Team
func CreateTeam(c *gin.Context, team *Team) (*Team, error) {
	collection := database.GetMongoClient().Database("esports-tournament-manager").Collection("teams")

	result, err := collection.InsertOne(c, team)
	if err != nil {
		return nil, err
	}

	team.ID = result.InsertedID.(primitive.ObjectID)
	return team, nil
}

// Retrieves a team by id
func GetTeamByID(c *gin.Context, id primitive.ObjectID) (*Team, error) {
	collection := database.GetMongoClient().Database("esports-tournament-manager").Collection("team")

	var team Team
	err := collection.FindOne(c, bson.M{"_id": id}).Decode(&team)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("Team not found")
		}
		return nil, err
	}

	return &team, nil
}

// Updates an existing team
func UpdateTeam(c *gin.Context, id primitive.ObjectID, updatedTeam *Team) error {
	collection := database.GetMongoClient().Database("esports-tournament-manager").Collection("team")

	update := bson.M{"Set": updatedTeam}
	_, err := collection.UpdateOne(c, bson.M{"_id": id}, update)
	if err != nil {
		return err
	}

	return nil
}

// Delets a team by it's ID
func DeleteTeam(c *gin.Context, id primitive.ObjectID) error {
	collection := database.GetMongoClient().Database("esports-tournament-manager").Collection("team")

	_, err := collection.DeleteOne(c, bson.M{"_id": id})
	if err != nil {
		return err
	}

	return nil
}