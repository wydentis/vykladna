package core_pgx_pool

import (
	"fmt"
	"time"

	core_utils_env "github.com/wydentis/vykladna/shared/utils/env"
)

type Config struct {
	Host     string        `envconfig:"HOST" required:"true"`
	Port     string        `envconfig:"PORT" default:"5432"`
	User     string        `envconfig:"USER" required:"true"`
	Password string        `envconfig:"PASSWORD" required:"true"`
	Database string        `envconfig:"DB" required:"true"`
	Timeout  time.Duration `envconfig:"TIMEOUT" default:"30s"`
}

func NewConfig() (Config, error) {
	var config Config

	if err := core_utils_env.Process(&config, "POSTGRES"); err != nil {
		return Config{}, fmt.Errorf("process envconfig: %w", err)
	}

	return config, nil
}

func NewConfigMust() Config {
	config, err := NewConfig()
	if err != nil {
		err := fmt.Errorf("get postgres connection pool config: %w", err)
		panic(err)
	}

	return config
}
