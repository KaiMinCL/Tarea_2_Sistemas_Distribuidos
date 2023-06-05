package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	pb "common/movimientos"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
)

type movimientosServer struct {
	db *mongo.Client
}

func NewMovimientosServer(db *mongo.Client) *movimientosServer {
	return &movimientosServer{
		db: db,
	}
}

// Implement the methods of the movimientos service
func (s *movimientosServer) RegisterMovimientosService(ctx context.Context, request *pb.MovimientoRequest) (*pb.MovimientoResponse, error) {
	// Process the received movimiento request
	nroClienteOrigen := request.GetNroClienteOrigen()
	nroClienteDestino := request.GetNroClienteDestino()
	monto := request.GetMonto()
	divisa := request.GetDivisa()
	tipoOperacion := request.GetTipoOperacion()

	// Create a BSON document from the received movimiento
	movimiento := bson.D{
		{Key: "nro_cliente_origen", Value: nroClienteOrigen},
		{Key: "nro_cliente_destino", Value: nroClienteDestino},
		{Key: "monto", Value: monto},
		{Key: "divisa", Value: divisa},
		{Key: "tipo_operacion", Value: tipoOperacion},
	}

	// Insert the movimiento document into the MongoDB collection
	collection := s.db.Database("TrustBank").Collection("Movimientos")
	_, err := collection.InsertOne(ctx, movimiento)
	if err != nil {
		log.Printf("Failed to insert movimiento into MongoDB: %v", err)
	}

	// Return a response to the client
	response := &pb.MovimientoResponse{
		Mensaje: fmt.Sprintf("Movimiento registrado: %s -> %s, Monto: %f", request.NroClienteOrigen, request.NroClienteDestino, request.Monto),
	}

	return response, nil
}

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	var grpcHost = os.Getenv("GRPC_HOST")
	var grpcPort = os.Getenv("GRPC_PORT")
	var dbConnectionString = os.Getenv("DB_CONNECTION_STRING")

	/// Set up the MongoDB client
	client, err := mongo.NewClient(options.Client().ApplyURI(dbConnectionString))
	if err != nil {
		log.Fatalf("Failed to create MongoDB client: %v", err)
	}

	ctx := context.TODO()
	err = client.Connect(ctx)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer client.Disconnect(ctx)

	// Iniciar el servidor en un puerto específico
	listener, err := net.Listen("tcp", grpcHost+":"+grpcPort)
	if err != nil {
		log.Fatalf("Error al iniciar el servidor: %v", err)
	}
	log.Println("Servidor gRPC iniciado en el puerto " + grpcPort)

	// Create a new gRPC server
	server := grpc.NewServer()

	///// BUGS BUGS

	// Register the movimientos server implementation with the gRPC server
	//movimientosServer := NewMovimientosServer(client)
	//pb.RegisterMovimientosService(server, movimientosServer)

	///// BUGS BUGS

	// Iniciar el servidor gRPC
	if err := server.Serve(listener); err != nil {
		log.Fatalf("Error al servir: %v", err)
	}
}
