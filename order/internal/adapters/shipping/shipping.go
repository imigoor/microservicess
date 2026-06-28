package shipping_adapter

import (
	"context"
	"log"
	"time"

	"github.com/ruandg/microservices-proto/golang/shipping"
	"github.com/ruandg/microservices/order/internal/application/core/domain"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

type Adapter struct {
	shipping shipping.ShippingClient
}

func NewAdapter(shippingServiceUrl string) (*Adapter, error) {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	conn, err := grpc.Dial(shippingServiceUrl, opts...)
	if err != nil {
		return nil, err
	}
	client := shipping.NewShippingClient(conn)
	return &Adapter{shipping: client}, nil
}

func (a *Adapter) Ship(order *domain.Order) (int32, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	var items []*shipping.ShippingItem
	for _, item := range order.OrderItems {
		items = append(items, &shipping.ShippingItem{
			ProductCode: item.ProductCode,
			Quantity:    item.Quantity,
		})
	}

	resp, err := a.shipping.Create(ctx, &shipping.CreateShippingRequest{
		OrderId:       order.ID,
		ShippingItems: items,
	})
	if err != nil {
		if status.Code(err) == codes.DeadlineExceeded {
			log.Printf("shipping timeout: deadline exceeded for order %d", order.ID)
		}
		log.Printf("shipping error: %v", err)
		return 0, err
	}
	return resp.DeliveryDays, nil
}
