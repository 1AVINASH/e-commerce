package images

import (
	"encoding/json"
	"gotemplate/utility/logger"
	middleware "gotemplate/utility/middlewares"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type ImageAPIs struct {
	repo *ImageRepository
	utils *ImageUtils
}

func NewImageApis(mux *chi.Mux, repo *ImageRepository) *ImageAPIs {
	pa := &ImageAPIs{repo: repo}
	pa.registerImageRoutes(mux)
	return pa
}

func (ia *ImageAPIs) registerImageRoutes(mux *chi.Mux) {
	mux.Route("/images", func(r chi.Router) {
		r.Group(func(protected chi.Router) {
			protected.Use(middleware.JWTAuth)
			r.Get("/", middleware.JSON(ia.GetProductPhotos))
			r.Post("/", middleware.JSON(ia.CreateImage))
			r.Post("/", middleware.JSON(ia.CreateImage))
		})
	})
}

func (ia *ImageAPIs) GetProductPhotos(w http.ResponseWriter, r *http.Request) (interface{}, int) {
	logger.Logger.Infof("Starting to get photos for product")
	var input GetProductPhotosInputDto

	// Decode JSON request into struct
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		return map[string]string{"error": "Invalid request body"}, http.StatusBadRequest
	}

	images, err := ia.repo.GetImagesForProduct(*input.ProductId)
	if err != nil {
		logger.Logger.Errorf("Ran into an error %v", err)
		return nil, http.StatusBadRequest
	}
	return images, http.StatusOK
}

func (ia *ImageAPIs) CreateImage(w http.ResponseWriter, r *http.Request) (interface{}, int) {
	logger.Logger.Infof("Creating Image")
	var input CreateProductPhotoInputDto

	// Decode JSON request into struct
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		return map[string]string{"error": "Invalid request body"}, http.StatusBadRequest
	}

	createImageInput := Image{
		ProductId: input.ProductId,
	}
	image, err := ia.repo.CreateImage(&createImageInput)
	if err != nil {
		return map[string]string{"error": "Something went wrong"}, http.StatusInternalServerError
	}

	return map[string]interface{}{"status": "created", "image": image}, http.StatusCreated
}

func (ia *ImageAPIs) UploadImage(w http.ResponseWriter, r *http.Request) (interface{}, int) {
	logger.Logger.Infof("Upload Image")
	r.Body = http.MaxBytesReader(w, r.Body, 10<<20) // 10 MB

	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return nil, http.StatusBadRequest
	}

	files := r.MultipartForm.File["files"]
	if len(files) == 0 {
		http.Error(w, "No files uploaded", http.StatusBadRequest)
		return nil, http.StatusBadRequest
	}

	var uploaded []map[string]string

	for _, fileHeader := range files {
		fileInfo, err := ia.utils.processFile(fileHeader)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return nil, http.StatusBadRequest
		}
		uploaded = append(uploaded, fileInfo)
	}


	return map[string]interface{}{"status": "created", "image": uploaded}, http.StatusCreated
}
