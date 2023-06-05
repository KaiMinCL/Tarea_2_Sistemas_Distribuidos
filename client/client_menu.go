package main

import (
	"TrustBankClient/utils"
	"fmt"
	"log"
	"os"
	"github.com/joho/godotenv"
)

func gretting(){
	var option int
	fmt.Println("Bienvenido a TrustBank! ")
	fmt.Println("1. Iniciar Sesion")
	fmt.Println("2. Salir")
	fmt.Print("Ingrese a un opcion: ")
	fmt.Scan(&option)

	switch option{
	case 1:
		
	case 2:
		fmt.Println("Gracias por usar TrustBank")
		os.Exit(0)
	default:
		fmt.Println("Opcion no reconocida")
		fmt.Println("")
		gretting()
	}
}

func log_in(URL string) (string){
	var (
		numero_de_identificacion string
		contrasena string
	)
	fmt.Print("Ingrese su numero de identificacion: ")
	fmt.Scan(&numero_de_identificacion)
	fmt.Print("Ingrese su contrasena: ")
	fmt.Scan(&contrasena)
	
	resp, err := utils.InicioSesion(URL, numero_de_identificacion, contrasena)
	if err != nil {
		fmt.Println("error: ", err)
		log.Fatal("Bug in InicioSesion")
	}
	if resp.Estado == "exitoso" {
		fmt.Println("Login exitoso!")
		return numero_de_identificacion
	} else {
		fmt.Println("Login no exitoso.")
		os.Exit(1)
	}
	return ""
}



func main(){
	err := godotenv.Load("../.env")
        if err != nil {
                log.Fatal("Error loading .env file")
        }
    var (
		option int
		numero_cliente string
		SERVER = os.Getenv("SERVER")
        PORT   = os.Getenv("PORT")
        URL = "http://" + SERVER + ":" + PORT + "/api"
		)
	
	START: gretting()
	numero_cliente = log_in(URL)
	
	fmt.Println("1. Realizar Deposito")
	fmt.Println("2. Realizar Transferancia")
	fmt.Println("3. Realizar Giro")
	fmt.Println("4. Salir")
	fmt.Print("Ingrese a un opcion: ")
	fmt.Scan(&option)

	switch option{
	case 1:
		var amount float64
		fmt.Print("Ingrese un monto: ")
		fmt.Scan(&amount)
		resp, err := utils.Deposito(URL, "USD", numero_cliente, amount)
		if err != nil{
			fmt.Println("error: ", err)
			log.Fatal("Deposito problema")
		}
		if resp.Estado == "deposito_enviado"{
			fmt.Println("El deposito ha sido enviado correctamente.")
		} else if resp.Estado == "cliente_no_encontrado" {
			fmt.Println("Error: cliente no existe")
			os.Exit(1)
		} else if resp.Estado == "billetera_no_encontrado"{
			fmt.Println("Error: billetera no existe")
			os.Exit(1)
		} else{
			fmt.Println("El estado no es bueno deposito")
			os.Exit(1)
		}

	case 2:
		var amount float64
		var destino string
		fmt.Print("Ingrese un monto: ")
		fmt.Scan(&amount)
		fmt.Print("Ingrese un destino: ")
		fmt.Scan(&destino)
		resp, err := utils.Transferencia(URL, numero_cliente, destino, "USD", amount)
		if err != nil{
			fmt.Println("error: ", err)
			log.Fatal("Transferencia problema")
		}
		if resp.Estado == "transferencia_eviada"{
			fmt.Println("La transferencia fue enviada correctamente")
		} else if resp.Estado == "cliente_origen_no_encontrado" {
			fmt.Println("Error: cliente origen no existe")
			os.Exit(1)
		} else if resp.Estado == "cliente_destino_no_encontrado"{
			fmt.Println("Error: cliente destino no existe")
			os.Exit(1)
		} else if resp.Estado == "billetera_origen_no_encontrado"{
			fmt.Println("Error: billetera origen no existe")
			os.Exit(1)
		} else if resp.Estado == "billetera_destino_no_encontrado"{
			fmt.Println("Error: billetera destino no existe")
			os.Exit(1)
		} else if resp.Estado == "billetera_sin_fondos_suficientes"{
			fmt.Println("Error: billetera sin fondos suficientes")
			os.Exit(1)
		} else{
			fmt.Println("El estado no es bueno transferencia")
			os.Exit(1)
		}

	case 3:
		var amount float64
		fmt.Print("Ingrese un monto: ")
		fmt.Scan(&amount)
		resp, err := utils.Giro(URL, numero_cliente, "USD", amount)
		if err != nil{
			fmt.Println("error: ", err)
			log.Fatal("Giro problema")
		}
		if resp.Estado == "giro_eviado"{
			fmt.Println("La transferencia fue enviada correctamente")
		} else if resp.Estado == "cliente_no_encontrado" {
			fmt.Println("Error: cliente no existe")
			os.Exit(1)
		} else if resp.Estado == "billetera_no_encontrado"{
			fmt.Println("Error:  billetra no existe")
			os.Exit(1)
		} else if resp.Estado == "billetera_sin_fondos_suficientes"{
			fmt.Println("Error: billetera sin fondos suficientes")
		} else{
			fmt.Println("El estado no es bueno transferencia")
			os.Exit(1)
		}

	case 4:
		goto START
	default:
		fmt.Println("Opcion no reconocida")
		fmt.Println("")
		gretting()
	}

}



