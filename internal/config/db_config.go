package config

import (
	"fmt"
	"github.com/spf13/viper"
)

type DbConfig struct {
	Url string
}

func readDbConfig() DbConfig {
	dbname := viper.Get("database.name")
	host := viper.Get("database.host")
	user := viper.Get("database.user")
	password := viper.Get("database.password")

	url := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", user, password, host, dbname)

	return DbConfig{
		Url: url,
	}
}
