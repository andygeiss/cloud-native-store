package api

import (
	"net/http"

	"github.com/andygeiss/cloud-native-store/internal/app/core/services"
)

func Route(service *services.ObjectService) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("DELETE /api/v1/store", Delete(service))
	mux.HandleFunc("GET /api/v1/store", Get(service))
	mux.HandleFunc("PUT /api/v1/store", Put(service))
	return mux
}
