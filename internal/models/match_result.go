package models

import (
	"encoding/json"
	"errors"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/haydnmeyburgh/cod-eSports-tournament-manager/internal/database"
	"github.com/haydnmeyburgh/cod-eSports-tournament-manager/internal/realtimemanager"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MatchResult struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	MatchID      primitive.ObjectID `bson:"match_id" binding:"required"`
	OrganiserID  primitive.ObjectID `bson:"organiser_id" binding:"required"`
	WinnerID     primitive.ObjectID `bson:"winner_id"`
	LoserID      primitive.ObjectID `bson:"loser_id"`
	WinnerScore  int                `bson:"winner_score"`
	LoserScore   int                `bson:"loser_score"`
	WebSocketHub *realtimemanager.WebSocketHub
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

	// Construct and broadcast the WebSocket message
	message := map[string]interface{}{
		"action":          "match_result_created",
		"match_result_id": matchResult.ID.Hex(),
		"match_id":        matchResult.MatchID.Hex(),
		"organiser_id":    matchResult.OrganiserID.Hex(),
		"winner_id":       matchResult.WinnerID.Hex(),
		"loser_id":        matchResult.LoserID.Hex(),
		"winner_score":    matchResult.WinnerScore,
		"loser_sccore":    matchResult.LoserScore,
	}

	// Marshal the message to JSON
	messageJSON, err := json.Marshal(message)
	if err != nil {
		log.Printf("Error marshaling WebSocket message: %v", err)
		return matchResult, nil
	}
	
	log.Printf("Broadcasting WebSocket message: %s", messageJSON)
	// Broadcast the message to WebSocket clients
	matchResult.WebSocketHub.Broadcast(messageJSON)

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

	// Construct and broadcast the WebSocket message
	message := map[string]interface{}{
		"action":          "match_result_created",
		"match_result_id": updatedMatchResult.ID.Hex(),
		"match_id":        updatedMatchResult.MatchID.Hex(),
		"organiser_id":    updatedMatchResult.OrganiserID.Hex(),
		"winner_id":       updatedMatchResult.WinnerID.Hex(),
		"loser_id":        updatedMatchResult.LoserID.Hex(),
		"winner_score":    updatedMatchResult.WinnerScore,
		"loser_sccore":    updatedMatchResult.LoserScore,
	}

	// Marshal the message to JSON
	messageJSON, err := json.Marshal(message)
	if err != nil {
		log.Printf("Error marshaling WebSocket message: %v", err)
		return err
	}

	log.Printf("Broadcasting WebSocket message: %s", messageJSON)
	// Broadcast the message to WebSocket clients
	updatedMatchResult.WebSocketHub.Broadcast(messageJSON)

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
