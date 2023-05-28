package grpc

import (
	"log"
	"net"

	"google.golang.org/grpc"
)

func main() {

	// Iniciar el servidor en un puerto específico
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Error al iniciar el servidor: %v", err)
	}
	log.Println("Servidor gRPC iniciado en el puerto 50051")

	// Configuración del servidor gRPC
	server := grpc.NewServer()

	// Iniciar el servidor gRPC
	if err := server.Serve(listener); err != nil {
		log.Fatalf("Error al servir: %v", err)
	}
}
