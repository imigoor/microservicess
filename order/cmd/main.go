package main

import (
	"log"

	"github.com/ruandg/microservices/order/config"
	"github.com/ruandg/microservices/order/internal/adapters/db"
	grpcAdapter "github.com/ruandg/microservices/order/internal/adapters/grpc"
	payment_adapter "github.com/ruandg/microservices/order/internal/adapters/payment"
	shipping_adapter "github.com/ruandg/microservices/order/internal/adapters/shipping"
	"github.com/ruandg/microservices/order/internal/application/core/api"
)

func main() {
	log.Println("Starting order service...")
	log.Println("Connecting to database...")
	dbAdapter, err := db.NewAdapter(config.GetDataSourceURL())
	if err != nil {
		log.Fatalf("Failed to connect to database. Error: %v", err)
	}
	log.Println("Database connected!")

	log.Println("Connecting to payment service...")
	paymentAdapter, err := payment_adapter.NewAdapter(config.GetPaymentServiceUrl())
	if err != nil {
		log.Fatalf("Failed to initialize payment stub. Error: %v", err)
	}
	log.Println("Payment service connected!")

	log.Println("Connecting to shipping service...")
	shippingAdapter, err := shipping_adapter.NewAdapter(config.GetShippingServiceUrl())
	if err != nil {
		log.Fatalf("Failed to initialize shipping stub. Error: %v", err)
	}
	log.Println("Shipping service connected!")

	application := api.NewApplication(dbAdapter, paymentAdapter, shippingAdapter)
	adapter := grpcAdapter.NewAdapter(application, config.GetApplicationPort())
	log.Println("Starting gRPC server...")
	adapter.Run()
}
