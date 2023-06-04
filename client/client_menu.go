package main

import (
	"TrustBankApi/client/utils"
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

func log_in(URL string){
	var (
		numero_de_identificacion string
		contrasena string
	)
	fmt.Print("Ingrese su numero de identificacion:")
	fmt.Scan(&numero_de_identificacion)
	fmt.Print("Ingrese su contrasena:")
	fmt.Scan(&contrasena)
	
	resp, err := utils.InicioSesion(URL, numero_de_identificacion, contrasena)
	if err != nil {
		fmt.Println("error: ", err)
	}
	if resp.Estado == "exitoso" {
		fmt.Println("Login exitoso!")
	} else {
		fmt.Println("Login no exitoso. Intente de nuevo")
		log_in(URL)
	}
}



func actions(URL string){
	var(
		money int
		destination string
	)
	
	fmt.Println("1. Realizar Deposito")
	fmt.Println("2. Realizar Transferancia")
	fmt.Println("3. Realizar Giro")
	fmt.Println("4. Salir")
	fmt.Print("Ingrese a un opcion: ")
	fmt.Scan(&option)

	switch option{
	case 1:
	
	case 2:

	case 3: 
		
	case 4:
		fmt.Println("Gracias por usar TrustBank")
	
	default:
		fmt.Println("Opcion no reconocida")
		fmt.Println("")
		gretting()
	}

}

func main(){
	err := godotenv.Load()
        if err != nil {
                log.Fatal("Error loading .env file")
        }
        var (
                SERVER = os.Getenv("SERVER")
                PORT   = os.Getenv("PORT")
                URL = "http://" + SERVER + ":" + PORT + "/api"
			)
	
	gretting()
	log_in(URL)
	actions(URL)

}



