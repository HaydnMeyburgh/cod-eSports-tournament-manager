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

type Match struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	TournamentID primitive.ObjectID `bson:"tournament_id,omitempty"`
	OrganiserID  primitive.ObjectID `bson:"organiser_id" binding:"required"`
	Team1ID      primitive.ObjectID `bson:"team1_id" binding:"required"`
	Team2ID      primitive.ObjectID `bson:"team2_id" binding:"required"`
	Score1       int                `bson:"score1"`
	Score2       int                `bson:"score2"`
	WebSocketHub *realtimemanager.WebSocketHub
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

	// Construct the WebSocket message for match creation
	message := map[string]interface{}{
		"action":        "match_created",
		"match_id":      match.ID.Hex(),
		"tournament_id": match.TournamentID.Hex(),
		"team1_id":      match.Team1ID.Hex(),
		"team2_id":      match.Team2ID.Hex(),
		"score1":        match.Score1,
		"score2":        match.Score2,
	}

	// Marshal the message to JSON
	messageJSON, err := json.Marshal(message)
	if err != nil {
		log.Printf("Error marshaling WebSocket message: %v", err)
		return match, nil
	}

	log.Printf("Broadcasting WebSocket message: %s", messageJSON)
	// Broadcast the message to WebSocket clients
	match.WebSocketHub.Broadcast(messageJSON)

	return match, nil
}

// Get Matches by organiserID
func GetMatchesByOrganiserID(c *gin.Context, organiserID primitive.ObjectID) ([]*Match, error) {
	collection := database.GetMongoClient().Database("esports-tournament-manager").Collection("matches")

	cursor, err := collection.Find(c, bson.M{"organiser_id": organiserID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(c)

	var matches []*Match
	for cursor.Next(c) {
		var match Match
		if err := cursor.Decode(&match); err != nil {
			return nil, err
		}
		matches = append(matches, &match)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return matches, nil
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

	// Construct the WebSocket message for match update
	message := map[string]interface{}{
		"action":        "match_updated",
		"match_id":      id.Hex(),
		"tournament_id": updatedMatch.TournamentID.Hex(),
		"team1_id":      updatedMatch.Team1ID.Hex(),
		"team2_id":      updatedMatch.Team2ID.Hex(),
		"score1":        updatedMatch.Score1,
		"score2":        updatedMatch.Score2,
	}

	// Marshal the message to JSON
	messageJSON, err := json.Marshal(message)
	if err != nil {
		return err
	}

	log.Printf("Broadcasting WebSocket message: %s", messageJSON)
	// Broadcast the message to WebSocket clients
	updatedMatch.WebSocketHub.Broadcast(messageJSON)

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
