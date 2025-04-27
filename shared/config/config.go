package config

import (
	"fmt"
	"os"

	yaml "gopkg.in/yaml.v3"
)

type Config struct {
	Server struct {
		GRPCPort  string `yaml:"grpc_port"`
		QueueSize int    `yaml:"queue_size"`
	} `yaml:"server"`
	Task struct {
		MaxDuration int `yaml:"max_duration"`
		MinDuration int `yaml:"min_duration"`
	} `yaml:"task"`
	Workers struct {
		Amount int `yaml:"amount"`
	} `yaml:"workers"`
}

func ReadConfig() (*Config, error) {
	f, err := os.ReadFile("../config.yaml")
	if err != nil {
		return nil, fmt.Errorf("config file does not exist")
	}

	c := Config{}
	if err := yaml.Unmarshal(f, &c); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config file: %v", err)
	}

	return &c, nil
}
