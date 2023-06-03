package main

import (
	"RabbitMqMovements/database"
	"RabbitMqMovements/models"
	"encoding/json"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/streadway/amqp"
	"google.golang.org/grpc"
)

func main() {

	// Load the environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	var rabbitPort = os.Getenv("RABBITMQ_PORT")
	var rabbitHost = os.Getenv("RABBITMQ_HOST")
	var rabbitUsername = os.Getenv("RABBITMQ_USERNAME")
	var rabbitPassword = os.Getenv("RABBITMQ_PASSWORD")
	var rabbitQueue = os.Getenv("RABBITMQ_QUEUE_NAME")
	var grpcHost = os.Getenv("GRPC_HOST")
	var grpcPort = os.Getenv("GRPC_PORT")

	// Configurar la conexión gRPC
	grpc, err := grpc.Dial(grpcHost+":"+grpcPort, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to gRPC server: %v", err)
	}
	defer grpc.Close()

	// Establecer conexión con RabbitMQ
	rabbit, err := amqp.Dial("amqp://" + rabbitUsername + ":" + rabbitPassword + "@" + rabbitHost + ":" + rabbitPort + "/")
	if err != nil {
		log.Fatalf("Error al establecer la conexión con RabbitMQ: %v", err)
	}
	defer rabbit.Close()

	// Crear un canal de comunicación
	ch, err := rabbit.Channel()
	if err != nil {
		log.Fatalf("Error al crear el canal: %v", err)
	}
	defer ch.Close()

	// Declarar la cola de mensajes
	queue, err := ch.QueueDeclare(
		rabbitQueue, // Nombre de la cola
		false,       // No duradera
		false,       // No eliminar cuando no hay consumidores
		false,       // No exclusiva
		false,       // No esperar confirmación
		nil,         // Argumentos adicionales
	)
	if err != nil {
		log.Fatalf("Error al declarar la cola: %v", err)
	}

	// Consumir mensajes de la cola
	msgs, err := ch.Consume(
		queue.Name, // Nombre de la cola
		"",         // Nombre del consumidor
		true,       // Auto-acknowledge
		false,      // No-exclusivo
		false,      // No-local
		false,      // No-wait
		nil,        // Argumentos adicionales
	)
	if err != nil {
		log.Fatalf("Error al registrar el consumidor: %v", err)
	}

	// Procesar mensajes recibidos
	forever := make(chan bool)

	go func() {
		for d := range msgs {
			// Deserializar el mensaje en una estructura Movimiento
			var movimiento models.Movimiento
			err := json.Unmarshal(d.Body, &movimiento)
			if err != nil {
				log.Println("Failed to unmarshal message:", err)
				d.Ack(false) // Descartar el mensaje
				continue
			}

			// Actualizar la billetera del cliente según el tipo de operación
			switch movimiento.Tipo {
			case "deposito":
				err = database.Deposit(movimiento.NroClienteDestino, movimiento.Monto, movimiento.Divisa)
			case "transferencia":
				err = database.Transfer(movimiento.NroClienteOrigen, movimiento.NroClienteDestino, movimiento.Monto, movimiento.Divisa)
			case "giro":
				err = database.Withdraw(movimiento.NroClienteOrigen, movimiento.Monto, movimiento.Divisa)
			default:
				log.Println("Invalid operation type:", movimiento.Tipo)
			}

			if err != nil {
				log.Println("Failed to process operation:", err)
			}

			d.Ack(false) // Confirmar el procesamiento del mensaje
		}
	}()

	log.Printf("Esperando mensajes. Presiona CTRL+C para salir.")
	<-forever
}
