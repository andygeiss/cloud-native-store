package config

import "github.com/andygeiss/cloud-native-store/internal/app/core/services"

type Config[K comparable, V any] struct {
	Server   Server `json:"server"`
	Services Services[K, V]
}

type Server struct {
	CertFile string `json:"cert_file"`
	KeyFile  string `json:"key_file"`
}

type Services[K comparable, V any] struct {
	ObjectService *services.ObjectService[K, V]
}
