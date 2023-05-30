package configs

import (
	"github.com/spf13/viper"
)

type Config struct {
	GithubToken string `mapstructure:"GITHUB_TOKEN"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AutomaticEnv()
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	err = viper.ReadInConfig()
	if err != nil {
		return
	}
	err = viper.Unmarshal(&config)
	return
}
