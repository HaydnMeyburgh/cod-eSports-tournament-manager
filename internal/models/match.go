package models

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/haydnmeyburgh/cod-eSports-tournament-manager/internal/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Match struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	TournamentID primitive.ObjectID `bson:"tournament_id,omitempty"`
	OrganizerID primitive.ObjectID   `bson:"organizer_id" binding:"required"`
	Team1ID   primitive.ObjectID `bson:"team1_id" binding:"required"`
	Team2ID   primitive.ObjectID `bson:"team2_id" binding:"required"`
	Score1    int                `bson:"score1"`
	Score2    int                `bson:"score2"`
}

// add a team to a tournament
func AddMatchToTournament(c *gin.Context, matchID, tournamentID primitive.ObjectID) error {
	collection := database.GetMongoClient().Database("esports-tournament-manager").Collection("tournament")

	update := bson.M{"$push": bson.M{"match": matchID}}

	_, err := collection.UpdateOne(c, bson.M{"_id": tournamentID}, update)
	if err != nil {
		return err
	}

	return nil
}

// - CreateMatch
func CreateMatch(c *gin.Context, match *Match) (*Match, error) {
	collection := database.GetMongoClient().Database("esports-tournament-manager").Collection("matches")

	result, err := collection.InsertOne(c, match)
	if err != nil {
		return nil, err
	}

	match.ID = result.InsertedID.(primitive.ObjectID)
	return match, nil
}

// - GetMatchByID
func GetMatchByID(c *gin.Context, id primitive.ObjectID) (*Match, error) {
	collection := database.GetMongoClient().Database("esports-tournament-manager").Collection("matches")

	var match Match
	err := collection.FindOne(c, bson.M{"_id": id}).Decode(&match)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("Match not found")
		}
		return nil, err
	}

	return &match, nil
}

// - UpdateMatch
func UpdateMatch(c *gin.Context, id primitive.ObjectID, updatedMatch *Match) error {
	collection := database.GetMongoClient().Database("esports-tournament-manager").Collection("matches")

	update := bson.M{"$set": updatedMatch}
	_, err := collection.UpdateOne(c, bson.M{"_id": id}, update)
	if err != nil {
		return err
	}

	return nil
}

// - DeleteMatch
func DeleteMatch(c *gin.Context, id primitive.ObjectID) error {
	collection := database.GetMongoClient().Database("esports-tournament-manager").Collection("matches")

	_, err := collection.DeleteOne(c, bson.M{"_id": id})
	if err != nil {
		return err
	}

	return nil
}

// - CreateMatchHandler
// - GetMatchHandler
// - UpdateMatchHandler
// - DeleteMatchHandler