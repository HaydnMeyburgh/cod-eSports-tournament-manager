package models

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/haydnmeyburgh/cod-eSports-tournament-manager/internal/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MatchResult struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	MatchID   primitive.ObjectID `bson:"match_id" binding:"required"`
	OrganiserID primitive.ObjectID   `bson:"organiser_id" binding:"required"`
	WinnerID  primitive.ObjectID `bson:"winner_id"`
	LoserID   primitive.ObjectID `bson:"loser_id"`
	Score1    int                `bson:"score1"`
	Score2    int                `bson:"score2"` 
}

// MatchResult-related functions
// - CreateMatchResult
func CreateMatchResult(c *gin.Context, matchResult *MatchResult) (*MatchResult, error) {
	collection := database.GetMongoClient().Database("esports-tournament-manager").Collection("match_results")

	result, err := collection.InsertOne(c, matchResult)
	if err != nil {
		return nil, err
	}

	matchResult.ID = result.InsertedID.(primitive.ObjectID)
	return matchResult, nil
}
// Gets match results by organiser id
func GetMatchResultsByOrganiserID(c *gin.Context, organiserID primitive.ObjectID) ([]*MatchResult, error) {
	collection := database.GetMongoClient().Database("esports-tournament-manager").Collection("match_results")


	cursor, err := collection.Find(c, bson.M{"organiser_id": organiserID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(c)

	var matchResults []*MatchResult
	for cursor.Next(c) {
		var matchResult MatchResult
		if err := cursor.Decode(&matchResult); err != nil {
			return nil, err
		}
		matchResults = append(matchResults, &matchResult)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return matchResults, nil
}

// - GetMatchResultByID
func GetMatchResultByID(c *gin.Context, id primitive.ObjectID) (*MatchResult, error) {
	collection := database.GetMongoClient().Database("esports-tournament-manager").Collection("match_results")

	var matchResult MatchResult
	err := collection.FindOne(c, bson.M{"_id": id}).Decode(&matchResult)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("Match result not found")
		}
		return nil, err
	}

	return &matchResult, nil
}

// - UpdateMatchResult
func UpdateMatchResult(c *gin.Context, id primitive.ObjectID, updatedMatchResult *MatchResult) error {
	collection := database.GetMongoClient().Database("esports-tournament-manager").Collection("match_results")

	update := bson.M{"$set": updatedMatchResult}
	_, err := collection.UpdateOne(c, bson.M{"_id": id}, update)
	if err != nil {
		return err
	}

	return nil
}

// - DeleteMatchResult
func DeleteMatchResult(c *gin.Context, id primitive.ObjectID) error {
	collection := database.GetMongoClient().Database("esports-tournament-manager").Collection("match_results")

	_, err := collection.DeleteOne(c, bson.M{"_id": id})
	if err != nil {
		return err
	}

	return nil
}
// MatchResult Handler (match_result_handler.go):
// - CreateMatchResultHandler
// - GetMatchResultHandler
// - UpdateMatchResultHandler
// - DeleteMatchResultHandler