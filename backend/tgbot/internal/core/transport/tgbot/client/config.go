package core_tgbot_client

import (
	"fmt"
	"time"

	core_utils_env "github.com/wydentis/vykladna/shared/utils/env"
)

type Config struct {
	Token   string        `env:"BOT_TOKEN,required"`
	APIURL  string        `env:"API_URL,,https://api.telegram.org"`
	Timeout time.Duration `env:"TIMEOUT,,10s"`
}

func NewConfig() (Config, error) {
	var cfg Config

	if err := core_utils_env.Process(&cfg, "TGBOT_CLIENT"); err != nil {
		return Config{}, fmt.Errorf("failed to get telegram client config: %w", err)
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
