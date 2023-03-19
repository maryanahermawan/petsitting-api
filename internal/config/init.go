package config

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"strings"
)

type Config struct {
	DbConfig     DbConfig
	ServerConfig ServerConfig
}

func Init() Config {
	viper.SetConfigName("application")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		log.Errorf("Failed to load config: %v", err)
	}

	return Config{
		DbConfig:     readDbConfig(),
		ServerConfig: readServerConfig(),
	}
}
