package controllers

import (
	"TrustBankApi/database"
	"TrustBankApi/models"
	"fmt"
	"log"
	_ "log"
	_ "math/rand"
	"net/http"
	"strconv"
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
			c.JSON(http.StatusNotFound, models.Response{Estado: "cliente_no_encontrado"})
			return
		}
		// Crear la respuesta del cliente sin la contraseña
		respuesta := models.Cliente{
			ID:                   cliente.ID,
			Nombre:               cliente.Nombre,
			FechaNacimiento:      cliente.FechaNacimiento,
			Direccion:            cliente.Direccion,
			NumeroIdentificacion: cliente.NumeroIdentificacion,
			Email:                cliente.Email,
			Telefono:             cliente.Telefono,
			Genero:               cliente.Genero,
			Nacionalidad:         cliente.Nacionalidad,
			Ocupacion:            cliente.Ocupacion,
		}

		c.JSON(http.StatusOK, respuesta)
	}
}

func SessionHandler(c *gin.Context) {
	//Varianble for the client's parameters
	var param_inicio models.ParametroInicio

	//Getting the parameters
	if err := c.BindJSON(&param_inicio); err != nil {

		fmt.Println("Problem with Json bindng:", err)
		c.AbortWithStatus(http.StatusBadRequest)
	} else {
		// Call the GetCliente function from the models package with the parameters
		session := database.VerifySession(param_inicio)

		if session {
			c.JSON(http.StatusOK, models.Response{Estado: "exitoso"})
		} else {
			c.JSON(http.StatusUnauthorized, models.Response{Estado: "no_exitoso"})
		}
	}
}

func DepositHandler(c *gin.Context) {
	//Varianble for the client's parameters
	var param_deposito models.ParametroDeposito

	//Getting the parameters
	if err := c.BindJSON(&param_deposito); err != nil {

		fmt.Println("Problem with Json bindng:", err)
		c.AbortWithStatus(http.StatusBadRequest)
	} else {
		// Verificar si el cliente existe

		var param_cliente models.ParametroCliente
		param_cliente.NumeroIdentificacion = param_deposito.NroCliente

		cliente, err := database.GetClient(param_cliente)
		if err != nil {
			c.JSON(http.StatusNotFound, models.Estado{Estado: "cliente_no_encontrado"})
			return
		}

		// Obtener la billetera del cliente
		billetera, err := database.GetWallet(param_deposito.NumeroCliente)
		if err != nil {
			c.JSON(http.StatusNotFound, models.Estado{Estado: "billetera_no_encontrada"})
			return
		}

		// Realizar el depósito enviando un mensaje a RabbitMQ
		//
		//
		// 	      TO DO
		//
		//
		//

		c.JSON(http.StatusOK, models.Estado{Estado: "deposito_enviado"})
	}
}

// VerificarFondosSuficientes verifica si la billetera de origen tiene fondos suficientes para la transferencia
func VerifyFunds(billetera models.Billetera, monto string) bool {
	// Convertir el saldo y el monto a números decimales
	saldo, err := strconv.ParseFloat(billetera.Saldo, 64)
	if err != nil {
		log.Println("Error al convertir el saldo de la billetera a número decimal:", err)
		return false
	}

	montoTransferencia, err := strconv.ParseFloat(monto, 64)
	if err != nil {
		log.Println("Error al convertir el monto de transferencia a número decimal:", err)
		return false
	}

	// Verificar si el saldo es suficiente para la transferencia
	if saldo >= montoTransferencia {
		return true
	}

	return false
}

func TransferHandler(c *gin.Context) {
	//Varianble for the client's parameters
	var param_transferencia models.ParametroTransferencia

	//Getting the parameters
	if err := c.BindJSON(&param_transferencia); err != nil {

		fmt.Println("Problem with Json bindng:", err)
		c.AbortWithStatus(http.StatusBadRequest)
	} else {

		//Verificar si existe el cliente Origen
		var param_cliente models.ParametroCliente
		param_cliente.NumeroIdentificacion = param_transferencia.NroClienteOrigen

		cliente, err := database.GetClient(param_cliente)
		if err != nil {
			c.JSON(http.StatusNotFound, models.Estado{Estado: "cliente_no_encontrado"})
			return
		}

		//Verificar si existe el cliente destino
		param_cliente.NumeroIdentificacion = param_transferencia.NroClienteDestino

		cliente, err = database.GetClient(param_cliente)
		if err != nil {
			c.JSON(http.StatusNotFound, models.Estado{Estado: "cliente_no_encontrado"})
			return
		}

		// Obtener la billetera del cliente de origen
		billeteraOrigen, err := models.GetWallet(param_transferencia.NroClienteOrigen)
		if err != nil {
			c.JSON(http.StatusNotFound, models.Estado{Estado: "billetera_destino_no_encontrada"})
			return
		}

		// Obtener la billetera del cliente de destino
		billeteraDestino, err := models.GetWallet(param_transferencia.NroClienteDestino)
		if err != nil {
			c.JSON(http.StatusNotFound, models.Estado{Estado: "billetera_destino_no_encontrada"})
			return
		}

		// Verificar si la billetera de origen tiene fondos suficientes
		if !VerifyFunds(billeteraOrigen, param_transferencia.Monto) {
			c.JSON(http.StatusUnprocessableEntity, models.Estado{Estado: "billetera_origen_sin_fondos_suficientes"})
			return
		}

		// Realizar la transferencia enviando un mensaje a RabbitMQ
		//
		//
		// 	      TO DO
		//
		//
		//

		c.JSON(http.StatusOK, models.Estado{Estado: "transferencia_enviada"})
	}
}

func WithdrawHandler(c *gin.Context) {
	//Varianble for the client's parameters
	var param_giro models.ParametroGiro

	//Getting the parameters
	if err := c.BindJSON(&param_giro); err != nil {

		fmt.Println("Problem with Json bindng:", err)
		c.AbortWithStatus(http.StatusBadRequest)
	} else {
		//Verificar si existe el cliente
		var param_cliente models.ParametroCliente
		param_cliente.NumeroIdentificacion = param_giro.NroCliente

		cliente, err := database.GetClient(param_cliente)
		if err != nil {
			c.JSON(http.StatusNotFound, models.Estado{Estado: "cliente_no_encontrado"})
			return
		}

		// Obtener la billetera del cliente
		billetera, err := models.GetWallet(param_giro.NroCliente)
		if err != nil {
			c.JSON(http.StatusNotFound, models.Estado{Estado: "billetera_destino_no_encontrada"})
			return
		}

		// Verificar si la billetera de origen tiene fondos suficientes
		if !VerifyFunds(billetera, param_giro.Monto) {
			c.JSON(http.StatusUnprocessableEntity, models.Estado{Estado: "billetera_origen_sin_fondos_suficientes"})
			return
		}

		// Realizar el giro enviando un mensaje a RabbitMQ
		//
		//
		// 	      TO DO
		//
		//
		//

		c.JSON(http.StatusOK, models.Estado{Estado: "giro_enviado"})
	}
}
