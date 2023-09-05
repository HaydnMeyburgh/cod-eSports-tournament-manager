package models

import (
	"context"
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/haydnmeyburgh/cod-eSports-tournament-manager/database"
	"github.com/haydnmeyburgh/cod-eSports-tournament-manager/internal/database"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-drivers/bson"
	"go.mongodb.org/mongo-drivers/bson/primitive"
)

type Team struct {
	ID primitive.ObjectID `bson:"_id,omitempty"`
	Name string `bson:"name" binding:"required"`
	Players []string `bson:"players"`
}

func CreateTeam(c *gin.Context, team *Team) (*Team, error) {
	collection := database.GetMongoClient().Database("esports-tournament-manager").Collection("teams")

	result, err := collection.InsertOne(c, team)
	if err != nil {
		return nil, err
	}

	team.ID = result.InsertedID.(primitive.ObjectID)
	return team, nil
}
