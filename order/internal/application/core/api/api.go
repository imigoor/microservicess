package api

import (
	"github.com/ruandg/microservices/order/internal/application/core/domain"
	"github.com/ruandg/microservices/order/internal/ports"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Application struct {
	db      ports.DBPort
	payment ports.PaymentPort
}

func NewApplication(db ports.DBPort, payment ports.PaymentPort) *Application {
	return &Application{
		db:      db,
		payment: payment,
	}
}

func (a Application) PlaceOrder(order domain.Order) (domain.Order, error) {
	var totalQuantity int32
	for _, item := range order.OrderItems {
		totalQuantity += item.Quantity
	}
	if totalQuantity > 50 {
		return domain.Order{}, status.Errorf(codes.InvalidArgument, "Order with more than 50 items is not allowed.")
	}

	err := a.db.Save(&order)
	if err != nil {
		order.Status = "Canceled"
		a.db.UpdateStatus(&order)
		return domain.Order{}, status.Errorf(codes.Internal, "failed to save order. %v", err)
	}

	paymentErr := a.payment.Charge(&order)
	if paymentErr != nil {
		order.Status = "Canceled"
		a.db.UpdateStatus(&order)
		return domain.Order{}, paymentErr
	}

	order.Status = "Paid"
	a.db.UpdateStatus(&order)
	return order, nil
}
