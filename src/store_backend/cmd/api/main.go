package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/dishbreak/sample-store-backend/controllers"
	"github.com/dishbreak/sample-store-backend/models"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
    db, err := models.Open(models.DBConfigFromEnv())
    if err != nil {
        panic(err)
    }
    defer db.Close()

    itemsSvc := models.NewItemService(db)
    
    r := chi.NewRouter()

    r.Use(middleware.Timeout(60 * time.Second))

    r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        fmt.Fprint(w, "Ok!")
    })

    r.Mount("/items/", controllers.NewItemController(itemsSvc))

    log.Print("Listening on Port 8080!")
    http.ListenAndServe(":8080", r)
}
