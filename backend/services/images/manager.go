package images

import (
	"fmt"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/disintegration/imaging"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type ImageManager struct {
	repo *ImageRepository
	apis *ImageAPIs
}

func NewImageManager(mux *chi.Mux) *ImageManager {
	repo := NewImageRepository()
	apis := NewImageApis(mux, repo)
	return &ImageManager{
		repo: repo,
		apis: apis,
	}
}

const (
	fullDir  = "./uploads/full"
	thumbDir = "./uploads/thumbs"
)

var thumbSizes = map[string]int{
	"small":  150,
	"medium": 300,
	"large":  600,
}

func main() {
	r := chi.NewRouter()

	// Ensure folders exist
	os.MkdirAll(fullDir, os.ModePerm)
	for size := range thumbSizes {
		os.MkdirAll(filepath.Join(thumbDir, size), os.ModePerm)
	}

	r.Post("/upload", uploadHandler)

	fmt.Println("Server running on :8080")
	http.ListenAndServe(":8080", r)
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 10<<20) // 10 MB

	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	files := r.MultipartForm.File["files"]
	if len(files) == 0 {
		http.Error(w, "No files uploaded", http.StatusBadRequest)
		return
	}

	var uploaded []map[string]string

	for _, fileHeader := range files {
		fileInfo, err := processFile(fileHeader)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		uploaded = append(uploaded, fileInfo)
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, "%v", uploaded)
}

func processFile(fileHeader *multipart.FileHeader) (map[string]string, error) {
	// Validate MIME type
	mime := fileHeader.Header.Get("Content-Type")
	if mime != "image/jpeg" && mime != "image/png" {
		return nil, fmt.Errorf("unsupported file type: %s", mime)
	}

	file, err := fileHeader.Open()
	if err != nil {
		return nil, fmt.Errorf("cannot open file")
	}
	defer file.Close()

	// Generate unique filename
	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	if ext == "" {
		if mime == "image/jpeg" {
			ext = ".jpg"
		} else {
			ext = ".png"
		}
	}
	filename := uuid.New().String() + ext

	// Save original
	fullPath := filepath.Join(fullDir, filename)
	dst, err := os.Create(fullPath)
	if err != nil {
		return nil, fmt.Errorf("cannot save file")
	}
	defer dst.Close()

	_, err = file.Seek(0, 0)
	if err != nil {
		return nil, fmt.Errorf("seek error")
	}
	_, err = dst.ReadFrom(file)
	if err != nil {
		return nil, fmt.Errorf("save error")
	}

	// Async thumbnails
	go createThumbnails(fullPath, filename)

	return map[string]string{
		"full": fullPath,
	}, nil
}

func createThumbnails(fullPath, filename string) {
	// Add slight delay to avoid locking the file
	time.Sleep(200 * time.Millisecond)

	img, err := imaging.Open(fullPath)
	if err != nil {
		fmt.Println("Thumbnail generation failed:", err)
		return
	}

	for sizeName, size := range thumbSizes {
		go func(sizeName string, size int) {
			thumb := imaging.Thumbnail(img, size, size, imaging.Lanczos)
			thumbPath := filepath.Join(thumbDir, sizeName, filename)
			if err := imaging.Save(thumb, thumbPath); err != nil {
				fmt.Printf("Failed to save %s thumbnail: %v\n", sizeName, err)
			}
		}(sizeName, size)
	}
}
