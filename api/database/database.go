package database

import (
	"TrustBankApi/models"
	"context"
	_ "errors"
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

	collection := client.Database("banco").Collection(collectionName)

	return collection
}

func GetCliente(param_cliente models.ParametroCliente) (models.Cliente, error) {

	var cliente models.Cliente
	collection := getDatabaseCollection("cliente")

	//define filter for the findOne command using the identification number
	filter := bson.M{"numero_identificacion": param_cliente.NumeroIdentificacion}

	err := collection.FindOne(context.Background(), filter).Decode(&cliente)
	if err != nil {
		fmt.Println("Error retrieving documents:", err)
		return cliente, err
	}
	return cliente, nil
}

// dESDE AQUÍ FALTA IMPLEMENTAR LAS FUNCIONES
func PostSession(param_inicio models.ParametroInicio) (int, error) {
	var estado = 0
	return estado, nil
}

func PostDeposito(param_inicio models.ParametroInicio) (int, error) {
	var estado = 0
	return estado, nil
}

func PostTranferencia(param_inicio models.ParametroInicio) (int, error) {
	var estado = 0
	return estado, nil
}

func PostGiro(param_inicio models.ParametroInicio) (int, error) {
	var estado = 0
	return estado, nil
}
