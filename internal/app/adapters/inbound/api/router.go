package api

import (
	"context"
	"net/http"

	"github.com/andygeiss/cloud-native-store/internal/app/config"
	"github.com/andygeiss/cloud-native-store/internal/app/core/services"
	"github.com/andygeiss/cloud-native-utils/security"
	"github.com/andygeiss/cloud-native-utils/templating"
)

// Route creates a new mux with the liveness and readiness probe (/liveness, /readiness),
// the static assets endpoint (/) and the store endpoints (/api/v1/store).
func Route(service *services.ObjectService, ctx context.Context, cfg *config.Config) *http.ServeMux {
	// Create a new mux with liveness and readyness endpoint.
	// Embed the assets into the mux.
	mux := security.Mux(ctx, cfg.Server.Efs)

	// Add the store endpoints to the mux.
	mux.HandleFunc("DELETE /api/v1/store", Delete(service))
	mux.HandleFunc("GET /api/v1/store", Get(service))
	mux.HandleFunc("PUT /api/v1/store", Put(service))

	// Create a new templating engine and parse the templates.
	engine := templating.NewEngine(cfg.Server.Efs)
	engine.Parse(cfg.Server.Templates)

	// Add the UI endpoints for HTMX.
	mux.HandleFunc("GET /ui", ViewIndex(engine))
	mux.HandleFunc("GET /ui/store", ViewStore(engine))
	return mux
}
