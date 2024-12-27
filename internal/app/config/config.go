package config

type Config struct {
	Key  [32]byte `json:"-"`
	Port string   `json:"port"`
}
