package core_auth

import (
	"fmt"
	"time"

	core_utils_env "github.com/wydentis/vykladna/shared/utils/env"
)

type Config struct {
	Secret     string        `env:"SECRET,required"`
	Issuer     string        `env:"ISSUER,,vykladna-app"`
	AccessTTL  time.Duration `env:"ACCESS_TTL,,15m"`
	RefreshTTL time.Duration `env:"REFRESH_TTL,,24h"`
}

func NewConfig() (Config, error) {
	var config Config

	if err := core_utils_env.Process(&config, "AUTH"); err != nil {
		return Config{}, fmt.Errorf("process envconfig: %w", err)
	}

	return config, nil
}

func NewConfigMust() Config {
	config, err := NewConfig()
	if err != nil {
		err := fmt.Errorf("get auth config: %w", err)
		panic(err)
	}

	return config
}
