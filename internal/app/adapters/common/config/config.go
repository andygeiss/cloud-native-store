package config

type Config struct {
	Server Server `json:"server"`
}

type Server struct {
	CertFile string `json:"cert_file"`
	KeyFile  string `json:"key_file"`
}
