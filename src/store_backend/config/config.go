package config

import (
	"fmt"

	"go-simpler.org/env"
)

type Config struct {
	Database DatabaseConfig `env:"DATABASE_"`
	OAuth2   OAuth2Config   `env:"OAUTH_"`
	Assets   AssetsConfig   `env:"ASSETS_"`
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

type OAuth2Config struct {
	AuthEnabled      bool   `env:"ENABLED" default:"true"`
	IssuerURL        string `env:"ISSUER_URL,required"`
	ExpectedAudience string `env:"EXPECTED_AUDIENCE" default:"store_backend"`
}

type AssetsConfig struct {
	UploadDir string `env:"UPLOAD_DIR" default:"./uploads"`
}

func Must() (cfg Config) {
	if err := env.Load(&cfg, nil); err != nil {
		panic(err)
	}
	return
}
