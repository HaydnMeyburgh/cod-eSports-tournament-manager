// package for the database setup and initialisation
package database

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
)

var client *mongo.Client

func ConnectToMongoDB() error {
	// Getting the MongoDB URI from the environment variable
	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		log.Fatal("MONGODB_URI environment variable is not set")
	}

	clientOptions := options.Client().ApplyURI(uri)
	
	// Creating the MongoDB client
	c, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return err
	}
	
	// Pinging the server to verify the connection
	if err = c.Ping(context.TODO(), nil); err != nil {
		return err
	}

	// Assigning the client to the package-level 'client' variable
	client = c

	log.Println("Connected to MongoDB!")

	return nil
}

func GetMongoClient() *mongo.Client {
	return client
}