package config

import (
	"bytes"
	"errors"
	"fmt"

	"gopkg.in/yaml.v2"
)

var (
	ErrMissingRedisConnectionString   = errors.New("missing Redis connection string")
	ErrMissingMongoDBConnectionString = errors.New("missing MongoDB connection string")
)

type Config struct {
	MongoDBConnectionString string `yaml:"MongoDBConnectionString"`
	RedisConnectionString   string `yaml:"RedisConnectionString"`
}

// validate checks the configuration for missing required fields.
func (c *Config) validate() error {
	if c.MongoDBConnectionString == "" {
		return ErrMissingMongoDBConnectionString
	}
	if c.RedisConnectionString == "" {
		return ErrMissingRedisConnectionString
	}
	return nil
}

func Load(data []byte) (*Config, error) {
	// Create a strict YAML decoder
	decoder := yaml.NewDecoder(bytes.NewReader(data))
	decoder.SetStrict(true)

	// Decode the configuration into a ServerConfig struct
	cfg := &Config{}
	err := decoder.Decode(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to decode configuration: %w", err)
	}

	// Validate the configuration
	err = cfg.validate()
	if err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}
	return cfg, nil
}
