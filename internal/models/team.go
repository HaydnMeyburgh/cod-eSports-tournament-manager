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

type Team struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	Name         string             `bson:"name" binding:"required"`
	OrganiserID  primitive.ObjectID `bson:"organiser_id" binding:"required"`
	Players      []string           `bson:"players"`
	TournamentID primitive.ObjectID `bson:"tournament_id,omitempty"`
	WebSocketHub *realtimemanager.WebSocketHub
}

// add a team to a tournament
func AddTeamToTournament(c *gin.Context, teamID, tournamentID primitive.ObjectID) error {
	collection := database.GetMongoClient().Database("esports-tournament-manager").Collection("tournament")

	update := bson.M{"$push": bson.M{"teams": teamID}}

	_, err := collection.UpdateOne(c, bson.M{"_id": tournamentID}, update)
	if err != nil {
		return err
	}

	return nil
}

// Creates a new Team
func CreateTeam(c *gin.Context, team *Team) (*Team, error) {
	collection := database.GetMongoClient().Database("esports-tournament-manager").Collection("teams")

	result, err := collection.InsertOne(c, team)
	if err != nil {
		return nil, err
	}

	team.ID = result.InsertedID.(primitive.ObjectID)

	// Broadcast the message to WebSocket clients
	message := map[string]interface{}{
		"action":        "team_created",
		"match_id":      team.ID.Hex(),
		"name":          team.Name,
		"players":       team.Players,
		"tournament_id": team.TournamentID.Hex(),
	}

	// Marshal the message to JSON
	messageJSON, err := json.Marshal(message)
	if err != nil {
		log.Printf("Error marshaling WebSocket message: %v", err)
		return team, nil
	}

	team.WebSocketHub.Broadcast(messageJSON)

	return team, nil
}

// Get Matches by organiser id
func GetTeamsByOrganiserID(c *gin.Context, organiserID primitive.ObjectID) ([]*Team, error) {
	collection := database.GetMongoClient().Database("esports-tournament-manager").Collection("teams")

	cursor, err := collection.Find(c, bson.M{"organiser_id": organiserID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(c)

	var teams []*Team
	for cursor.Next(c) {
		var team Team
		if err := cursor.Decode(&team); err != nil {
			return nil, err
		}
		teams = append(teams, &team)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return teams, nil
}

// Retrieves a team by id
func GetTeamByID(c *gin.Context, id primitive.ObjectID) (*Team, error) {
	collection := database.GetMongoClient().Database("esports-tournament-manager").Collection("teams")

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
	collection := database.GetMongoClient().Database("esports-tournament-manager").Collection("teams")

	update := bson.M{"$set": updatedTeam}
	_, err := collection.UpdateOne(c, bson.M{"_id": id}, update)
	if err != nil {
		return err
	}

	// Broadcast the message to WebSocket clients
	message := map[string]interface{}{
		"action":        "team_created",
		"match_id":      updatedTeam.ID.Hex(),
		"name":          updatedTeam.Name,
		"players":       updatedTeam.Players,
		"tournament_id": updatedTeam.TournamentID.Hex(),
	}

	// Marshal the message to JSON
	messageJSON, err := json.Marshal(message)
	if err != nil {
		return err
	}

	updatedTeam.WebSocketHub.Broadcast(messageJSON)

	return nil
}

// Delets a team by it's ID
func DeleteTeam(c *gin.Context, id primitive.ObjectID) error {
	collection := database.GetMongoClient().Database("esports-tournament-manager").Collection("teams")

	_, err := collection.DeleteOne(c, bson.M{"_id": id})
	if err != nil {
		return err
	}

	return nil
}
