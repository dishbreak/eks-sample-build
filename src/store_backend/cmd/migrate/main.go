package main

import (
	"log"

	"github.com/dishbreak/sample-store-backend/config"
	"github.com/dishbreak/sample-store-backend/migrations"
	"github.com/dishbreak/sample-store-backend/models"
	"github.com/pressly/goose/v3"
)

func main() {
	cfg := config.Must()

	db, err := models.Open(cfg.Database)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	goose.SetBaseFS(migrations.EmbeddedMigrations)

	if err := goose.SetDialect("postgres"); err != nil {
		panic(err)
	}

	if err := goose.Up(db, "sql"); err != nil {
		panic(err)
	}

	log.Println("database successfully migrated!")
}
