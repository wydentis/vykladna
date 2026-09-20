package core_logger

import (
	"fmt"
	"log/slog"

	core_utils_env "github.com/wydentis/vykladna/shared/utils/env"
)

type Config struct {
	Level  slog.Level `env:"LEVEL,,info"`
	Folder string     `env:"FOLDER,,./logs"`
}

func NewConfig() (Config, error) {
	var cfg Config

	err := core_utils_env.Process(&cfg, "LOG")
	if err != nil {
		return Config{}, fmt.Errorf("failed to get logger config: %w", err)
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
