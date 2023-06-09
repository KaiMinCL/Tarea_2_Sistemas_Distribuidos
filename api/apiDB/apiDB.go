package apiDB

import (
	"common/database"
	"common/models"
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
)

func GetClient(param_cliente models.ParametroCliente) (models.Cliente, error) {

	var cliente models.Cliente
	_, collection := database.GetDatabaseCollection("Clientes")

	filter := bson.M{"numero_identificacion": param_cliente.NumeroIdentificacion}

	err := collection.FindOne(context.Background(), filter).Decode(&cliente)
	if err != nil {
		fmt.Println("Error al recuperar los documentos:", err)
		return cliente, err
	}

	return cliente, nil
}

func VerifySession(param_inicio models.ParametroInicio) bool {
	var cliente models.Cliente
	_, collection := database.GetDatabaseCollection("Clientes")

	//define filter for the findOne command using the identification number
	filter := bson.M{"numero_identificacion": param_inicio.NumeroIdentificacion, "contrasena": param_inicio.Contrasena}

	err := collection.FindOne(context.Background(), filter).Decode(&cliente)
	if err != nil {
		fmt.Println("Error retrieving documents:", err)
		return false
	}
	return true
}
