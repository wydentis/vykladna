package core_http_server

import (
	"fmt"
	"time"

	core_utils_env "github.com/wydentis/vykladna/shared/utils/env"
)

type Config struct {
	Addr            string        `env:"ADDR,,0.0.0.0"`
	Port            int           `env:"PORT,,8080"`
	ShutdownTimeout time.Duration `env:"SHUTDOWN_TIMEOUT,,30s"`
}

func NewConfig() (Config, error) {
	var cfg Config

	err := core_utils_env.Process(&cfg, "HTTP_SERVER")
	if err != nil {
		return Config{}, fmt.Errorf("failed to get server config: %w", err)
	}

	return cfg, nil
}

func NewConfigMust() Config {
	cfg, err := NewConfig()
	if err != nil {
		panic(err)
	}

	return cfg
}
