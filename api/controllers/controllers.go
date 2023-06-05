package controllers

import (
	"TrustBankAPI/apiDB"
	"common/database"
	"common/models"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/streadway/amqp"
)

var rabbitConn *amqp.Connection
var rabbitCh *amqp.Channel
var rabbitQueue string

func GetClient(c *gin.Context) {

	//Varianble for the client's parameters
	var param_cliente models.ParametroCliente

	//Getting the parameters
	if err := c.BindJSON(&param_cliente); err != nil {

		fmt.Println("Problem with JSON bindng:", err)
		c.AbortWithStatus(http.StatusBadRequest)
	} else {
		// Call the GetCliente function from the models package with the parameters
		cliente, err := apiDB.GetClient(param_cliente)

		if err != nil {
			c.JSON(http.StatusNotFound, models.Response{Estado: "cliente_no_encontrado"})
			return
		}
		// Crear la respuesta del cliente sin la contraseña
		respuesta := models.Cliente{
			Id:                   cliente.Id,
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
		session := apiDB.VerifySession(param_inicio)

		if session {
			c.JSON(http.StatusOK, models.Response{Estado: "exitoso"})
		} else {
			c.JSON(http.StatusUnauthorized, models.Response{Estado: "no_exitoso"})
		}
	}
}

func InitRabbitMQ() error {
	var rabbitPort = os.Getenv("RABBITMQ_PORT")
	var rabbitHost = os.Getenv("RABBITMQ_HOST")
	var rabbitUsername = os.Getenv("RABBITMQ_USERNAME")
	var rabbitPassword = os.Getenv("RABBITMQ_PASSWORD")
	rabbitQueue = os.Getenv("RABBITMQ_QUEUE_NAME")

	// Establish a connection to RabbitMQ
	conn, err := amqp.Dial("amqp://" + rabbitUsername + ":" + rabbitPassword + "@" + rabbitHost + ":" + rabbitPort + "/")
	if err != nil {
		return err
	}
	rabbitConn = conn

	// Create a channel
	ch, err := rabbitConn.Channel()
	if err != nil {
		return err
	}
	rabbitCh = ch

	// Declare the queue
	_, err = rabbitCh.QueueDeclare(
		rabbitQueue, // Name of the queue
		false,       // Durable (false for non-durable)
		false,       // AutoDelete (false for non-auto-delete)
		false,       // Exclusive (false for non-exclusive)
		false,       // NoWait (false to wait for server response)
		nil,         // Arguments
	)
	if err != nil {
		return err
	}

	return nil
}

func CloseRabbitMQ() {
	if rabbitCh != nil {
		rabbitCh.Close()
	}
	if rabbitConn != nil {
		rabbitConn.Close()
	}
}

// EnviarMensajeRabbitMQ envía un mensaje a RabbitMQ
func SendRabbitMessage(movimiento models.Movimiento) error {

	body, err := json.Marshal(movimiento)
	if err != nil {
		return err
	}

	err = rabbitCh.Publish(
		"",          // Exchange
		rabbitQueue, // Routing key
		false,       // Mandatory
		false,       // Immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		return err
	}

	return nil
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

		_, err := apiDB.GetClient(param_cliente)
		if err != nil {
			c.JSON(http.StatusNotFound, models.Response{Estado: "cliente_no_encontrado"})
			return
		}

		// Obtener la billetera del cliente
		_, err = apiDB.GetWallet(param_deposito.NroCliente, param_deposito.Divisa)
		if err != nil {
			c.JSON(http.StatusNotFound, models.Response{Estado: "billetera_no_encontrada"})
			return
		}

		monto, _ := strconv.ParseFloat(param_deposito.Monto, 64)

		// Realizar el depósito enviando un mensaje a RabbitMQ
		movimiento := models.Movimiento{
			NroClienteOrigen:  param_deposito.NroCliente,
			NroClienteDestino: param_deposito.NroCliente,
			Monto:             monto,
			Divisa:            param_deposito.Divisa,
			Tipo:              "deposito",
		}

		err = SendRabbitMessage(movimiento)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.Response{Estado: "error_deposito"})
			return
		}

		c.JSON(http.StatusOK, models.Response{Estado: "deposito_enviado"})
	}
}

// VerificarFondosSuficientes verifica si la billetera de origen tiene fondos suficientes para la transferencia
func VerifyFunds(billetera models.Billetera, monto string) bool {
	// Convertir el saldo y el monto a números decimales

	montoTransferencia, err := strconv.ParseFloat(monto, 64)
	if err != nil {
		log.Println("Error al convertir el monto de transferencia a número decimal:", err)
		return false
	}

	// Verificar si el saldo es suficiente para la transferencia
	if billetera.Saldo >= montoTransferencia {
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

		_, err := apiDB.GetClient(param_cliente)
		if err != nil {
			c.JSON(http.StatusNotFound, models.Response{Estado: "cliente_origen_no_encontrado"})
			return
		}

		//Verificar si existe el cliente destino
		param_cliente.NumeroIdentificacion = param_transferencia.NroClienteDestino

		_, err = apiDB.GetClient(param_cliente)
		if err != nil {
			c.JSON(http.StatusNotFound, models.Response{Estado: "cliente_destino_no_encontrado"})
			return
		}

		// Obtener la billetera del cliente de origen
		billeteraOrigen, err := database.GetWallet(param_transferencia.NroClienteOrigen, param_transferencia.Divisa)
		if err != nil {
			c.JSON(http.StatusNotFound, models.Response{Estado: "billetera_origen_no_encontrada"})
			return
		}

		// Obtener la billetera del cliente de destino
		_, err = database.GetWallet(param_transferencia.NroClienteDestino, param_transferencia.Divisa)
		if err != nil {
			c.JSON(http.StatusNotFound, models.Response{Estado: "billetera_destino_no_encontrada"})
			return
		}

		// Verificar si la billetera de origen tiene fondos suficientes
		if !VerifyFunds(billeteraOrigen, param_transferencia.Monto) {
			c.JSON(http.StatusUnprocessableEntity, models.Response{Estado: "billetera_origen_sin_fondos_suficientes"})
			return
		}

		montoTransferencia, err := strconv.ParseFloat(param_transferencia.Monto, 64)
		if err != nil {
			log.Println("Error al convertir el monto de transferencia a número decimal:", err)
			return
		}

		// Realizar la transferencia enviando un mensaje a RabbitMQ
		movimiento := models.Movimiento{
			NroClienteOrigen:  param_transferencia.NroClienteOrigen,
			NroClienteDestino: param_transferencia.NroClienteDestino,
			Monto:             montoTransferencia,
			Divisa:            param_transferencia.Divisa,
			Tipo:              "transferencia",
		}

		err = SendRabbitMessage(movimiento)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.Response{Estado: "error_transferencia"})
			return
		}

		c.JSON(http.StatusOK, models.Response{Estado: "transferencia_enviada"})
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

		_, err := apiDB.GetClient(param_cliente)
		if err != nil {
			c.JSON(http.StatusNotFound, models.Response{Estado: "cliente_no_encontrado"})
			return
		}

		// Obtener la billetera del cliente
		billetera, err := database.GetWallet(param_giro.NroCliente, param_giro.Divisa)
		if err != nil {
			c.JSON(http.StatusNotFound, models.Response{Estado: "billetera_no_encontrada"})
			return
		}

		// Verificar si la billetera de origen tiene fondos suficientes
		if !VerifyFunds(billetera, param_giro.Monto) {
			c.JSON(http.StatusUnprocessableEntity, models.Response{Estado: "billetera_sin_fondos_suficientes"})
			return
		}

		montoGiro, err := strconv.ParseFloat(param_giro.Monto, 64)

		// Create a message payload
		movimiento := models.Movimiento{
			NroClienteOrigen:  param_giro.NroCliente,
			NroClienteDestino: param_giro.NroCliente,
			Monto:             montoGiro,
			Divisa:            param_giro.Divisa,
			Tipo:              "giro",
		}

		err = SendRabbitMessage(movimiento)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.Response{Estado: "error_transferencia"})
			return
		}

		c.JSON(http.StatusOK, models.Response{Estado: "giro_enviado"})
	}
}
