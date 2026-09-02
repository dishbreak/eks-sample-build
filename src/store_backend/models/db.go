package models

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"

	_ "github.com/jackc/pgx/v4/stdlib"
)

type DbConfig struct {
	Host     string
	Username string
	Password string
	DbName   string
	Port     int
	SslMode  string
}

func DefaultDBConfig() *DbConfig {
	c := DbConfig{
		Host:     "localhost",
		Username: "username",
		Password: "facepunch",
		DbName:   "store_backend",
		Port:     5432,
		SslMode:  "disable",
	}

	return &c
}

func DBConfigFromEnv() *DbConfig {
	c := DefaultDBConfig()
	if username, ok := os.LookupEnv("DATABASE_USERRNAME"); ok {
		c.Username = username
	}

	if password, ok := os.LookupEnv("DATABASE_PASSWORD"); ok {
		c.Password = password
	}

	if hostname, ok := os.LookupEnv("DATABASE_HOST"); ok {
		c.Host = hostname
	}

	if dbname, ok := os.LookupEnv("DATABASE_DBNAME"); ok {
		c.DbName = dbname
	}

	if sslMode, ok := os.LookupEnv("DATABASE_SSLMODE"); ok {
		c.SslMode = sslMode
	}

	if val, ok := os.LookupEnv("DATABASE_PORT"); ok {
		if port, err := strconv.Atoi(val); err == nil {
			c.Port = port
		}
	}

	return c
}

func (d *DbConfig) ToDSN() string {
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=%s",
		d.Host, d.Username, d.Password, d.DbName, d.Port, d.SslMode,
	)
}

func Open(d *DbConfig) (*sql.DB, error) {
	return sql.Open("pgx", d.ToDSN())
}
