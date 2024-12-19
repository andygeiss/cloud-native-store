// This is the main package for initializing and running the server.
package main

import (
	"log"

	"github.com/andygeiss/cloud-native-store/internal/app/adapters/common/config"
	"github.com/andygeiss/cloud-native-store/internal/app/adapters/inbound/api"
	"github.com/andygeiss/cloud-native-store/internal/app/adapters/outbound/inmemory"
	"github.com/andygeiss/cloud-native-store/internal/app/core/services"
	"github.com/andygeiss/cloud-native-utils/consistency"
	"github.com/andygeiss/cloud-native-utils/security"
)

func main() {
	// Set up the server configuration, including paths to the TLS certificate and key files.
	cfg := &config.Config[string, any]{
		Server: config.Server{CertFile: ".tls/server.crt", KeyFile: ".tls/server.key"},
	}

	logger := consistency.NewJsonFileLogger[string, any](".cache/transactions.json")

	port := inmemory.NewObjectStore[string, any](1)

	service := services.
		NewObjectService[string, any]().
		WithTransactionalLogger(logger).
		WithPort(port)

	cfg.Services = config.Services[string, any]{
		ObjectService: service,
	}

	if err := service.Setup(); err != nil {
		log.Fatalf("error during setup: %v", err)
	}
	defer service.Teardown()

	// Initialize the API routing with the given configuration.
	mux := api.Route(cfg)

	// Create a new secure server instance, listening on localhost.
	srv := security.NewServer(mux, "localhost")
	defer srv.Close() // Ensure the server is properly closed when the program exits.
	log.Printf("start listening...")

	// Start the server with TLS enabled using the provided certificate and key files.
	if err := srv.ListenAndServeTLS(cfg.Server.CertFile, cfg.Server.KeyFile); err != nil {
		log.Fatalf("listening failed: %v", err)
	}
}
