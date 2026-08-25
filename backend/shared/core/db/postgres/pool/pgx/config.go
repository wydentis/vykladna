package core_pgx_pool

import (
	"fmt"
	"time"

	core_utils_env "github.com/wydentis/vykladna/shared/utils/env"
)

type Config struct {
	Host     string        `env:"HOST,required"`
	Port     int           `env:"PORT,,5432"`
	User     string        `env:"USER,required"`
	Password string        `env:"PASSWORD,required"`
	Database string        `env:"DB,required"`
	Timeout  time.Duration `env:"TIMEOUT,,30s"`
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
