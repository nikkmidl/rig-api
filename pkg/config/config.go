package config

import (
	"github.com/spf13/viper"
)

// Config stores all configuration of the application.
// The values are read by viper from a config file or environment variable.
type EnvConfig struct {
	Environment       string   `mapstructure:"ENVIRONMENT"`
	AllowedOrigins    []string `mapstructure:"ALLOWED_ORIGINS"`
	GRPCServerAddress string   `mapstructure:"GRPC_SERVER_ADDRESS"`
	HTTPServerAddress string   `mapstructure:"HTTP_SERVER_ADDRESS"`
	GithubToken       string   `mapstructure:"GITHUB_TOKEN"`
	BlockedUsers      []string `mapstructure:"BLOCKED_USERS"`
}

var Config EnvConfig

func init() {
	var err error
	Config, err = LoadConfig(".")

	if err != nil {
		panic("Could not read configuration")
	}
}

// LoadConfig reads configuration from file or environment variables.
func LoadConfig(path string) (EnvConfig, error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		return Config, err
	}

	err = viper.Unmarshal(&Config)
	return Config, err
}
