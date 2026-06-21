package main

import (
	"log"

	"github.com/ruandg/microservices/order/config"
	"github.com/ruandg/microservices/order/internal/adapters/db"
	payment_adapter "github.com/ruandg/microservices/order/internal/adapters/payment"
	"github.com/ruandg/microservices/order/internal/adapters/grpc"
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
	application := api.NewApplication(dbAdapter, paymentAdapter)
	grpcAdapter := grpc.NewAdapter(application, config.GetApplicationPort())
	log.Println("Starting gRPC server...")
	grpcAdapter.Run()
}