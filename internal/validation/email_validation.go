package validation

import (
	"context"

	"github.com/haydnmeyburgh/cod-eSports-tournament-manager/internal/database"
	"go.mongodb.org/mongo-driver/bson"
)

func EmailExistsInDB(email string) (bool, error) {
	// Handle to mongoDB collection
	collection := database.GetMongoClient().Database("eSports-tournament-manager").Collection("users")

	// Define a filter to find documents with the given email
	filter := bson.M{"email": email}

	// Count documents with given email
	count, err := collection.CountDocuments(context.TODO(), filter)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}
