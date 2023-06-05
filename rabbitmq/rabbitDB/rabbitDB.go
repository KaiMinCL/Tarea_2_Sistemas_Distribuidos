package rabbitDB

import (
	"context"
	"fmt"
	"log"
	"os"

	"../../common/database"
	"../../common/movimientos/movimientosGRPC"

	"go.mongodb.org/mongo-driver/bson"
	"google.golang.org/grpc"
)

func sendMovimientoRequest(connection *grpc.ClientConn, request *movimientosGRPC.MovimientoRequest) error {

	// Create a new instance of the MovimientosService client
	client := movimientosGRPC.NewMovimientosServiceClient(connection)

	// Call the RegistrarMovimiento RPC on the gRPC server
	response, err := client.RegistrarMovimiento(context.Background(), request)
	if err != nil {
		log.Fatalf("Failed to send MovimientoRequest: %v", err)
		return err
	}

	log.Printf("Received response: %s", response.Mensaje)

	return nil
}

// Create a gRPC client connection
func createGRPCClient() (*grpc.ClientConn, error) {
	var grpcHost = os.Getenv("GRPC_HOST")
	var grpcPort = os.Getenv("GRPC_PORT")

	grpcServerAddr := grpcHost + ":" + grpcPort

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

	// Create a new MovimientoRequest
	request := &movimientosGRPC.MovimientoRequest{
		NroCliente:    nroCliente,
		Monto:         monto,
		Divisa:        divisa,
		TipoOperacion: "deposito",
	}

	conn, err := createGRPCClient()
	if err != nil {
		log.Println("Failed to connect to gRPC server:", err)
		return err
	}

	err = sendMovimientoRequest(conn, request)
	if err != nil {
		log.Println("Failed to update movimientos:", err)
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

	// Create a new MovimientoRequest
	request := &movimientosGRPC.MovimientoRequest{
		NroClienteOrigen:  nroClienteOrigen,
		NroClienteDestino: nroClienteDestino,
		Monto:             monto,
		Divisa:            divisa,
		TipoOperacion:     "transferencia",
	}

	conn, err := createGRPCClient()
	if err != nil {
		log.Println("Failed to connect to gRPC server:", err)
		return err
	}

	err = sendMovimientoRequest(conn, request)
	if err != nil {
		log.Println("Failed to update movimientos:", err)
		return err
	}

	fmt.Printf("Transferencia en MongoDB: nroClienteOrigen=%s, nroClienteDestino=%s, monto=%f, divisa=%s\n", nroClienteOrigen, nroClienteDestino, monto, divisa)
	return nil
}

// Función para realizar el giro en la billetera del cliente
func Withdraw(nroCliente string, monto float64, divisa string) error {
	_, collection := database.GetDatabaseCollection("Billeteras")

	filter := bson.M{"nro_cliente": nroCliente}
	update := bson.M{"$inc": bson.M{"saldo": monto}}

	_, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		log.Println("Failed to update wallet in MongoDB:", err)
		return err
	}

	// Create a new MovimientoRequest
	request := &movimientosGRPC.MovimientoRequest{
		NroCliente:    nroCliente,
		Monto:         monto,
		Divisa:        divisa,
		TipoOperacion: "giro",
	}

	conn, err := createGRPCClient()
	if err != nil {
		log.Println("Failed to connect to gRPC server:", err)
		return err
	}

	err = sendMovimientoRequest(conn, request)
	if err != nil {
		log.Println("Failed to update movimientos:", err)
		return err
	}

	fmt.Printf("Giro en MongoDB: nroCliente=%s, monto=%f, divisa=%s\n", nroCliente, monto, divisa)
	return nil
}
