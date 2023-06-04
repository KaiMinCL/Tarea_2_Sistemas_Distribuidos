package utils

import (
	"fmt"
	"net/http"
	"log"
	"io/ioutil"
	"encoding/json"
	"TrustBankApi/api/models"
	"bytes"
	)

func InicioSesion(URL, id, passwd string) (models.Response, error){
	var (
		param models.ParametroInicio
		respuesta models.Response
	)

	param.NumeroIdentificacion = id 
	param.Contrasena = passwd

	JSONString, err := json.Marshal(param)
	if err != nil {
        return respuesta, err
    }
	
	resp, err := http.Post(URL + "/inicio_sesion", "application/json", bytes.NewBuffer(JSONString))
	if err != nil {
		log.Fatal(err)
		fmt.Println("error: ", err)
        return respuesta, err
    }

	body, err := ioutil.ReadAll(resp.Body)
        if err != nil {
        fmt.Println("error:", err)
        return respuesta, err
    }
	
	err = json.Unmarshal(body, &respuesta)
        if err != nil {
        fmt.Println("error:", err)
        return respuesta, err
    }
	return respuesta, nil
}

func Deposito(URL, divisa, numero_cliente string, amount int) (models.Response, error){
	var (
		param models.ParametroDeposito
		respuesta models.Response
	)

	param.NroCliente = numero_cliente
	param.Divisa = divisa
	param.Monto = amount

	JSONString, err := json.Marshal(param)
	if err != nil {
        return respuesta, err
    }
	
	resp, err := http.Post(URL + "/deposito", "application/json", bytes.NewBuffer(JSONString))
	if err != nil {
		log.Fatal(err)
		fmt.Println("error: ", err)
        return respuesta, err
    }

	body, err := ioutil.ReadAll(resp.Body)
        if err != nil {
        fmt.Println("error:", err)
        return respuesta, err
    }
	
	err = json.Unmarshal(body, &respuesta)
        if err != nil {
        fmt.Println("error:", err)
        return respuesta, err
    }
	return respuesta, nil
}

func Transferencia(URL, origen, destino, divisa string, amount int) (models.Response, error){
	var (
		param models.ParametroTransferencia
		respuesta models.Response
	)

	param.NroClienteOrigen = origen
	param.NroClienteDestino = destino
	param.Divisa = divisa
	param.Monto = amount

	JSONString, err := json.Marshal(param)
	if err != nil {
        return respuesta, err
    }
	
	resp, err := http.Post(URL + "/deposito", "application/json", bytes.NewBuffer(JSONString))
	if err != nil {
		log.Fatal(err)
		fmt.Println("error: ", err)
        return respuesta, err
    }

	body, err := ioutil.ReadAll(resp.Body)
        if err != nil {
        fmt.Println("error:", err)
        return respuesta, err
    }
	
	err = json.Unmarshal(body, &respuesta)
        if err != nil {
        fmt.Println("error:", err)
        return respuesta, err
    }
	return respuesta, nil
}

func Giro(URL, numero_cliente, divisa string, amount int) (models.Response, error){
	var(
		param models.ParametroGiro
		respuesta models.Response
	)
	
	param.NroCliente = numero_cliente
	param.Divisa = divisa
	param.Monto = amount

	JSONString, err := json.Marshal(param)
	if err != nil {
        return respuesta, err
    }
	
	resp, err := http.Post(URL + "/deposito", "application/json", bytes.NewBuffer(JSONString))
	if err != nil {
		log.Fatal(err)
		fmt.Println("error: ", err)
        return respuesta, err
    }

	body, err := ioutil.ReadAll(resp.Body)
        if err != nil {
        fmt.Println("error:", err)
        return respuesta, err
    }
	
	err = json.Unmarshal(body, &respuesta)
        if err != nil {
        fmt.Println("error:", err)
        return respuesta, err
    }
	return respuesta, nil

}
