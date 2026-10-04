package core_tgbot_server

import (
	"fmt"

	core_utils_env "github.com/wydentis/vykladna/shared/utils/env"
)

type Config struct {
	Secret    string `env:"WEBHOOK_SECRET"`
	Workers   int    `env:"WORKERS,,8"`
	QueueSize int    `env:"QUEUE_SIZE,,256"`
}

func NewConfig() (Config, error) {
	var cfg Config

	if err := core_utils_env.Process(&cfg, "TGBOT_SERVER"); err != nil {
		return Config{}, fmt.Errorf("failed to get telegram server config: %w", err)
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
