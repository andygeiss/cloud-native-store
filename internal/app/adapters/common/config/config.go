package config

import "github.com/andygeiss/cloud-native-store/internal/app/core/services"

type Config struct {
	Key      [32]byte `json:"-"`
	Server   Server   `json:"server"`
	Services Services `json:"-"`
}

type Server struct {
	CertFile string `json:"cert_file"`
	KeyFile  string `json:"key_file"`
}

type Services struct {
	ObjectService *services.ObjectService `json:"-"`
}
