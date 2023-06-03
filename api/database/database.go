package database

import (
	"TrustBankApi/models"
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Connect to MongoDB and retrieve the collection needed
func getDatabaseCollection(collectionName string) *mongo.Collection {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	var CONNECTION_STRING = os.Getenv("CONNECTION_STRING")

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(CONNECTION_STRING))

	if err != nil {
		panic(err)
	}

	collection := client.Database("TrustBank").Collection(collectionName)

	return collection
}

func GetClient(param_cliente models.ParametroCliente) (models.Cliente, error) {

	var cliente models.Cliente
	collection := getDatabaseCollection("Clientes")

	//define filter for the findOne command using the identification number
	filter := bson.M{"numero_identificacion": param_cliente.NumeroIdentificacion}

	err := collection.FindOne(context.Background(), filter).Decode(&cliente)
	if err != nil {
		fmt.Println("Error al recuperar los documentos:", err)
		return cliente, err
	}

	return cliente, nil
}

func GetWallet(numeroCliente string) (models.Billetera, error) {
	collection := getDatabaseCollection("Billeteras")

	filter := bson.M{"nro_cliente": numeroCliente}
	var billetera models.Billetera
	err := collection.FindOne(context.Background(), filter).Decode(&billetera)
	if err != nil {
		return models.Billetera{}, err
	}

	return billetera, nil
}

// DESDE AQUÍ FALTA IMPLEMENTAR LAS FUNCIONES
func VerifySession(param_inicio models.ParametroInicio) bool {
	var cliente models.Cliente
	collection := getDatabaseCollection("Clientes")

	//define filter for the findOne command using the identification number
	filter := bson.M{"numero_identificacion": param_inicio.NumeroIdentificacion, "contrasena": param_inicio.Contrasena}

	err := collection.FindOne(context.Background(), filter).Decode(&cliente)
	if err != nil {
		fmt.Println("Error retrieving documents:", err)
		return false
	}
	return true
}
