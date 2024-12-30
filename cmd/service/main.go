// This is the main package for initializing and running the server.
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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

	// Create a new context with a cancel function.
	ctx, cancel := signal.NotifyContext(
		context.Background(),
		// SIGTERM is sent by Kubernetes to gracefully stop a container.
		syscall.SIGTERM,
		// SIGINT is sent by a user terminal to interrupt a running process.
		syscall.SIGINT,
		// SIGQUIT is sent by a user terminal to make a core dump.
		syscall.SIGQUIT,
		// SIGKILL is sent by a user terminal to kill a process immediately.
		syscall.SIGKILL,
	)
	defer cancel()

	// Set up the service. If an error occurs during setup, log it and terminate the program.
	if err := service.Setup(); err != nil {
		log.Fatalf("error during setup: %v", err)
	}
	// Ensure proper cleanup of resources when the service is no longer needed.
	defer service.Teardown()

	// Initialize the API router using the configuration object.
	mux := api.Route(service, ctx)

	// Create a new secure server.
	server := security.NewServer(mux)
	defer server.Close()

	// Register the server shutdown function.
	// This function is called when the server is shut down.
	server.RegisterOnShutdown(func() {
		// TODO: Add cleanup logic here.
	})

	// Run the context check in a separate goroutine.
	go func() {
		<-ctx.Done()
		// Wait for the readiness check to fail.
		<-time.After(5 * time.Second)
		// Shut down the server using a new context.
		// Do not use the original (already closed) context.
		server.Shutdown(context.Background())
	}()

	// Start the HTTP server in the main goroutine.
	log.Printf("start listening at port %s ...", cfg.Server.Port)
	if err := server.ListenAndServe(); err != nil {
		// Check if the server was closed intentionally.
		if err == http.ErrServerClosed {
			log.Println("server is closed.")
			return
		}

		// Log the error and terminate the program.
		log.Fatalf("listening failed: %v", err)
	}
}
