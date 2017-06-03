package main

import (
	"github.com/Sirupsen/logrus"
	"github.com/spf13/viper"
)

var version = "development"

var log = logrus.New()

var config Config

type Config struct {
	Database      string
	Listen        string
	IsDevelopment bool
}

func init() {
	log.Formatter = &logrus.TextFormatter{
		DisableColors: true,
	}
	if version == "development" {
		config.IsDevelopment = true
	}

	viper.SetDefault("database", "postgres://127.0.0.1:5432/app?sslmode=disable")
	viper.SetDefault("listen", ":4000")
	viper.AutomaticEnv()

	config.Database = viper.GetString("database")
	if viper.GetBool("dev") {
		config.IsDevelopment = true
	}
	config.Listen = viper.GetString("listen")

	log.WithFields(logrus.Fields{
		"version":       version,
		"isDevelopment": config.IsDevelopment,
		"database":      config.Database,
		"listen":        config.Listen,
	}).Print("config")
}
