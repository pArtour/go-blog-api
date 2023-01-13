package config

import "github.com/kelseyhightower/envconfig"

// Config is the configuration for the application
type Config struct {
	HttpAddress string `envconfig:"HTTP_ADDRESS" default:":8080"`
	DBAddress   string
}

// NewConfig returns a new Config struct with the values from the environment
func NewConfig() (*Config, error) {
	var cfg Config
	err := envconfig.Process("app", &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}
