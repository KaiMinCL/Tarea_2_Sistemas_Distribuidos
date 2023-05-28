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
	router.GET("/api/cliente", controllers.GetCliente)
	router.POST("/api/inicio_sesion", controllers.PostSession)
	router.POST("/api/deposito", controllers.PostDeposito)
	router.POST("/api/transferancia", controllers.PostTranferencia)
	router.POST("/api/giro", controllers.PostGiro)

	router.Run(server + ":" + port)
}
