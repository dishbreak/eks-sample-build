package config

import (
	"fmt"

	"go-simpler.org/env"
)

type Config struct {
	Database DatabaseConfig `env:"DATABASE_"`
}

type DatabaseConfig struct {
	Host     string `env:"HOST" default:"localhost"`
	Username string `env:"USERNAME" default:"username"`
	Password string `env:"PASSWORD" default:"facepunch"`
	DbName   string `env:"DBNAME" default:"store_backend"`
	Port     int    `env:"PORT" default:"5432"`
	SslMode  string `env:"SSLMODE" default:"disable"`
}

func (d DatabaseConfig) ToDSN() string {
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=%s",
		d.Host, d.Username, d.Password, d.DbName, d.Port, d.SslMode,
	)
}

func Must() (cfg Config) {
	if err := env.Load(&cfg, nil); err != nil {
		panic(err)
	}
	return
}
