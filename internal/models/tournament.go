package models

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/haydnmeyburgh/cod-eSports-tournament-manager/internal/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Tournament struct {
	ID          primitive.ObjectID   `bson:"_id,omitempty"`
	Name        string               `bson:"name" binding:"required"`
	Description string               `bson:"description"`
	StartDate   string               `bson:"start_date" binding:"required"`
	EndDate     string               `bson:"end_date" binding:"required"`
	OrganizerID primitive.ObjectID   `bson:"organizer_id" binding:"required"`
	Teams       []primitive.ObjectID `bson:"teams"`
	Matches     []primitive.ObjectID `bson:"matches"`
}
// Adds teams to tournaments
func AddTeamsToTournament(c *gin.Context, tournamentID primitive.ObjectID, teamIDs []primitive.ObjectID) error {
	// Add teams to the tournament
	collection := database.GetMongoClient().Database("esports-tournament-manager").Collection("tournaments")

	update := bson.M{"$push": bson.M{"teams": bson.M{"$each": teamIDs}}}
	_, err := collection.UpdateOne(c, bson.M{"_id": tournamentID}, update)
	if err != nil {
		return err
	}

	return nil
}

// - Creates a new tournament
func CreateTournament(c *gin.Context, tournament *Tournament) (*Tournament, error) {
	collection := database.GetMongoClient().Database("esports-tournament-manager").Collection("tournaments")

	result, err := collection.InsertOne(c, tournament)
	if err != nil {
		return nil, err
	}

	tournament.ID = result.InsertedID.(primitive.ObjectID)
	return tournament, nil
}

// - GetTournamentByID
func GetTournamentByID(c *gin.Context, id primitive.ObjectID) (*Tournament, error) {
	collection := database.GetMongoClient().Database("esports-tournament-manager").Collection("tournaments")

	var tournament Tournament
	err := collection.FindOne(c, bson.M{"_id": id}).Decode(&tournament)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("Tournament not found")
		}
		return nil, err
	}

	return &tournament, nil
}

// - UpdateTournament
func UpdateTournament(c *gin.Context, id primitive.ObjectID, updatedTournament *Tournament) error {
	collection := database.GetMongoClient().Database("esports-tournament-manager").Collection("tournaments")

	update := bson.M{"$set": updatedTournament}
	_, err := collection.UpdateOne(c, bson.M{"_id": id}, update)
	if err != nil {
		return err
	}

	return nil
}

// - DeleteTournament
func DeleteTournament(c *gin.Context, id primitive.ObjectID) error {
	collection := database.GetMongoClient().Database("esports-tournament-manager").Collection("tournaments")

	_, err := collection.DeleteOne(c, bson.M{"_id": id})
	if err != nil {
		return err
	}

	return nil
}

// - CreateTournamentHandler
// - GetTournamentHandler
// - UpdateTournamentHandler
// - DeleteTournamentHandler
