package infra

import (
	"errors"
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	GitHubToken string `mapstructure:"GITHUB_TOKEN"`
}

var c *Config

func Load(path string) error {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		return errors.New("cannot read config")
	}

	var cfg Config
	err = viper.Unmarshal(&cfg)
	if err != nil {
		return fmt.Errorf("cannot unmarshal config: %w", err)
	}

	c = &cfg
	return nil
}

func Get() *Config {
	return c
}
