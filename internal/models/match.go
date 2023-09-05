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
	TournamentID primitive.ObjectID `bson:"tournament_id" binding:"required"`
	Team1ID   primitive.ObjectID `bson:"team1_id" binding:"required"`
	Team2ID   primitive.ObjectID `bson:"team2_id" binding:"required"`
	Score1    int                `bson:"score1"`
	Score2    int                `bson:"score2"`
}

// - CreateMatch
func CreateMatch(c *gin.Context, match *Match) (*Match, error) {
	collection := database.GetMongoClient().Database("esports-tournament-manager").Collection("matches")

	result, err := collection.InsertOne(ctx, match)
	if err != nil {
		return nil, err
	}

	match.ID = result.InsertedID.(primitive.ObjectID)
	return match, nil
}

// - GetMatchByID
// - UpdateMatch
// - DeleteMatch

// - CreateMatchHandler
// - GetMatchHandler
// - UpdateMatchHandler
// - DeleteMatchHandler