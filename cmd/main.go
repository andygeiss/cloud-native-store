// This is the main package for initializing and running the server.
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/andygeiss/cloud-native-store/internal/app/adapters/inbound/api"
	"github.com/andygeiss/cloud-native-store/internal/app/adapters/outbound/inmemory"
	"github.com/andygeiss/cloud-native-store/internal/app/config"
	"github.com/andygeiss/cloud-native-store/internal/app/core/services"
	"github.com/andygeiss/cloud-native-utils/security"
)

func main() {
	// Create a new configuration object.
	cfg := &config.Config{
		Key: security.Getenv("ENCRYPTION_KEY"),
	}

	// Create a new object service and configure it with the transactional logger and the in-memory port.
	port := inmemory.NewObjectStore(1)
	service := services.
		NewObjectService(cfg).
		// WithTransactionalLogger(logger).
		WithPort(port)

	// Set up the service. If an error occurs during setup, log it and terminate the program.
	if err := service.Setup(); err != nil {
		log.Fatalf("error during setup: %v", err)
	}
	// Ensure proper cleanup of resources when the service is no longer needed.
	defer service.Teardown()

	// Initialize the API router using the configuration object.
	mux := api.Route(service)

	// Start the HTTP server.
	log.Printf("start listening...")
	if err := http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("PORT")), mux); err != nil {
		log.Fatalf("listening failed: %v", err)
	}
}
