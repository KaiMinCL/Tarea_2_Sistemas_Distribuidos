package database

import (
	"context"
	"fmt"
	"log"
	"os"

	"common/models"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Función para enviar una solicitud gRPC al servidor de movimientos
/*
func sendMovimientoRequest(client movimientos.MovimientoRequest, movimiento *movimientos.MovimientoRequest) {
	_, err := client.RegistrarMovimiento(context.Background(), movimiento)
	if err != nil {
		log.Printf("Failed to send movimiento request: %v", err)
	} else {
		log.Println("Movimiento request sent successfully")
	}
}
*/

// Connect to MongoDB and retrieve the collection needed
func getDatabaseCollection(collectionName string) (*mongo.Client, *mongo.Collection) {
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

	return client, collection
}

func Deposit(nroCliente string, monto float64, divisa string) error {
	_, collection := getDatabaseCollection("billeteras")
	filter := bson.M{"nro_cliente": nroCliente}
	update := bson.M{"$inc": bson.M{"saldo": monto}}

	_, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		log.Println("Failed to update wallet in MongoDB:", err)
		return err
	}

	fmt.Printf("Depósito en MongoDB: nroCliente=%s, monto=%f, divisa=%s\n", nroCliente, monto, divisa)
	return nil
}

func Transfer(nroClienteOrigen string, nroClienteDestino string, monto float64, divisa string) error {
	client, collection := getDatabaseCollection("billeteras")

	session, err := client.StartSession()
	if err != nil {
		log.Fatal("Failed to start MongoDB session:", err)
	}
	defer session.EndSession(context.Background())

	err = session.StartTransaction()
	if err != nil {
		log.Fatal("Failed to start MongoDB transaction:", err)
	}

	// Update the sender's wallet
	filterSender := bson.M{"nro_cliente": nroClienteOrigen}
	updateSender := bson.M{"$inc": bson.M{"saldo": -monto}}

	_, err = collection.UpdateOne(context.Background(), filterSender, updateSender)
	if err != nil {
		session.AbortTransaction(context.Background())
		log.Println("Failed to update sender's wallet in MongoDB:", err)
		return err
	}

	// Update the recipient's wallet
	filterRecipient := bson.M{"nro_cliente": nroClienteDestino}
	updateRecipient := bson.M{"$inc": bson.M{"saldo": monto}}

	_, err = collection.UpdateOne(context.Background(), filterRecipient, updateRecipient)
	if err != nil {
		session.AbortTransaction(context.Background())
		log.Println("Failed to update recipient's wallet in MongoDB:", err)
		return err
	}

	session.CommitTransaction(context.Background())

	fmt.Printf("Transferencia en MongoDB: nroClienteOrigen=%s, nroClienteDestino=%s, monto=%f, divisa=%s\n", nroClienteOrigen, nroClienteDestino, monto, divisa)
	return nil
}

// Función para realizar el giro en la billetera del cliente
func Withdraw(nroCliente string, monto float64, divisa string) error {
	// Obtener la conexión a MongoDB y la colección necesaria
	client, collection := getDatabaseCollection("billeteras")
	defer client.Disconnect(context.Background())

	// Iniciar una sesión de transacción en MongoDB
	session, err := client.StartSession()
	if err != nil {
		log.Println("Failed to start MongoDB session:", err)
		return err
	}
	defer session.EndSession(context.Background())

	// Iniciar una transacción en la sesión
	err = session.StartTransaction()
	if err != nil {
		log.Println("Failed to start MongoDB transaction:", err)
		return err
	}

	// Obtener la billetera del cliente
	filter := bson.M{"nro_cliente": nroCliente}
	update := bson.M{"$inc": bson.M{"saldo": -monto}}

	_, err = collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		session.AbortTransaction(context.Background())
		log.Println("Failed to update wallet in MongoDB:", err)
		return err
	}

	// Crear el documento del movimiento
	movimiento := models.Movimiento{
		NroCliente: nroCliente,
		Monto:      monto,
		Divisa:     divisa,
		Tipo:       "giro",
	}

	// Insertar el documento del movimiento en la colección de movimientos
	movimientosCollection := client.Database("TrustBank").Collection("movimientos")
	_, err = movimientosCollection.InsertOne(context.Background(), movimiento)
	if err != nil {
		session.AbortTransaction(context.Background())
		log.Println("Failed to insert movement in MongoDB:", err)
		return err
	}

	// Confirmar la transacción en MongoDB
	err = session.CommitTransaction(context.Background())
	if err != nil {
		log.Println("Failed to commit MongoDB transaction:", err)
		return err
	}

	log.Printf("Giro en MongoDB: nroCliente=%s, monto=%f, divisa=%s\n", nroCliente, monto, divisa)
	return nil
}
