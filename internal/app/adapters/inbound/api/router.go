package api

import (
	"context"
	"net/http"

	"github.com/andygeiss/cloud-native-store/internal/app/core/services"
	"github.com/andygeiss/cloud-native-utils/security"
)

// Route creates a new mux with the health check endpoint (/health)
// and the store endpoints (/api/v1/store).
func Route(service *services.ObjectService, ctx context.Context) *http.ServeMux {
	// Create a new mux with health check endpoint.
	mux := security.Mux(ctx)
	mux.HandleFunc("DELETE /api/v1/store", Delete(service))
	mux.HandleFunc("GET /api/v1/store", Get(service))
	mux.HandleFunc("PUT /api/v1/store", Put(service))
	return mux
}
