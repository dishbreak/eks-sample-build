package models

import (
	"database/sql"
	"fmt"
)

type Image struct {
	Id     int
	Path   string
	ItemId int
}

type ImageService interface {
	GetById(imageId int) (Image, error)
	GetByItemId(itemId int) ([]Image, error)
	Create(itemId int, img Image) (Image, error)
	DeleteById(itemId int) error
}

type imageServiceImpl struct {
	db *sql.DB
}

// Create implements [ImageService].
func (i imageServiceImpl) Create(itemId int, img Image) (Image, error) {
	result := Image{}
	if itemId != img.ItemId {
		return result, fmt.Errorf("failed to create image: passed in item id %d does not match item %x", itemId, img)
	}

	row := i.db.QueryRow(
		`INSERT INTO store_images
		 (id, path, item_id)
		 VALUES 
		 (DEFAULT, $1, $2)
		 RETURNING id, path, item_id`,
		img.Path, img.ItemId,
	)

	if err := row.Scan(&result.Id, &result.Path, &result.ItemId); err != nil {
		return result, fmt.Errorf("failed to create image: %w", err)
	}

	return result, nil
}

// DeleteById implements [ImageService].
func (i imageServiceImpl) DeleteById(itemId int) error {
	_, err := i.db.Exec(`DELETE from store_images WHERE id = $1`, itemId)
	if err != nil {
		return fmt.Errorf("failed to delete image: %w", err)
	}
	return nil
}

// GetById implements [ImageService].
func (i imageServiceImpl) GetById(imageId int) (Image, error) {
	result := Image{}
	row := i.db.QueryRow(`select id, path, item_id from store_images where id = $1`, imageId)
	if err := row.Scan(&result.Id, &result.Path, &result.ItemId); err != nil {
		return result, fmt.Errorf("failed to find image: %w", err)
	}

	return result, nil
}

// GetByItemId implements [ImageService].
func (i imageServiceImpl) GetByItemId(itemId int) ([]Image, error) {
	result := make([]Image, 0)

	rows, err := i.db.Query(`select id, path, item_id from store_images where item_id = $1`, itemId)
	if err != nil {
		return result, fmt.Errorf(`failed to find images for item_id %d: %w`, itemId, err)
	}

	for rows.Next() {
		img := Image{}
		if err := rows.Scan(&img.Id, &img.Path, &img.ItemId); err != nil {
			return result, fmt.Errorf("failed to read image record from db: %w", err)
		}
		result = append(result, img)
	}

	if rows.Err() != nil {
		return result, fmt.Errorf("failed to parse query results: %w", rows.Err())
	}

	return result, nil
}

func NewImageService(db *sql.DB) ImageService {
	return imageServiceImpl{
		db: db,
	}
}
