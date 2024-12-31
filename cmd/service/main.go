// This is the main package for initializing and running the server.
package main

import (
	"context"
	"embed"
	"log"
	"net/http"
	"os"

	"github.com/andygeiss/cloud-native-store/internal/app/adapters/inbound/api"
	"github.com/andygeiss/cloud-native-store/internal/app/adapters/outbound/gcp"
	"github.com/andygeiss/cloud-native-store/internal/app/config"
	"github.com/andygeiss/cloud-native-store/internal/app/core/services"
	"github.com/andygeiss/cloud-native-utils/security"
	"github.com/andygeiss/cloud-native-utils/service"
)

//go:embed assets
var efs embed.FS

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
			Efs:  efs,
			Port: os.Getenv("PORT"),
		},
	}

	// Create a new Cloud Spanner adapter.
	objectPort := gcp.NewCloudSpanner(cfg)

	// Create a new Object Service.
	svc := services.
		NewObjectService(cfg).
		WithPort(objectPort)

	// Create a new context with a cancel function.
	ctx, cancel := service.Context()
	defer cancel()

	// Set up the service. If an error occurs during setup, log it and terminate the program.
	if err := svc.Setup(); err != nil {
		log.Fatalf("error during setup: %v", err)
	}
	defer svc.Teardown()

	// Initialize the API router using the configuration object.
	mux := api.Route(svc, ctx, cfg)

	// Create a new secure server.
	srv := security.NewServer(mux)
	defer srv.Close()

	// Register the server shutdown function.
	srv.RegisterOnShutdown(func() {
		// TODO: Add cleanup logic here.
	})

	// Register the service shutdown function.
	service.RegisterOnContextDone(ctx, func() {
		// Avoid using the already closed context to prevent a panic.
		// Thus we use a new context to shut down the server.
		srv.Shutdown(context.Background())
	})

	// Start the HTTP server in the main goroutine.
	log.Printf("start listening at port %s ...", cfg.Server.Port)
	if err := srv.ListenAndServe(); err != nil {
		// Check if the server was closed intentionally.
		if err == http.ErrServerClosed {
			log.Println("server is closed.")
			return
		}

		// Log the error and terminate the program.
		log.Fatalf("listening failed: %v", err)
	}
}
