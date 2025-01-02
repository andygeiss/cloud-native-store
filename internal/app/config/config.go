package config

import "embed"

type Config struct {
	PortCloudSpanner PortCloudSpanner `json:"port_cloud_spanner"`
	Server           Server           `json:"server"`
	Service          Service          `json:"service"`
}

type PortCloudSpanner struct {
	DatabaseID string `json:"database_id"`
	InstanceID string `json:"instance_id"`
	ProjectID  string `json:"project_id"`
	Table      string `json:"table"`
}

type Server struct {
	Efs       embed.FS `json:"-"`
	Port      string   `json:"port"`
	Templates string   `json:"templates"`
}

type Service struct {
	Key [32]byte `json:"-"`
}
