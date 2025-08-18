package images

import (
	"database/sql"
	"fmt"
	"gotemplate/infra/postgres"
)

type ImageRepository struct{}

func NewImageRepository() *ImageRepository {
	return &ImageRepository{}
}

// GetUser fetches a single image by ID
func (r *ImageRepository) GetImagesForProduct(productId int64) ([]*Image, error) {
	var images []*Image
	query := `SELECT id, product_id, path FROM product_photos WHERE product_id = $1`
	rows, err := postgres.DB.Query(query, productId)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No image found
		}
		return nil, fmt.Errorf("GetImages: %w", err)
	}
	for rows.Next() {
		img := Image{}
		err := rows.Scan(&img.ID, &img.ProductId, &img.Path)
		if err != nil {
			return images, err
		}
		images = append(images, &img)
	}
	return images, nil
}

// CreateUser inserts a new user into the DB
func (r *ImageRepository) CreateImage(i *Image) (*Image, error) {
	query := `INSERT INTO product_photos (product_id, path) VALUES ($1, $2) RETURNING id`
	err := postgres.DB.QueryRow(query, *i.ProductId, *i.Path).Scan(&i.ID)
	if err != nil {
		return nil, fmt.Errorf("CreateImage: %w", err)
	}
	return i, nil
}
