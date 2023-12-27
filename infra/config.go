package infra

import (
	"context"
	"fmt"

	"github.com/sethvargo/go-envconfig"
)

type Config struct {
	GitHubToken string `env:"GITHUB_TOKEN,required"`
}

var c *Config

func Load(ctx context.Context) error {
	var cfg Config
	if err := envconfig.Process(ctx, &cfg); err != nil {
		return fmt.Errorf("cannot unmarshal config: %w", err)
	}

	c = &cfg
	return nil
}

func Get() *Config {
	return c
}
