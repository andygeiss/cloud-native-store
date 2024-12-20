package api

import (
	"encoding/json"
	"net/http"

	"github.com/andygeiss/cloud-native-store/internal/app/adapters/common/config"
)

// Delete defines an HTTP handler function for deleting an object by key.
func Delete[K comparable, V any](cfg *config.Config[K, V]) http.HandlerFunc {

	type request struct {
		Key K `json:"key"`
	}

	type response struct{}

	return func(w http.ResponseWriter, r *http.Request) {
		var req request
		var res response

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if err := cfg.Services.ObjectService.Delete(r.Context(), req.Key); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if err := json.NewEncoder(w).Encode(res); err != nil {
		}
	}
}

// Get defines an HTTP handler function for retrieving an object by key.
func Get[K comparable, V any](cfg *config.Config[K, V]) http.HandlerFunc {

	type request struct {
		Key K `json:"key"`
	}

	type response struct {
		Value V `json:"value,omitempty"`
	}

	// Return the HTTP handler function.
	return func(w http.ResponseWriter, r *http.Request) {
		var req request
		var res response

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		val, err := cfg.Services.ObjectService.Get(r.Context(), req.Key)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		res.Value = val

		if err := json.NewEncoder(w).Encode(res); err != nil {
		}
	}
}

// Put defines an HTTP handler function for creating or updating an object.
func Put[K comparable, V any](cfg *config.Config[K, V]) http.HandlerFunc {

	type request struct {
		Key   K `json:"key"`
		Value V `json:"value"`
	}

	type response struct{}

	return func(w http.ResponseWriter, r *http.Request) {
		var req request
		var res response

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if err := cfg.Services.ObjectService.Put(r.Context(), req.Key, req.Value); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if err := json.NewEncoder(w).Encode(res); err != nil {
		}
	}
}
