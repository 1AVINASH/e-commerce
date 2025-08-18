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
	"github.com/google/uuid"
)

type ImageUtils struct {
	repo *ImageRepository
}

func NewImageUtils(repo *ImageRepository) *ImageUtils {
	return &ImageUtils{repo: repo}
}

func (iu *ImageUtils) processFile(fileHeader *multipart.FileHeader) (map[string]string, error) {
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
	go iu.createThumbnails(fullPath, filename)

	return map[string]string{
		"full": fullPath,
	}, nil
}

func (iu *ImageUtils) createThumbnails(fullPath, filename string) {
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
