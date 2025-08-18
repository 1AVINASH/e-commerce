package images

import (
	"os"
	"path/filepath"

	"github.com/go-chi/chi/v5"
)

type ImageManager struct {
	repo  *ImageRepository
	apis  *ImageAPIs
	utils *ImageUtils
}

func NewImageManager(mux *chi.Mux) *ImageManager {
	repo := NewImageRepository()
	apis := NewImageApis(mux, repo)
	utils := NewImageUtils(repo)
	os.MkdirAll(fullDir, os.ModePerm)
	for size := range thumbSizes {
		os.MkdirAll(filepath.Join(thumbDir, size), os.ModePerm)
	}
	return &ImageManager{
		repo:  repo,
		apis:  apis,
		utils: utils,
	}
}
