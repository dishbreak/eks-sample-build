-- +goose Up
ALTER TABLE store_items ADD CONSTRAINT store_items_pkey PRIMARY KEY (id);

-- +goose Down
ALTER TABLE store_items DROP CONSTRAINT store_items_pkey;
