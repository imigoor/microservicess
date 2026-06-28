package main

import (
	"log"

	"github.com/ruandg/microservices/shipping/config"
	"github.com/ruandg/microservices/shipping/internal/adapters/db"
	grpcAdapter "github.com/ruandg/microservices/shipping/internal/adapters/grpc"
	"github.com/ruandg/microservices/shipping/internal/application/core/api"
)

func main() {
	log.Println("Starting shipping service...")
	dbAdapter, err := db.NewAdapter(config.GetDataSourceURL())
	if err != nil {
		log.Fatalf("Failed to connect to database. Error: %v", err)
	}
	application := api.NewApplication(dbAdapter)
	adapter := grpcAdapter.NewAdapter(application, config.GetApplicationPort())
	adapter.Run()
}
