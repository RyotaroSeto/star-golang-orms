package infra

import (
	"context"
	"fmt"

	"github.com/sethvargo/go-envconfig"
)

type Config struct {
	// GitHubToken string `mapstructure:"GITHUB_TOKEN"`
	GitHubToken string `env:"GITHUB_TOKEN,required"`
}

var c *Config

func Load(ctx context.Context) error {
	// viper.AddConfigPath(path)
	// viper.SetConfigName("app")
	// viper.SetConfigType("env")

	// viper.AutomaticEnv()
	// err := viper.ReadInConfig()
	// if err != nil {
	// 	return errors.New("cannot read config")
	// }

	// var cfg Config
	// err = viper.Unmarshal(&cfg)
	// if err != nil {
	// 	return fmt.Errorf("cannot unmarshal config: %w", err)
	// }
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
