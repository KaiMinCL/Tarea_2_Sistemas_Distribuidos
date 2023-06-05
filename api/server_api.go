package main

import (
	"TrustBankAPI/controllers"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/streadway/amqp"
)

var (
	rabbitConn  *amqp.Connection
	rabbitCh    *amqp.Channel
	rabbitQueue string
)

func main() {
	// Load the environment variables
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	var httpHost = os.Getenv("HTTP_HOST")
	var httpPort = os.Getenv("HTTP_PORT")

	err = controllers.InitRabbitMQ()
	if err != nil {
		log.Fatalf("Error initializing RabbitMQ: %v", err)
	}
	defer controllers.CloseRabbitMQ()

	//Define the new router for the Gin framework
	router := gin.Default()

	//Set the endpoints
	router.GET("/api/cliente", controllers.GetClient)
	router.POST("/api/inicio_sesion", controllers.SessionHandler)
	router.POST("/api/deposito", controllers.DepositHandler)
	router.POST("/api/transferencia", controllers.TransferHandler)
	router.POST("/api/giro", controllers.WithdrawHandler)

	router.Run(httpHost + ":" + httpPort)
}
