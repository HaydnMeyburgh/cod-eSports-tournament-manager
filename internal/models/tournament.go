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

type Tournament struct {
	ID           primitive.ObjectID   `bson:"_id,omitempty"`
	Name         string               `bson:"name" binding:"required"`
	Description  string               `bson:"description"`
	StartDate    string               `bson:"start_date" binding:"required"`
	EndDate      string               `bson:"end_date" binding:"required"`
	OrganiserID  primitive.ObjectID   `bson:"organiser_id" binding:"required"`
	Teams        []primitive.ObjectID `bson:"teams"`
	Matches      []primitive.ObjectID `bson:"matches"`
	WebSocketHub *realtimemanager.WebSocketHub
}

// Adds teams to tournaments
func AddTeamsToTournament(c *gin.Context, tournamentID primitive.ObjectID, teamIDs []primitive.ObjectID) error {
	collection := database.GetMongoClient().Database("esports-tournament-manager").Collection("tournaments")

	update := bson.M{"$push": bson.M{"teams": bson.M{"$each": teamIDs}}}
	_, err := collection.UpdateOne(c, bson.M{"_id": tournamentID}, update)
	if err != nil {
		return err
	}

	return nil
}

// Add matches to a tournament
func AddMatchesToTournament(c *gin.Context, tournamentID primitive.ObjectID, matchIDs []primitive.ObjectID) error {
	collection := database.GetMongoClient().Database("esports-tournament-manager").Collection("tournaments")

	update := bson.M{"$push": bson.M{"matches": bson.M{"$each": matchIDs}}}
	_, err := collection.UpdateOne(c, bson.M{"_id": tournamentID}, update)
	if err != nil {
		return err
	}

	return nil
}

// Get user id from context
func GetUserIDFromContext(c *gin.Context) (string, error) {
	userID, exists := c.Get("user_id")
	if !exists {
		return "", errors.New("User ID not found in context")
	}

	userIDStr, ok := userID.(string)
	if !ok {
		return "", errors.New("User ID is not a string")
	}

	return userIDStr, nil
}

// - Creates a new tournament
func CreateTournament(c *gin.Context, tournament *Tournament) (*Tournament, error) {
	collection := database.GetMongoClient().Database("esports-tournament-manager").Collection("tournaments")

	result, err := collection.InsertOne(c, tournament)
	if err != nil {
		return nil, err
	}

	tournament.ID = result.InsertedID.(primitive.ObjectID)

	if len(tournament.Teams) > 0 {
		err = AddTeamsToTournament(c, tournament.ID, tournament.Teams)
		if err != nil {
			return nil, err
		}
	}

	if len(tournament.Matches) > 0 {
		err = AddMatchesToTournament(c, tournament.ID, tournament.Matches)
		if err != nil {
			return nil, err
		}
	}

	// Broadcast the message to WebSocket clients
	message := map[string]interface{}{
		"action":        "tournament_created",
		"tournament_id": tournament.ID.Hex(),
		"name":          tournament.Name,
		"desccription":  tournament.Description,
		"start_date":    tournament.StartDate,
		"end_date":      tournament.EndDate,
	}

	// Marshal the message to JSON
	messageJSON, err := json.Marshal(message)
	if err != nil {
		log.Printf("Error marshaling WebSocket message: %v", err)
		return tournament, nil
	}

	tournament.WebSocketHub.Broadcast(messageJSON)

	return tournament, nil
}

// GetAllTournaments
func GetTournamentsByOrganiserID(c *gin.Context, OrganiserID primitive.ObjectID) ([]*Tournament, error) {
	collection := database.GetMongoClient().Database("esports-tournament-manager").Collection("tournaments")

	cursor, err := collection.Find(c, bson.M{"organiser_id": OrganiserID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(c)

	var tournaments []*Tournament
	for cursor.Next(c) {
		var tournament Tournament
		if err := cursor.Decode(&tournament); err != nil {
			return nil, err
		}
		tournaments = append(tournaments, &tournament)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return tournaments, nil
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

	// Broadcast the message to WebSocket clients
	message := map[string]interface{}{
		"action":        "tournament_created",
		"tournament_id": updatedTournament.ID.Hex(),
		"name":          updatedTournament.Name,
		"desccription":  updatedTournament.Description,
		"start_date":    updatedTournament.StartDate,
		"end_date":      updatedTournament.EndDate,
	}

	// Marshal the message to JSON
	messageJSON, err := json.Marshal(message)
	if err != nil {
		return err
	}

	updatedTournament.WebSocketHub.Broadcast(messageJSON)

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
