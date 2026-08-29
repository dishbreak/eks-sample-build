package main

import (
	"embed"
	"log"

	"github.com/dishbreak/sample-store-backend/models"
	"github.com/pressly/goose/v3"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

func main() {
    db, err := models.Open(models.DBConfigFromEnv())
    if err != nil {
        panic(err)
    }
    defer db.Close()

    goose.SetBaseFS(embedMigrations)

    if err := goose.SetDialect("postgres"); err != nil {
        panic(err)
    }

    if err := goose.Up(db, "migrations"); err != nil {
        panic(err)
    }
    
    log.Println("database successfully migrated!")
}
