package controllers


import (
	"banco/models"
	"banco/database"
    "fmt"
    _ "log"
    _ "math/rand"
	"net/http"
    _ "strings"
	_ "time"
	"github.com/gin-gonic/gin"
)

func GetCliente(c *gin.Context){

	//Varianble for the client's parameters
	var param_cliente models.ParametroCliente

	//Getting the parameters
	if err := c.BindJSON(&param_cliente); err != nil {
	
		fmt.Println("Problem with Json bindng:", err)
		c.AbortWithStatus(http.StatusBadRequest)
	} else {
        // Call the GetCliente function from the models package with the parameters
		cliente, err := database.GetCliente(param_cliente)
		if err != nil{
			c.AbortWithStatus(http.StatusBadRequest)
		} else {
			c.JSON(http.StatusOK, cliente)
		}
    }
}

