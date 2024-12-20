package api

import (
	"encoding/json"
	"net/http"

	"github.com/andygeiss/cloud-native-store/internal/app/adapters/common/config"
	"github.com/andygeiss/cloud-native-utils/security"
)

// Delete defines an HTTP handler function for deleting an object by key.
// It expects a JSON request body with the "key" field and deletes the corresponding object.
func Delete(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Define the structure of the request payload.
		var req struct {
			Key string `json:"key"` // The key to identify the object to delete.
		}
		// Define an empty response structure.
		var res struct{}

		// Decode the JSON request body into the req struct.
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			// Respond with a 400 Bad Request status if the body cannot be parsed.
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Attempt to delete the object identified by the key.
		if err := cfg.Services.ObjectService.Delete(r.Context(), req.Key); err != nil {
			// Respond with a 500 Internal Server Error if the deletion fails.
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Respond with a 200 OK status if the deletion is successful.
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(res)
	}
}

// Get defines an HTTP handler function for retrieving an object by key.
// It expects a JSON request body with the "key" field and retrieves the corresponding object.
func Get(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Define the structure of the request payload.
		var req struct {
			Key string `json:"key"` // The key to identify the object to retrieve.
		}
		// Define the structure of the response payload.
		var res struct {
			Value string `json:"value,omitempty"` // The retrieved value, if found.
		}

		// Decode the JSON request body into the req struct.
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			// Respond with a 400 Bad Request status if the body cannot be parsed.
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Retrieve the encrypted value associated with the key.
		encValue, err := cfg.Services.ObjectService.Get(r.Context(), req.Key)
		if err != nil {
			// Respond with a 404 Not Found status if the object is not found.
			w.WriteHeader(http.StatusNotFound)
			return
		}

		// Decrypt the retrieved value.
		decValue, err := security.Decrypt([]byte(encValue), cfg.Key)
		if err != nil {
			// Respond with a 500 Internal Server Error if decryption fails.
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Set the decrypted value in the response.
		res.Value = string(decValue)

		// Respond with a 200 OK status and encode the response as JSON.
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(res)
	}
}

// Put defines an HTTP handler function for creating or updating an object.
// It expects a JSON request body with "key" and "value" fields.
func Put(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Define the structure of the request payload.
		var req struct {
			Key   string `json:"key"`   // The key to identify the object.
			Value string `json:"value"` // The value to store or update.
		}
		// Define an empty response structure.
		var res struct{}

		// Decode the JSON request body into the req struct.
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			// Respond with a 400 Bad Request status if the body cannot be parsed.
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Encrypt the value before storing it.
		encValue := security.Encrypt([]byte(req.Value), cfg.Key)

		// Store or update the object identified by the key with the encrypted value.
		if err := cfg.Services.ObjectService.Put(r.Context(), req.Key, string(encValue)); err != nil {
			// Respond with a 500 Internal Server Error if storing/updating fails.
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Respond with a 200 OK status and encode the response as JSON.
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(res)
	}
}
