package controllers

import (
	"TrustBankApi/database"
	"TrustBankApi/models"
	"fmt"
	_ "log"
	_ "math/rand"
	"net/http"
	_ "strings"
	_ "time"

	"github.com/gin-gonic/gin"
)

func GetCliente(c *gin.Context) {

	//Varianble for the client's parameters
	var param_cliente models.ParametroCliente

	//Getting the parameters
	if err := c.BindJSON(&param_cliente); err != nil {

		fmt.Println("Problem with Json bindng:", err)
		c.AbortWithStatus(http.StatusBadRequest)
	} else {
		// Call the GetCliente function from the models package with the parameters
		cliente, err := database.GetCliente(param_cliente)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"estado": "cliente_no_encontrado",
			})
		} else {
			c.JSON(http.StatusOK, cliente)
		}
	}
}

func PostSession(c *gin.Context) {
	//Varianble for the client's parameters
	var param_inicio models.ParametroInicio

	//Getting the parameters
	if err := c.BindJSON(&param_inicio); err != nil {

		fmt.Println("Problem with Json bindng:", err)
		c.AbortWithStatus(http.StatusBadRequest)
	} else {
		// Call the GetCliente function from the models package with the parameters
		estado, err := database.PostSession(param_inicio)
		if err != nil {
			c.JSON(401, gin.H{
				"estado": "no exitoso",
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"estado": "exitoso",
			})
		}
	}
}

func PostDeposito(c *gin.Context) {
	//Varianble for the client's parameters
	var param_deposito models.ParametroDeposito

	//Getting the parameters
	if err := c.BindJSON(&param_deposito); err != nil {

		fmt.Println("Problem with Json bindng:", err)
		c.AbortWithStatus(http.StatusBadRequest)
	} else {
		// Call the GetCliente function from the models package with the parameters
		estado, err := database.PostDeposito(param_deposito)
		if err != nil {
			c.JSON(404, gin.H{
				"estado": err,
			})
			// Agregar demas estados
		} else {
			c.JSON(http.StatusOK, gin.H{
				"estado": "deposito enviado",
			})
		}
	}
}

func PostTranferencia(c *gin.Context) {
	//Varianble for the client's parameters
	var param_transferencia models.ParametroTransferencia

	//Getting the parameters
	if err := c.BindJSON(&param_transferencia); err != nil {

		fmt.Println("Problem with Json bindng:", err)
		c.AbortWithStatus(http.StatusBadRequest)
	} else {
		// Call the GetCliente function from the models package with the parameters
		estado, err := database.PostTranferencia(param_transferencia)
		if err != nil {
			c.JSON(404, gin.H{
				"estado": err,
			})

			// Agregar demas estados
		} else {
			c.JSON(http.StatusOK, gin.H{
				"estado": "transferencia_enviada",
			})
		}
	}
}

func PostGiro(c *gin.Context) {
	//Varianble for the client's parameters
	var param_giro models.ParametroGiro

	//Getting the parameters
	if err := c.BindJSON(&param_giro); err != nil {

		fmt.Println("Problem with Json bindng:", err)
		c.AbortWithStatus(http.StatusBadRequest)
	} else {
		// Call the GetCliente function from the models package with the parameters
		estado, err := database.PostGiro(param_giro)
		if err != nil {
			c.JSON(404, gin.H{
				"estado": err,
			})

			// Agregar demas estados
		} else {
			c.JSON(http.StatusOK, gin.H{
				"estado": "giro_enviado",
			})
		}
	}
}
