package models

import (
	_ "fmt"
	_ "time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Cliente struct {
	Id                   primitive.ObjectID `bson:"_id" json:"id"`
	Nombre               string             `bson:"nombre" json:"nombre"`
	Contrasena           string             `bson:"contrasena" json:"contrasena"`
	FechaNacimiento      string             `bson:"fecha_nacimiento" json:"fecha_nacimiento"`
	Direccion            string             `bson:"direccion" json:"direccion"`
	NumeroIdentificacion string             `bson:"numero_identificacion" json:"numero_identificacion"`
	Email                string             `bson:"email" json:"email"`
	Telefono             string             `bson:"telefono" json:"telefono"`
	Genero               string             `bson:"genero" json:"genero"`
	Nacionalidad         string             `bson:"nacionalidad" json:"nacionalidad"`
	Ocupacion            string             `bson:"ocupacion" json:"ocupacion"`
}

type Billetera struct {
	Id         primitive.ObjectID `bson:"_id" json:"id"`
	NroCliente string             `bson:"nro_cliente" json:"nro_cliente"`
	Saldo      float64            `bson:"saldo" json:"saldo"`
	Divisa     string             `bson:"divisa" json:"divisa"`
	Nombre     string             `bson:"nombre" json:"nombre"`
	Activo     string             `bson:"activo" json:"activo"`
}

type Movimiento struct {
	Id                primitive.ObjectID `bson:"_id" json:"id"`
	NroCliente        string             `bson:"nro_cliente" json:"nro_cliente"`
	NroClienteOrigen  string             `json:"nro_cliente_origen"`
	NroClienteDestino string             `json:"nro_cliente_destino"`
	Monto             float64            `bson:"monto" json:"monto"`
	Divisa            string             `bson:"divisa" json:"divisa"`
	Tipo              string             `bson:"tipo" json:"tipo"`
	Fecha_hora        string             `bson:"fecha_hora" json:"fecha_hora"`
	IdBilletera       primitive.ObjectID `bson:"id_billetera" json:"id_billetera"`
}

type ParametroCliente struct {
	NumeroIdentificacion string `bson:"numero_identificacion" json:"numero_identificacion"`
}

type ParametroInicio struct {
	NumeroIdentificacion string `bson:"numero_identificacion" json:"numero_identificacion"`
	Contrasena           string `bson:"contrasena" json:"contrasena"`
}

type ParametroDeposito struct {
	NroCliente string `bson:"nro_cliente" json:"nro_cliente"`
	Monto      float64 `bson:"monto" json:"monto"`
	Divisa     string `bson:"divisa" json:"divisa"`
}

type ParametroTransferencia struct {
	NroClienteOrigen  string `bson:"nro_cliente_origen" json:"nro_cliente_origen"`
	NroClienteDestino string `bson:"nro_cliente_destino" json:"nro_cliente_destino"`
	Monto             float64 `bson:"monto" json:"monto"`
	Divisa            string `bson:"divisa" json:"divisa"`
}

type ParametroGiro struct {
	NroCliente string `bson:"nro_cliente" json:"nro_cliente"`
	Monto      float64 `bson:"monto" json:"monto"`
	Divisa     string `bson:"divisa" json:"divisa"`
}

type Response struct {
	Estado string `json:"estado"`
}
