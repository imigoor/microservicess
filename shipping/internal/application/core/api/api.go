package api

import (
	"log"

	"github.com/ruandg/microservices/shipping/internal/application/core/domain"
	"github.com/ruandg/microservices/shipping/internal/ports"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Application struct {
	db ports.DBPort
}

func NewApplication(db ports.DBPort) *Application {
	return &Application{db: db}
}

func (a Application) CreateShipping(shipping domain.Shipping) (domain.Shipping, error) {
	log.Printf("Processando envio para o OrderID: %d", shipping.OrderID)
	
	var totalItems int32 = 0
	for _, item := range shipping.Items {
		totalItems += item.Quantity
	}
	
	log.Printf("Calculado: %d dias de entrega para %d itens totais", shipping.DeliveryDays, totalItems)

	err := a.db.Save(&shipping)
	if err != nil {
		log.Printf("Erro ao salvar shipping no banco: %v", err)
		return domain.Shipping{}, status.Errorf(codes.Internal, "failed to save shipping. %v", err)
	}
	
	log.Printf("Envio salvo com sucesso no banco de dados!")
	return shipping, nil
}
