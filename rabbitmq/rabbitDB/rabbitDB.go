package rabbitDB

import (
	"context"
	"fmt"
	"log"

	"common/database"
	"common/models"

	"go.mongodb.org/mongo-driver/bson"
	"google.golang.org/grpc"
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

// Create a gRPC client connection
func createGRPCClient() (*grpc.ClientConn, error) {
	grpcServerAddr := "<gRPC server address>" // Replace with the actual gRPC server address

	// Dial the gRPC server
	conn, err := grpc.Dial(grpcServerAddr, grpc.WithInsecure())
	if err != nil {
		return nil, fmt.Errorf("failed to connect to gRPC server: %v", err)
	}

	return conn, nil
}

func Deposit(nroCliente string, monto float64, divisa string) error {
	_, collection := database.GetDatabaseCollection("Billeteras")
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
	client, collection := database.GetDatabaseCollection("Billeteras")

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
	client, collection := database.GetDatabaseCollection("Billeteras")
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
