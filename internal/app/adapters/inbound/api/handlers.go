package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/andygeiss/cloud-native-store/internal/app/core/services"
)

// Delete defines an HTTP handler function for deleting an object by key.
// It expects a JSON request body with the "key" field and deletes the corresponding object.
func Delete(service *services.ObjectService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Key string `json:"key"`
		}
		var res struct{}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if err := service.Delete(r.Context(), req.Key); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Printf("service.Delete error: %v", err)
			return
		}

		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(res)
	}
}

// Get defines an HTTP handler function for retrieving an object by key.
// It expects a JSON request body with the "key" field and retrieves the corresponding object.
func Get(service *services.ObjectService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Key string `json:"key"`
		}
		var res struct {
			Value string `json:"value,omitempty"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		value, err := service.Get(r.Context(), req.Key)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			log.Printf("service.Get error: %v", err)
			return
		}

		res.Value = value

		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(res)
	}
}

// Put defines an HTTP handler function for creating or updating an object.
// It expects a JSON request body with "key" and "value" fields.
func Put(service *services.ObjectService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Key   string `json:"key"`
			Value string `json:"value"`
		}
		var res struct{}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if err := service.Put(r.Context(), req.Key, req.Value); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Printf("service.Put error: %v", err)
			return
		}

		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(res)
	}
}
