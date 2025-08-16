package products

import (
	"github.com/go-chi/chi/v5"
)

type ProductsManager struct {
	repo *ProductRepository
	apis *ProductAPIs
}

func NewProductsManager(mux *chi.Mux) *ProductsManager {
	repo := NewProductsRepository()
	apis := NewProductApis(mux, repo)
	return &ProductsManager{
		repo: repo,
		apis: apis,
	}
}
