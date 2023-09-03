package validation

import (
	"context"

	"github.com/haydnmeyburgh/cod-eSports-tournament-manager/internal/database"
	"go.mongodb.org/mongo-driver/bson"
)

func EmailExistsInDB(email string) (bool, error) {
	collection := database.GetMongoClient().Database("esports-tournament-manager").Collection("users")
	filter := bson.M{"email": email}
	count, err := collection.CountDocuments(context.Background(), filter)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
