package controllers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/dishbreak/sample-store-backend/models"
	"github.com/go-chi/chi/v5"
)

type items struct {
    itemsSvc models.ItemService
}

func (i items) GetById(w http.ResponseWriter, r *http.Request) {
    item, ok := ItemFromCtx(r.Context())
    if !ok {
        http.Error(w, http.StatusText(http.StatusUnprocessableEntity), http.StatusUnprocessableEntity)
        return
    }

    enc := json.NewEncoder(w)
    enc.Encode(item)
}

func (i items) Get(w http.ResponseWriter, r *http.Request) {
    result, err := i.itemsSvc.GetAll()
    if err != nil {
        http.Error(w, "failed to get items", http.StatusInternalServerError)
        log.Printf("failed to get items: %s", err)
        return
    }

    enc := json.NewEncoder(w)
    enc.Encode(result)
}

func (i items) Delete(w http.ResponseWriter, r *http.Request) {
    item, ok := ItemFromCtx(r.Context())
    if !ok {
        http.Error(w, http.StatusText(http.StatusUnprocessableEntity), http.StatusUnprocessableEntity)
        return
    }

    err := i.itemsSvc.Delete(item)
    if err != nil {
        http.Error(w, "failed to delete item", http.StatusInternalServerError)
        log.Printf("failed to delete item: %s", err)
        return
    }
}

func (i items) Update(w http.ResponseWriter, r *http.Request) {
    item, ok := ItemFromCtx(r.Context())
    if !ok {
        http.Error(w, http.StatusText(http.StatusUnprocessableEntity), http.StatusUnprocessableEntity)
        return
    }

    dec := json.NewDecoder(r.Body)
    incoming := models.Item{}
    
    if err := dec.Decode(&incoming); err != nil {
        http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
        return
    }

    item.Title = incoming.Title
    item.Description = incoming.Description
    
    if err := i.itemsSvc.Update(item); err != nil {
        http.Error(w, "failed to update item", http.StatusInternalServerError)
        log.Printf("failed to update item: %s", err)
    }
}

func (i items) Create (w http.ResponseWriter, r *http.Request) {
    incoming := models.Item{}
    dec := json.NewDecoder(r.Body)
    if err := dec.Decode(&incoming); err != nil {
        http.Error(w, "unable to parse input", http.StatusBadRequest)
        return
    }
    created, err := i.itemsSvc.Create(incoming)
    if err != nil {
        http.Error(w, "unable to create item", http.StatusBadRequest)
        log.Printf("unable to create item: %s", err)
        return
    }

    enc := json.NewEncoder(w)
    enc.Encode(created)
}

func NewItemController(itemsSvc models.ItemService) http.Handler {
    ic := items{itemsSvc: itemsSvc}

    r := chi.NewRouter()
    r.Get("/", ic.Get)
    r.Post("/", ic.Create)
    r.Route("/{itemId}", func(r chi.Router) {
        r.Use(ic.ItemCtx)
        r.Get("/", ic.Get)
        r.Put("/", ic.Update)
        r.Delete("/", ic.Delete)
    })
    return r
}

func (i items) ItemCtx(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        var itemId int
        itemParam := chi.URLParam(r, "itemId")
        if parsed, err := strconv.Atoi(itemParam); err != nil {
            http.Error(w, "invalid item id", http.StatusBadRequest)
            return
        } else {
            itemId = parsed
        }
        item, err := i.itemsSvc.GetById(itemId)
        if err != nil {
            http.Error(w, "no matching item", http.StatusNotFound)
            return
        }

        ctx := CtxWithItem(r.Context(), item)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}

func CtxWithItem(ctx context.Context, item models.Item) context.Context {
    return context.WithValue(ctx, "item", item)
}

func ItemFromCtx(ctx context.Context) (models.Item, bool) {
    item, ok := ctx.Value("item").(models.Item)
    return  item, ok
}
