package controllers

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"mime"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/dishbreak/sample-store-backend/models"
	"github.com/go-chi/chi/v5"
)

type images struct {
	imagesSvc     models.ImageService
	imagePath     string
	maxUploadSize int64
	timeProvider  func() time.Time
}

type ImagesOption func(i *images)

func WithImageUploadDir(path string) ImagesOption {
	return func(i *images) {
		i.imagePath = path
	}
}

func WithMaxUploadBytes(maxUpload int64) ImagesOption {
	return func(i *images) {
		i.maxUploadSize = maxUpload
	}
}

func WithImagesService(imagesSvc models.ImageService) ImagesOption {
	return func(i *images) {
		i.imagesSvc = imagesSvc
	}
}

// useful for testing to provide deterministic results
func WithTimeProvider(cb func() time.Time) ImagesOption {
	return func(i *images) {
		i.timeProvider = cb
	}
}

func (i *images) GetAllForItemId(w http.ResponseWriter, r *http.Request) {
	itemId := ItemIdFromContext(r.Context())

	results, err := i.imagesSvc.GetByItemId(itemId)
	if err != nil {
		http.Error(w, "no images found for item", http.StatusNotFound)
	}

	enc := json.NewEncoder(w)
	enc.Encode(results)
}

func (i *images) Get(w http.ResponseWriter, r *http.Request) {
	imageId := ImageIdFromContext(r.Context())

	image, err := i.imagesSvc.GetById(imageId)
	if err != nil {
		http.Error(w, "image not found", http.StatusNotFound)
		return
	}

	enc := json.NewEncoder(w)
	enc.Encode(image)
}

func (i *images) Upload(w http.ResponseWriter, r *http.Request) {
	itemId := ItemIdFromContext(r.Context())

	contentType := r.Header.Get("Content-Type")
	if contentType == "" {
		http.Error(w, "content not supported", http.StatusUnsupportedMediaType)
		return
	}

	mediaType, params, err := mime.ParseMediaType(contentType)
	if err != nil {
		http.Error(w, "invalid value for Content-Type", http.StatusBadRequest)
		return
	}

	if mediaType != "multipart/form-data" {
		http.Error(w, "unsupported media type", http.StatusUnsupportedMediaType)
		return
	}

	var boundary string
	if parsed, ok := params["boundary"]; !ok {
		http.Error(w, "missing boundary on content type", http.StatusBadRequest)
		return
	} else {
		boundary = parsed
	}

	r.Body = http.MaxBytesReader(w, r.Body, i.maxUploadSize)
	defer r.Body.Close()
	multiPartReader := multipart.NewReader(r.Body, boundary)

	for {
		p, err := multiPartReader.NextPart()
		if err == io.EOF {
			break
		}

		if err != nil {
			http.Error(w, "bad multipart file upload", http.StatusBadRequest)
			log.Printf("bad multipart file upload: %s", err)
			return
		}

		if p.FileName() == "" {
			continue
		}

		if path, err := i.savePart(p); err != nil {
			http.Error(w, "failed to complete upload", http.StatusBadRequest)
			log.Printf("failed to complete upload: %s", err)
			return
		} else {
			image := models.Image{Path: path, ItemId: itemId}
			i.imagesSvc.Create(itemId, image)
		}
	}
}

func (i *images) savePart(p *multipart.Part) (string, error) {
	ext := strings.ToLower(
		filepath.Ext(filepath.Base(p.FileName())),
	)
	r, contentType, err := detectContentType(p)

	if !validateContentAndExt(ext, contentType) {
		return "", errors.New("unsupported file upload")
	}

	b, err := io.ReadAll(r)
	if err != nil {
		return "", fmt.Errorf("failed to read part: %w", err)
	}

	hasher := sha256.New()
	tempFile, err := os.CreateTemp(i.imagePath, "upload-*")
	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %w", err)
	}

	writer := io.MultiWriter(hasher, tempFile)
	if _, err := writer.Write(b); err != nil {
		return "", fmt.Errorf("failed to write to temp file: %w", err)
	}
	tempPath := tempFile.Name()
	removeTemp := true

	defer func() {
		tempFile.Close()

		if removeTemp {
			_ = os.Remove(tempPath)
		}
	}()

	hash := hex.EncodeToString(hasher.Sum(nil))
	stamp := i.timeProvider().UTC().Format("20060102T150405.000000000Z")

	finalName := fmt.Sprintf("%s-%s%s", stamp, hash, ext)
	finalPath := filepath.Join(i.imagePath, finalName)

	if err := os.Rename(tempPath, finalPath); err != nil {
		return "", fmt.Errorf("failed to rename file: %w", err)
	}
	removeTemp = false

	return "", nil
}

func detectContentType(p *multipart.Part) (io.Reader, string, error) {
	header := make([]byte, 512)

	n, err := p.Read(header)
	if err != nil && err != io.ErrUnexpectedEOF && err != io.EOF {
		return nil, "unknown", fmt.Errorf("failed to read file header: %w", err)
	}

	header = header[:n]
	contentType := http.DetectContentType(header)
	replacementReader := io.MultiReader(bytes.NewReader(header), p)

	return replacementReader, contentType, nil
}

var acceptedTypesForExt map[string]string = map[string]string{
	".jpg":  "image/jpeg",
	".jpeg": "image/jpeg",
	".png":  "image/png",
}

func validateContentAndExt(ext, contentType string) bool {
	expectedContentType, ok := acceptedTypesForExt[ext]
	if !ok {
		return false
	}
	return contentType == expectedContentType
}

func (i *images) Delete(w http.ResponseWriter, r *http.Request) {
	imageId := ImageIdFromContext(r.Context())

	if err := i.imagesSvc.DeleteById(imageId); err != nil {
		http.Error(w, "failed to delete image", http.StatusInternalServerError)
		return
	}
}

func NewImageController(imagesSvc models.ImageService, opts ...ImagesOption) http.Handler {
	ic := &images{
		imagesSvc:     imagesSvc,
		imagePath:     "./assets",
		maxUploadSize: 100 << 20,
		timeProvider:  time.Now,
	}

	for _, opt := range opts {
		opt(ic)
	}

	r := chi.NewRouter()

	r.Group(func(r chi.Router) {
		r.Use(IntegerPathParam("itemId"))
		r.Get("/item/{itemId}", ic.GetAllForItemId)
		r.Post("/item/{itemId}", ic.Upload)
	})

	r.Group(func(r chi.Router) {
		r.Use(IntegerPathParam("imageId"))
		r.Get("/image/{imageId}", ic.Get)
		r.Delete("/image/{imageId}", ic.Delete)
	})

	return r
}
