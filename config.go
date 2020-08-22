package main

import (
	"github.com/spf13/viper"
)

type config struct {
	DB   *dbConfig
	File *fileConfig
}

// YAML Usage:
//
// db:
//   host: xxx
func loadConfig() *config {
	c := &config{
		DB:   &dbConfig{},
		File: &fileConfig{},
	}

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("res")

	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	if err := viper.Unmarshal(c); err != nil {
		panic(err)
	}

	return c
}
