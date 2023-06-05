package database

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Connect to MongoDB and retrieve the collection needed
func GetDatabaseCollection(collectionName string) (*mongo.Client, *mongo.Collection) {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	var connectionString = os.Getenv("DB_CONNECTION_STRING")

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(connectionString))

	if err != nil {
		panic(err)
	}

	collection := client.Database("TrustBank").Collection(collectionName)

	return client, collection
}
