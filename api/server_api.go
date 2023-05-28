package main

import (
	"TrustBankApi/controllers"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load the environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	var port = os.Getenv("PORT")
	var server = os.Getenv("SERVER")

	//Define the new router for the Gin framework
	router := gin.Default()

	//Set the endpoints
	router.GET("/api/cliente", controllers.GetClient)
	router.POST("/api/inicio_sesion", controllers.SessionHandler)
	router.POST("/api/deposito", controllers.DepositHandler)
	router.POST("/api/transferancia", controllers.TransferHandler)
	router.POST("/api/giro", controllers.WithdrawHandler)

	router.Run(server + ":" + port)
}
