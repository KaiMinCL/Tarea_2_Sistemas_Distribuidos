package apiDB

import (
	"context"
	"fmt"
	"common/database"
	"common/models"
	"go.mongodb.org/mongo-driver/bson"
)

func GetClient(param_cliente models.ParametroCliente) (models.Cliente, error) {

	var cliente models.Cliente
	collection := database.GetDatabaseCollection("Clientes")

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
	collection := database.GetDatabaseCollection("Billeteras")

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
	collection := database.GetDatabaseCollection("Clientes")

	//define filter for the findOne command using the identification number
	filter := bson.M{"numero_identificacion": param_inicio.NumeroIdentificacion, "contrasena": param_inicio.Contrasena}

	err := collection.FindOne(context.Background(), filter).Decode(&cliente)
	if err != nil {
		fmt.Println("Error retrieving documents:", err)
		return false
	}
	return true
}
