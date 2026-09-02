package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/dishbreak/sample-store-backend/config"
	"github.com/dishbreak/sample-store-backend/controllers/images"
	"github.com/dishbreak/sample-store-backend/controllers/items"
	"github.com/dishbreak/sample-store-backend/models"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {

	cfg := config.Must()

	db, err := models.Open(cfg.Database)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	itemsSvc := models.NewItemService(db)
	imagesSvc := models.NewImageService(db)

	r := chi.NewRouter()

	r.Use(middleware.Timeout(60 * time.Second))
	r.Use(middleware.RedirectSlashes)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "Ok!")
	})

	r.Mount("/items/", items.NewController(itemsSvc))
	r.Mount("/images/", images.NewController(imagesSvc))

	log.Print("Listening on Port 8080!")
	http.ListenAndServe(":8080", r)
}
