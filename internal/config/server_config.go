package config

import (
	"fmt"
	"github.com/spf13/viper"
)

type ServerConfig struct {
	Address string
}

func readServerConfig() ServerConfig {
	address := fmt.Sprintf("%s:%s", viper.GetString("server.host"), viper.GetString("server.port"))
	return ServerConfig{
		Address: address,
	}
}
