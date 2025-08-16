package products

import (
	"database/sql"
	"fmt"
	"gotemplate/infra/postgres"
)

type ProductRepository struct{}

func NewProductsRepository() *ProductRepository {
	return &ProductRepository{}
}

// GetProduct fetches a single product by ID
func (r *ProductRepository) GetProduct(id int64) (*Product, error) {
	var product Product
	query := `SELECT id, title, description, original_price, discount, thumbnail role FROM products WHERE id = $1`
	err := postgres.DB.QueryRow(query, id).Scan(&product.ID, &product.Title, &product.OriginalPrice, &product.Discount, &product.Thumbnail)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No user found
		}
		return nil, fmt.Errorf("GetUser: %w", err)
	}
	return &product, nil
}

// ListProducts
func (r *ProductRepository) GetProducts() ([]*Product, error) {
	query := `SELECT id, title, description, original_price, discount, thumbnail role FROM products`
	rows, err := postgres.DB.Query(query)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No user found
		}
		return nil, fmt.Errorf("Get Products: %w", err)
	}

	products := []*Product{}
	for rows.Next() {
		var product Product
		if err := rows.Scan(&product.ID, &product.Title, &product.OriginalPrice, &product.Discount, &product.Thumbnail); err != nil {
			return nil, err
		}
		products = append(products, &product)
	}
	return products, nil
}

// CreateProduct inserts a new product into the DB
func (r *ProductRepository) CreateProduct(p *Product) (*Product, error) {
	query := `INSERT INTO products (title, description, original_price, discount) VALUES ($1, $2, $3, $4, $5) RETURNING id`
	err := postgres.DB.QueryRow(query, *p.Title, *p.Description, *p.OriginalPrice, *p.Discount).Scan(&p.ID)
	if err != nil {
		return nil, fmt.Errorf("CreateProduct: %w", err)
	}
	return p, nil
}

// UpdateProduct updates a product into the DB
func (r *ProductRepository) UpdateProduct(p *Product) (*Product, error) {
	setClause := []string{}
	if p.Title!=nil {
		setClause = append(setClause, fmt.Sprintf("title=$%d", len(setClause)+1))
	}
	if p.Description!=nil {
		setClause = append(setClause, fmt.Sprintf("description=$%d", len(setClause)+1))
	}
	if p.OriginalPrice!=nil {
		setClause = append(setClause, fmt.Sprintf("original_price=$%d", len(setClause)+1))
	}
	if p.Discount!=nil {
		setClause = append(setClause, fmt.Sprintf("discount=$%d", len(setClause)+1))
	}
	if p.Thumbnail!=nil {
		setClause = append(setClause, fmt.Sprintf("thumbnail=$%d", len(setClause)+1))
	}
	query := `UPDATE products set title=$1, description=$2, original_price=$3, discount=$4, thumbnail=$5) VALUES ($1, $2, $3, $4, $5) where id=$1`
	err := postgres.DB.QueryRow(query, *p.Title, *p.Description, *p.OriginalPrice, *p.Discount, *p.Thumbnail)
	if err != nil {
		return nil, fmt.Errorf("CreateProduct: %w", err)
	}
	return p, nil
}
