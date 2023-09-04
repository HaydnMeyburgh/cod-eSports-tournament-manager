// package for the database setup and initialisation
package database

import (
	"context"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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
	c, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return err
	}

	// Pinging the server to verify the connection
	if err = c.Ping(context.Background(), nil); err != nil {
		return err
	}

	// Assigning the client to the package-level 'client' variable
	client = c

	log.Println("Connected to MongoDB!")

	// Initialise the database (Create it if it doesn't exist)
	if err := InitialiseDatabase("esports-tournament-manager"); err != nil {
		return err
	}

	return nil
}

func ConnectToTestMongoDB() error {
	// Getting the MongoDB URI from the environment variable
	uri := os.Getenv("MONGODBTEST_URI")
	if uri == "" {
		log.Fatal("MONGODBTEST_URI environment variable is not set")
	}

	clientOptions := options.Client().ApplyURI(uri)

	// Creating the MongoDB client
	c, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return err
	}

	// Pinging the server to verify the connection
	if err = c.Ping(context.Background(), nil); err != nil {
		return err
	}

	// Assigning the client to the package-level 'client' variable
	client = c

	log.Println("Connected to MongoDB!")

	// Initialise the database (Create it if it doesn't exist)
	if err := InitialiseDatabase("test_database"); err != nil {
		return err
	}

	return nil
}

func GetMongoClient() *mongo.Client {
	return client
}

func InitialiseDatabase(dbName string) error {
	// List existing databases
	databases, err := client.ListDatabaseNames(context.TODO(), bson.M{})
	if err != nil {
		return err
	}

	// Check if database already exists
	exists := false
	for _, db := range databases {
		if db == dbName {
			exists = true
			break
		}
	}

	// If the database does not exist, create it
	if !exists {
		err := client.Database(dbName).CreateCollection(context.TODO(), "users")
		if err != nil {
			return err
		}
	}

	return nil
}
