package core_config

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	TimeZone *time.Location `envconfig:"TIME_ZONE"`
}

func NewConfig() (*Config, error) {
	cfg := &Config{}
	if err := envconfig.Process("", cfg); err != nil {
		return &Config{}, fmt.Errorf("process envconfig: %w", err)
	}

	return cfg, nil
}

func NewConfigMust() *Config {
	cfg, err := NewConfig()
	if err != nil {
		err = fmt.Errorf("get core config: %w", err)
		panic(err)
	}
	return cfg
}
