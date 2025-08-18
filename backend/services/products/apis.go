package products

import (
	"encoding/json"
	"gotemplate/utility/logger"
	middleware "gotemplate/utility/middlewares"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type ProductAPIs struct {
	repo *ProductRepository
}

func NewProductApis(mux *chi.Mux, repo *ProductRepository) *ProductAPIs {
	pa := &ProductAPIs{repo: repo}
	pa.registerProductRoutes(mux)
	return pa
}

func (pa *ProductAPIs) registerProductRoutes(mux *chi.Mux) {
	mux.Route("/products", func(r chi.Router) {
		r.Group(func(protected chi.Router) {
			protected.Use(middleware.JWTAuth)
			r.Get("/", middleware.JSON(pa.getProducts))
			r.Get("/", middleware.JSON(pa.getProduct))
		})
	})
}

func (pa *ProductAPIs) getProducts(w http.ResponseWriter, r *http.Request) (interface{}, int) {
	logger.Logger.Infof("Starting to get products")
	products, err := pa.repo.GetProducts()
	if err != nil {
		logger.Logger.Errorf("Ran into an error %v", err)
		return nil, http.StatusBadRequest
	}
	return products, http.StatusOK
}

func (pa *ProductAPIs) getProduct(w http.ResponseWriter, r *http.Request) (interface{}, int) {
	logger.Logger.Infof("Starting to get product")
	var input GetProductsInputDto

	// Decode JSON request into struct
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		return map[string]string{"error": "Invalid request body"}, http.StatusBadRequest
	}
	product, err := pa.repo.GetProduct(input.Id)
	if err != nil {
		logger.Logger.Errorf("Ran into an error %v", err)
		return nil, http.StatusBadRequest
	}
	return product, http.StatusOK
}

func (pa *ProductAPIs) createProduct(w http.ResponseWriter, r *http.Request) (interface{}, int) {
	logger.Logger.Infof("Creating Product")
	var input CreateProductsInputDto

	// Decode JSON request into struct
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		return map[string]string{"error": "Invalid request body"}, http.StatusBadRequest
	}
	createProductInput := Product{
		Title:         input.Title,
		Description:   input.Description,
		OriginalPrice: input.OriginalPrice,
		Discount:      input.Discount,
	}
	product, err := pa.repo.CreateProduct(&createProductInput)
	if err != nil {
		return map[string]string{"error": "Something went wrong"}, http.StatusInternalServerError
	}
	jsonData, err := json.Marshal(product)
	if err != nil {
		return nil, http.StatusInternalServerError
	}

	return map[string]interface{}{"data": jsonData}, http.StatusCreated
}

func (pa *ProductAPIs) updateProduct(w http.ResponseWriter, r *http.Request) (interface{}, int) {
	logger.Logger.Infof("Updating Product")
	var input UpdateProductsInputDto

	// Decode JSON request into struct
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		return map[string]string{"error": "Invalid request body"}, http.StatusBadRequest
	}
	createProductInput := Product{
		ID:            input.ID,
		Title:         input.Title,
		Description:   input.Description,
		OriginalPrice: input.OriginalPrice,
		Discount:      input.Discount,
		Thumbnail:     input.Thumbnail,
	}
	product, err := pa.repo.CreateProduct(&createProductInput)
	if err != nil {
		return map[string]string{"error": "Something went wrong"}, http.StatusInternalServerError
	}
	jsonData, err := json.Marshal(product)
	if err != nil {
		return nil, http.StatusInternalServerError
	}

	return map[string]interface{}{"data": jsonData}, http.StatusCreated
}
