package main

import (
	"context"
	"log"
	"net"
	"os"

	"common/movimientos"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
)

type MovimientosServer struct {
	MongoClient   *mongo.Client
	MongoDatabase *mongo.Database
}

// Implementa los métodos del servicio
func (s *MovimientosServer) RegistrarMovimiento(ctx context.Context, request *protos.MovimientoRequest) (*protos.MovimientoResponse, error) {

	// Obtener la colección "movimientos"
	collection := s.MongoDatabase.Collection("movimientos")

	// Crear un documento BSON a partir del movimiento recibido
	movimiento := bson.D{
		{Key: "nro_cliente_origen", Value: request.NroClienteOrigen},
		{Key: "nro_cliente_destino", Value: request.NroClienteDestino},
		{Key: "monto", Value: request.Monto},
		{Key: "divisa", Value: request.Divisa},
		{Key: "tipo_operacion", Value: request.TipoOperacion},
	}

	// Insertar el documento en la colección "movimientos"
	_, err := collection.InsertOne(ctx, movimiento)
	if err != nil {
		log.Printf("Error al insertar el movimiento en MongoDB: %v", err)
		return nil, err
	}

	log.Printf("Movimiento registrado: %+v", request)

	// Devolver la misma solicitud recibida como respuesta
	return request, nil
}

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	var grpcHost = os.Getenv("GRPC_HOST")
	var grpcPort = os.Getenv("GRPC_PORT")
	var dbConnectionString = os.Getenv("DB_CONNECTION_STRING")

	// Establecer la conexión con MongoDB
	clientOptions := options.Client().ApplyURI(dbConnectionString)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatalf("Error al conectar con MongoDB: %v", err)
	}

	database := client.Database("TrustBank")

	// Crear una instancia del servidor gRPC de movimientos
	server := grpc.NewServer()

	// Registrar el servidor de movimientos
	movimientos.RegisterMovimientosServiceServer(server, &MovimientosServer{
		MongoClient:   client,
		MongoDatabase: database,
	})

	// Iniciar el servidor en un puerto específico
	listener, err := net.Listen("tcp", grpcHost+":"+grpcPort)
	if err != nil {
		log.Fatalf("Error al iniciar el servidor: %v", err)
	}
	log.Println("Servidor gRPC iniciado en el puerto " + grpcPort)

	// Iniciar el servidor gRPC
	if err := server.Serve(listener); err != nil {
		log.Fatalf("Error al servir: %v", err)
	}
}
