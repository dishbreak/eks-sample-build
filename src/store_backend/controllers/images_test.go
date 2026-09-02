package controllers_test

import (
	"bytes"
	"embed"
	"encoding/json"
	"errors"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/textproto"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/dishbreak/sample-store-backend/controllers"
	"github.com/dishbreak/sample-store-backend/models"
	"github.com/stretchr/testify/assert"
)

type mockImageService struct {
	images    map[int]models.Image
	nextIndex int
	t         *testing.T
}

// Create implements [models.ImageService].
func (m *mockImageService) Create(itemId int, img models.Image) (models.Image, error) {
	assert.Equal(m.t, itemId, img.ItemId)
	img.Id = m.nextIndex
	m.nextIndex++
	m.images[img.Id] = img
	return img, nil
}

// DeleteById implements [models.ImageService].
func (m *mockImageService) DeleteById(itemId int) error {
	delete(m.images, itemId)
	return nil
}

// GetById implements [models.ImageService].
func (m *mockImageService) GetById(imageId int) (models.Image, error) {
	img, ok := m.images[imageId]
	if !ok {
		return img, errors.New("image not found")
	}
	return img, nil
}

// GetByItemId implements [models.ImageService].
func (m *mockImageService) GetByItemId(itemId int) ([]models.Image, error) {
	results := make([]models.Image, 0)
	for _, v := range m.images {
		if itemId != v.ItemId {
			continue
		}
		results = append(results, v)
	}

	return results, nil
}

func NewMockImageService(t *testing.T, images ...models.Image) models.ImageService {
	m := &mockImageService{
		images:    make(map[int]models.Image),
		nextIndex: 1,
		t:         t,
	}

	for _, img := range images {
		m.Create(img.ItemId, img)
	}

	return m
}

func TestGetById(t *testing.T) {
	t.Run("404 on nonexistent image", func(t *testing.T) {
		c := controllers.NewImageController(
			NewMockImageService(t),
		)

		rr := httptest.NewRecorder()

		req, err := http.NewRequest("GET", "/image/67", nil)
		assert.Nil(t, err, "failed to form test request")

		c.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusNotFound, rr.Result().StatusCode)
		assert.Equal(t, "image not found\n", rr.Body.String())
	})

	t.Run("return valid image", func(t *testing.T) {
		c := controllers.NewImageController(
			NewMockImageService(t, models.Image{ItemId: 45, Path: "some/fake/path"}),
		)

		rr := httptest.NewRecorder()

		req, err := http.NewRequest("GET", "/image/1", nil)
		assert.Nil(t, err, "failed to form test request")

		c.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Result().StatusCode)

		dec := json.NewDecoder(rr.Body)
		img := models.Image{}
		dec.Decode(&img)

		assert.Equal(t, 1, img.Id)
		assert.Equal(t, 45, img.ItemId)
		assert.Equal(t, "some/fake/path", img.Path)
	})
}

func TestGetByItemId(t *testing.T) {
	images := []models.Image{
		{ItemId: 23, Path: "some/fake/path/1"},
		{ItemId: 42, Path: "some/fake/path/2"},
		{ItemId: 23, Path: "some/fake/path/3"},
		{ItemId: 45, Path: "some/fake/path/4"},
		{ItemId: 22, Path: "some/fake/path/5"},
		{ItemId: 22, Path: "some/fake/path/6"},
		{ItemId: 27, Path: "some/fake/path/7"},
		{ItemId: 22, Path: "some/fake/path/8"},
		{ItemId: 23, Path: "some/fake/path/9"},
		{ItemId: 23, Path: "some/fake/path/0"},
	}

	t.Run("non-existent itemId returns empty slice", func(t *testing.T) {
		c := controllers.NewImageController(
			NewMockImageService(t, images...),
		)

		rr := httptest.NewRecorder()
		req, err := http.NewRequest("GET", "/item/76/", nil)
		assert.Nil(t, err, "failed to create test request")

		c.ServeHTTP(rr, req)

		// turns out that Go automatically sets 404 when the result is empty. who knew?
		assert.Equal(t, http.StatusNotFound, rr.Result().StatusCode)

		dec := json.NewDecoder(rr.Body)

		images := make([]models.Image, 0)
		dec.Decode(&images)

		assert.NotNil(t, images)
		assert.Empty(t, images)
	})

	t.Run("existent itemId returns only the images associated with the item", func(t *testing.T) {
		c := controllers.NewImageController(
			NewMockImageService(t, images...),
		)

		rr := httptest.NewRecorder()
		req, err := http.NewRequest("GET", "/item/23", nil)
		assert.Nil(t, err, "failed to create test request")

		c.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Result().StatusCode)

		dec := json.NewDecoder(rr.Body)
		images := make([]models.Image, 0)
		dec.Decode(&images)
		//no guarantee on the order of the results, so sort by id
		slices.SortStableFunc(images, func(a, b models.Image) int {
			return a.Id - b.Id
		})

		expected := []models.Image{
			{Id: 1, ItemId: 23, Path: "some/fake/path/1"},
			{Id: 3, ItemId: 23, Path: "some/fake/path/3"},
			{Id: 9, ItemId: 23, Path: "some/fake/path/9"},
			{Id: 10, ItemId: 23, Path: "some/fake/path/0"},
		}
		assert.Equal(t, expected, images)
	})
}

//go:embed files/*
var EmbedTestFiles embed.FS

func TestDelete(t *testing.T) {
	c := controllers.NewImageController(
		NewMockImageService(t, models.Image{ItemId: 77, Path: "some/faker/path"}),
	)

	rr := httptest.NewRecorder()
	req, err := http.NewRequest("DELETE", "/image/1", nil)
	assert.Nil(t, err, "failed to create test request")

	c.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Result().StatusCode)

	// a subsequent get of the image should return 404
	req, err = http.NewRequest("GET", "/image/1", nil)
	rr = httptest.NewRecorder()

	c.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Result().StatusCode)
}

type fileUpload struct {
	Path, ContentType string
}

func prepareForm(url string, filePaths ...fileUpload) *http.Request {
	var buf bytes.Buffer
	mw := multipart.NewWriter(&buf)

	// this code is a bit panicky. that's ok because it's for a test. if you're
	// reading this and want to use it, make sure you're doing something
	// sensible with the errors instead of panicking at the slightest
	// provocation.
	for _, fp := range filePaths {
		func(fp fileUpload) {
			fd, err := EmbedTestFiles.Open(fp.Path)
			if err != nil {
				panic(err)
			}
			defer fd.Close()

			stat, err := fd.Stat()
			if err != nil {
				panic(err)
			}

			header := make(textproto.MIMEHeader)
			header.Set("Content-Disposition", mime.FormatMediaType(
				"form-data",
				map[string]string{
					"name":     "images",
					"filename": filepath.Base(fp.Path),
				},
			))
			header.Set("Content-Type", fp.ContentType)
			header.Set("Content-Length", strconv.FormatInt(
				stat.Size(), 10,
			))

			part, err := mw.CreatePart(header)
			if err != nil {
				panic(err)
			}

			n, err := io.Copy(part, fd)
			if err != nil {
				panic(err)
			}

			if int64(n) != stat.Size() {
				panic(errors.New("file size changed during write"))
			}
		}(fp)
	}
	if err := mw.Close(); err != nil {
		panic(err)
	}

	req, err := http.NewRequest("POST", url, io.NopCloser(&buf))
	if err != nil {
		panic(err)
	}

	req.Header.Set("Content-Type", mw.FormDataContentType())
	req.ContentLength = int64(buf.Len())

	return req
}

func TestUpload(t *testing.T) {
	staticTime, _ := time.Parse(time.RFC3339, "2026-09-01T21:08:10.000Z")
	staticTimeProvider := func() time.Time {
		return staticTime
	}
	t.Run("will handle a normal image upload", func(t *testing.T) {
		dir := t.TempDir()

		c := controllers.NewImageController(
			NewMockImageService(t),
			controllers.WithTimeProvider(staticTimeProvider),
			controllers.WithImageUploadDir(dir),
		)

		rr := httptest.NewRecorder()

		uploads := []fileUpload{
			{Path: "files/417.jpg", ContentType: "images/jpeg"},
			{Path: "files/424.jpg", ContentType: "images/jpeg"},
			{Path: "files/426.jpg", ContentType: "images/jpeg"},
		}

		req := prepareForm(
			"/item/7",
			uploads...,
		)
		c.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Result().StatusCode)
		entries, err := os.ReadDir(dir)
		assert.Nil(t, err)

		fileNames := make([]string, len(entries))
		for i, e := range entries {
			fileNames[i] = e.Name()
		}
		slices.SortFunc(fileNames, strings.Compare)
		expectedFiles := []string{
			"20260901T210810.000000000Z-06a079d5beb2b2d13e7343d12c81eaac58a1268424f8289c484abcdc1065309c.jpg",
			"20260901T210810.000000000Z-77e4af13a7017b797e90d9a0e86d3028f8a1238fbe03fca03a04d7f4e0d49c4b.jpg",
			"20260901T210810.000000000Z-fbdbce156d2935502704497f8cda394cb97d76c90bba1d7b444a2544f01f64c6.jpg",
		}
		assert.Equal(t, expectedFiles, fileNames)

	})

	t.Run("will reject an upload of a non-image file", func(t *testing.T) {
		dir := t.TempDir()

		c := controllers.NewImageController(
			NewMockImageService(t),
			controllers.WithTimeProvider(staticTimeProvider),
			controllers.WithImageUploadDir(dir),
		)

		rr := httptest.NewRecorder()

		req := prepareForm(
			"/item/7",
			fileUpload{
				Path:        "files/not_an_img.txt",
				ContentType: "text/plain",
			},
		)

		c.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Result().StatusCode)
	})

	t.Run("will detect and reject an upload of a fake jpg", func(t *testing.T) {
		dir := t.TempDir()

		c := controllers.NewImageController(
			NewMockImageService(t),
			controllers.WithTimeProvider(staticTimeProvider),
			controllers.WithImageUploadDir(dir),
		)

		rr := httptest.NewRecorder()

		req := prepareForm(
			"/item/7",
			fileUpload{
				Path:        "files/not_an_img.jpg",
				ContentType: "images/jpeg",
			},
		)

		c.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Result().StatusCode)
	})

	t.Run("will reject a mismatch between file extension and content type", func(t *testing.T) {
		dir := t.TempDir()

		c := controllers.NewImageController(
			NewMockImageService(t),
			controllers.WithTimeProvider(staticTimeProvider),
			controllers.WithImageUploadDir(dir),
		)

		rr := httptest.NewRecorder()

		req := prepareForm(
			"/item/7",
			fileUpload{
				Path:        "files/not_an_img.txt",
				ContentType: "images/jpeg",
			},
		)

		c.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Result().StatusCode)
	})

	t.Run("will reject arbitrary data", func(t *testing.T) {
		dir := t.TempDir()

		c := controllers.NewImageController(
			NewMockImageService(t),
			controllers.WithTimeProvider(staticTimeProvider),
			controllers.WithImageUploadDir(dir),
		)

		rr := httptest.NewRecorder()

		formBody := "helloThere"
		req, err := http.NewRequest("POST", "/item/7", strings.NewReader(formBody))
		assert.Nil(t, err)
		req.Header.Set("Content-Type", "text/plain")

		c.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusUnsupportedMediaType, rr.Result().StatusCode)
	})
}
