package config

type Config struct {
	Key    [32]byte `json:"-"`
	Server Server   `json:"server"`
}

type Server struct {
	CertFile string `json:"cert_file"`
	KeyFile  string `json:"key_file"`
}
