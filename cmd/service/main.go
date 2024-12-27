// This is the main package for initializing and running the server.
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/andygeiss/cloud-native-store/internal/app/adapters/inbound/api"
	"github.com/andygeiss/cloud-native-store/internal/app/adapters/outbound/gcp"
	"github.com/andygeiss/cloud-native-store/internal/app/config"
	"github.com/andygeiss/cloud-native-store/internal/app/core/services"
	"github.com/andygeiss/cloud-native-utils/security"
)

func main() {
	// Create a new configuration object.
	cfg := &config.Config{
		PortCloudSpanner: config.PortCloudSpanner{
			ProjectID:  os.Getenv("GCP_PROJECT_ID"),
			DatabaseID: os.Getenv("GCP_SPANNER_DATABASE_ID"),
			InstanceID: os.Getenv("GCP_SPANNER_INSTANCE_ID"),
			Table:      "KeyValueStore",
		},
		Service: config.Service{
			Key: security.Getenv("ENCRYPTION_KEY"),
		},
		Server: config.Server{
			Port: os.Getenv("PORT"),
		},
	}

	// Create a new Cloud Spanner adapter.
	objectPort := gcp.NewCloudSpanner(cfg)

	// Create a new Object Service.
	service := services.
		NewObjectService(cfg).
		WithPort(objectPort)

	// Set up the service. If an error occurs during setup, log it and terminate the program.
	if err := service.Setup(); err != nil {
		log.Fatalf("error during setup: %v", err)
	}
	// Ensure proper cleanup of resources when the service is no longer needed.
	defer service.Teardown()

	// Initialize the API router using the configuration object.
	mux := api.Route(service)

	// Start the HTTP server.
	log.Printf("start listening at port %s ...", cfg.Server.Port)
	if err := http.ListenAndServe(fmt.Sprintf(":%s", cfg.Server.Port), mux); err != nil {
		log.Fatalf("listening failed: %v", err)
	}
}
