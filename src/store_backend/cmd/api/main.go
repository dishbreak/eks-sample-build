package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
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

	log.Println("loading config")
	cfg := config.Must()

	oidcMiddleware := myMiddleware.PassThru
	if cfg.OAuth2.AuthEnabled {
		provider, err := oidc.NewProvider(context.Background(), cfg.OAuth2.IssuerURL)
		if err != nil {
			panic(err)
		}

		verifier := provider.Verifier(&oidc.Config{
			ClientID: cfg.OAuth2.ExpectedAudience,
		})
		oidcMiddleware = myMiddleware.OIDCVerification(verifier)
	}

	uploadFs := mustPrepareDir(cfg.Assets)

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

	adminMiddleWare := myMiddleware.RequireScope("store.admin")
	readOnlyMiddleware := myMiddleware.RequireScope("store.read")

	r.With(middleware.StripPrefix("/assets")).Get("/assets/*", http.FileServer(uploadFs).ServeHTTP)
	r.Mount("/items", items.NewController(
		itemsSvc,
		items.WithOIDCVerifier(oidcMiddleware),
		items.WithAdminMiddleware(adminMiddleWare),
		items.WithReadOnlyMiddleware(readOnlyMiddleware)))
	r.Mount("/images", images.NewController(
		imagesSvc,
		images.WithOIDCVerifier(oidcMiddleware),
		images.WithAdminMiddleware(adminMiddleWare),
		images.WithReadOnlyMiddleware(readOnlyMiddleware)))

	log.Print("Listening on Port 8080!")
	http.ListenAndServe(":8080", r)
}

func mustPrepareDir(assetCfg config.AssetsConfig) http.FileSystem {
	path, err := filepath.Abs(assetCfg.UploadDir)
	if err != nil {
		panic(fmt.Errorf("failed to expand %s: %w", assetCfg.UploadDir, err))
	}
	log.Printf("preparing uploads dir %s", path)
	stat, err := os.Stat(path)
	if errors.Is(err, os.ErrNotExist) {
		log.Printf("%s does not exist, creating...", path)
		err := os.MkdirAll(path, os.FileMode(os.O_RDWR))
		if err != nil {
			panic(fmt.Errorf("failed to create dir %s: %w", path, err))
		}
	}

	if !stat.IsDir() {
		panic(fmt.Errorf("cannot use %s, is not dir", path))
	}

	return http.Dir(path)
}
