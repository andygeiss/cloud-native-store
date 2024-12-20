// This is the main package for initializing and running the server.
package main

import (
	"log"
	"os"

	"github.com/andygeiss/cloud-native-store/internal/app/adapters/common/config"
	"github.com/andygeiss/cloud-native-store/internal/app/adapters/inbound/api"
	"github.com/andygeiss/cloud-native-store/internal/app/adapters/outbound/inmemory"
	"github.com/andygeiss/cloud-native-store/internal/app/core/services"
	"github.com/andygeiss/cloud-native-utils/consistency"
	"github.com/andygeiss/cloud-native-utils/security"
)

func main() {

	// Create a new configuration object, specifying paths to the TLS certificate
	// and key files needed for secure communication.
	cfg := &config.Config{
		Key:    security.Getenv("ENCRYPTION_KEY"),
		Server: config.Server{CertFile: os.Getenv("SERVER_CERTIFICATE"), KeyFile: os.Getenv("SERVER_KEY")},
	}

	// Initialize a JSON file logger to log transactional data.
	logger := consistency.NewJsonFileLogger[string, string](os.Getenv("TRANSACTIONAL_LOG"))

	// Create a new object service and configure it with the transactional logger and the in-memory port.
	port := inmemory.NewObjectStore(1)
	service := services.
		NewObjectService().
		WithTransactionalLogger(logger).
		WithPort(port)

	// Add the service instance to the configuration.
	cfg.Services = config.Services{
		ObjectService: service,
	}

	// Set up the service. If an error occurs during setup, log it and terminate the program.
	if err := service.Setup(); err != nil {
		log.Fatalf("error during setup: %v", err)
	}
	// Ensure proper cleanup of resources when the service is no longer needed.
	defer service.Teardown()

	// Initialize the API router using the configuration object.
	mux := api.Route(cfg)

	// Create a new secure server instance, binding it to the localhost address.
	srv := security.NewServer(mux, "localhost")
	// Ensure the server is properly closed when the program exits.
	defer srv.Close()

	// Start the server using TLS for secure communication, providing the certificate
	// and key files specified in the configuration. Log an error if server startup fails.
	log.Printf("start listening...")
	if err := srv.ListenAndServeTLS(cfg.Server.CertFile, cfg.Server.KeyFile); err != nil {
		log.Fatalf("listening failed: %v", err)
	}
}
