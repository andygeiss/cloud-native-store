package main

import (
	"log"

	"github.com/andygeiss/cloud-native-store/internal/app/adapters/common/config"
	"github.com/andygeiss/cloud-native-store/internal/app/adapters/inbound/api"
	"github.com/andygeiss/cloud-native-utils/security"
)

func main() {
	cfg := &config.Config{
		Server: config.Server{CertFile: ".tls/server.crt", KeyFile: ".tls/server.key"},
	}
	mux := api.Route(cfg)
	srv := security.NewServer(mux, "localhost")
	defer srv.Close()
	log.Printf("Start listening...")
	if err := srv.ListenAndServeTLS(cfg.Server.CertFile, cfg.Server.KeyFile); err != nil {
		log.Fatalf("ListenAndServeTLS failed: %v", err)
	}
}
