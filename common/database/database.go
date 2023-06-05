package database

import (
	"common/models"
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
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

func GetWallet(numeroCliente, divisaBilletera string) (models.Billetera, error) {
	_, collection := GetDatabaseCollection("Billeteras")

	filter := bson.M{"nro_cliente": numeroCliente, "divisa": divisaBilletera}

	var billetera models.Billetera
	err := collection.FindOne(context.Background(), filter).Decode(&billetera)
	if err != nil {
		return models.Billetera{}, err
	}

	return billetera, nil
}
