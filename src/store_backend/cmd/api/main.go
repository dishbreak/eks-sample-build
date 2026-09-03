package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/dishbreak/sample-store-backend/config"
	"github.com/dishbreak/sample-store-backend/controllers/images"
	"github.com/dishbreak/sample-store-backend/controllers/items"
	myMiddleware "github.com/dishbreak/sample-store-backend/middleware"
	"github.com/dishbreak/sample-store-backend/models"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {

	cfg := config.Must()

	provider, err := oidc.NewProvider(context.Background(), cfg.OAuth2.IssuerURL)
	if err != nil {
		panic(err)
	}

	verifier := provider.Verifier(&oidc.Config{
		ClientID: cfg.OAuth2.ExpectedAudience,
	})

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

	oidcMiddleware := myMiddleware.OIDCVerification(verifier)
	adminMiddleWare := myMiddleware.RequireScope("store.admin")
	readOnlyMiddleware := myMiddleware.RequireScope("store.read")

	r.Mount("/items/", items.NewController(
		itemsSvc,
		items.WithOIDCVerifier(oidcMiddleware),
		items.WithAdminMiddleware(adminMiddleWare),
		items.WithReadOnlyMiddleware(readOnlyMiddleware)))
	r.Mount("/images/", images.NewController(
		imagesSvc,
		images.WithOIDCVerifier(oidcMiddleware),
		images.WithAdminMiddleware(adminMiddleWare),
		images.WithReadOnlyMiddleware(readOnlyMiddleware)))

	log.Print("Listening on Port 8080!")
	http.ListenAndServe(":8080", r)
}
