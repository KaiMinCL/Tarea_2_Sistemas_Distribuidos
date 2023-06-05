package main

import (
	"context"
	"log"
	"net"
	"os"
	"time"

	"common/database"
	pb "common/movimientos/movimientosGRPC"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedMovimientosServiceServer
}

func (s *server) RegistrarMovimiento(ctx context.Context, req *pb.MovimientoRequest) (*pb.MovimientoResponse, error) {
	// Create a connection to MongoDB
	client, err := mongo.NewClient(options.Client().ApplyURI(getMongoURI()))
	if err != nil {
		log.Fatalf("Error creating MongoDB client: %v", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		log.Fatalf("Error connecting to MongoDB: %v", err)
	}
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			log.Fatalf("Error disconnecting from MongoDB: %v", err)
		}
	}()

	// Seleccionar la base de datos y la colección
	db := client.Database("TrustBank")
	collection := db.Collection("Movimientos")

	var movimiento = bson.M{}

	if req.TipoOperacion == "deposito" || req.TipoOperacion == "giro" {
		billetera, _ := database.GetWallet(req.NroCliente, req.Divisa)

		// Crear el documento para insertar en MongoDB
		movimiento = bson.M{
			"nro_cliente":    req.GetNroCliente(),
			"monto":          req.GetMonto(),
			"divisa":         req.GetDivisa(),
			"tipo_operacion": req.GetTipoOperacion(),
			"fecha_hora":     time.Now(),
			"id_billetera":   billetera.Id,
		}

	} else {
		billeteraOrigen, _ := database.GetWallet(req.NroClienteOrigen, req.Divisa)

		billeteraDestino, _ := database.GetWallet(req.NroClienteDestino, req.Divisa)

		// Crear el documento para insertar en MongoDB
		movimiento = bson.M{
			"nro_cliente_origen":   req.GetNroClienteOrigen(),
			"nro_cliente_destino":  req.GetNroClienteDestino(),
			"monto":                req.GetMonto(),
			"divisa":               req.GetDivisa(),
			"tipo_operacion":       req.GetTipoOperacion(),
			"fecha_hora":           time.Now(),
			"id_billetera_origen":  billeteraOrigen.Id,
			"id_billetera_destino": billeteraDestino.Id,
		}
	}

	// Insertar el documento en la colección
	_, err = collection.InsertOne(ctx, movimiento)
	if err != nil {
		log.Fatalf("Error al insertar el movimiento en MongoDB: %v", err)
	}

	log.Printf("Movimiento almacenado: " + req.GetTipoOperacion())

	return &pb.MovimientoResponse{Mensaje: "Movimiento registrado exitosamente"}, nil
}

func getMongoURI() string {
	mongoURI := os.Getenv("DB_CONNECTION_STRING")
	if mongoURI == "" {
		log.Fatal("DB_CONNECTION_STRING environment variable not set")
	}
	return mongoURI
}

func main() {

	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	var grpcHost = os.Getenv("GRPC_HOST")
	var grpcPort = os.Getenv("GRPC_PORT")

	// Iniciar el servidor en un puerto específico
	listener, err := net.Listen("tcp", grpcHost+":"+grpcPort)
	if err != nil {
		log.Fatalf("Error al iniciar el servidor: %v", err)
	}
	log.Println("Servidor gRPC iniciado en el puerto " + grpcPort)

	s := grpc.NewServer()
	pb.RegisterMovimientosServiceServer(s, &server{})
	log.Printf("gRPC server started on port %s", grpcPort)
	if err := s.Serve(listener); err != nil {
		log.Fatalf("Error starting gRPC server: %v", err)
	}
}
