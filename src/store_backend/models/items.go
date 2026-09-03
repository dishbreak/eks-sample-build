package models

import (
	"database/sql"
	"fmt"
)

type Item struct {
	Id          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type itemServiceImpl struct {
	db *sql.DB
}

type ItemService interface {
	GetAll() ([]Item, error)
	GetById(id int) (Item, error)
	Update(i Item) error
	Create(i Item) (Item, error)
	Delete(i Item) error
}

func NewItemService(db *sql.DB) ItemService {
	return &itemServiceImpl{
		db: db,
	}
}

func (i *itemServiceImpl) GetAll() ([]Item, error) {
	result := make([]Item, 0)
	rows, err := i.db.Query("SELECT id, title, description from store_items;")
	if err != nil {
		return result, fmt.Errorf("failed to retrieve items from server: %w", err)
	}

	for rows.Next() {
		item := Item{}
		if err := rows.Scan(&item.Id, &item.Title, &item.Description); err != nil {
			return result, fmt.Errorf("failed to process query result: %w", err)
		}
		result = append(result, item)
	}

	if err := rows.Err(); err != nil {
		return result, fmt.Errorf("failed while iterating over query results: %w", err)
	}

	return result, nil
}

// Create implements [ItemService].
func (i *itemServiceImpl) Create(it Item) (Item, error) {
	row := i.db.QueryRow("insert into store_items (id, title, description) VALUES (DEFAULT, $1, $2) RETURNING id, title, description;", it.Title, it.Description)
	result := Item{}
	if err := row.Scan(&result.Id, &result.Title, &result.Description); err != nil {
		return result, fmt.Errorf("failed to insert item: %w", err)
	}
	return result, nil
}

// Delete implements [ItemService].
func (i *itemServiceImpl) Delete(it Item) error {
	_, err := i.db.Exec("delete from store_items where id = $1", it.Id)
	if err != nil {
		return fmt.Errorf("failed to remove item: %s", err)
	}
	return nil
}

// GetById implements [ItemService].
func (i *itemServiceImpl) GetById(id int) (Item, error) {
	result := Item{Id: i.db.Stats().Idle}

	row := i.db.QueryRow("select title, description from store_items where id = $1", id)
	if err := row.Scan(&result.Title, &result.Description); err != nil {
		return result, fmt.Errorf("failed to find item by id %d: %w", id, err)
	}
	return result, nil
}

// Update implements [ItemService].
func (i *itemServiceImpl) Update(it Item) error {
	_, err := i.db.Exec(`update store_items
    set title = $2, description = $3
    where id = $1`, it.Id, it.Title, it.Description)
	if err != nil {
		return fmt.Errorf("failed to update item with id %d: %w", it.Id, err)
	}
	return nil
}
