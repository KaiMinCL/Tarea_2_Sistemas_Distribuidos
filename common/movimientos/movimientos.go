package movimientos

import (
	"context"
	"log"
)

type Server struct {
}

func (s *Server) RegistrarMovimiento(ctx context.Context, movimientoRequest *MovimientoRequest) (*MovimientoRequest, error) {
	log.Printf("Tipo de operación del movimiento recibido desde el cliente: %s", movimientoRequest.TipoOperacion)
	return &MovimientoRequest{TipoOperacion: "Confirmacion"}, nil
}
