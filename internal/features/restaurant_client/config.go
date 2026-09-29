package restaurant_client

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	BaseURL   string `envconfig:"URL"        required:"true"`
	OrderPath string `envconfig:"ORDER_PATH" required:"true"`
}

func NewConfig() (Config, error) {
	var config Config

	if err := envconfig.Process("RESTAURANT", &config); err != nil {
		return Config{}, fmt.Errorf("process envconfig: %w", err)
	}

	return config, nil
}

func NewConfigMust() Config {
	config, err := NewConfig()

	if err != nil {
		panic(fmt.Errorf("get restaurant client config: %w", err))
	}

	return config
}
