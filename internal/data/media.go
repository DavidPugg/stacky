package data

import (
	"fmt"
	"image"
	"image/jpeg"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/google/uuid"
)

func (d *Data) CreateMediaLocally(img image.Image, header *multipart.FileHeader) (string, error) {
	randID := uuid.New().String()
	ext := filepath.Ext(header.Filename)

	imagePath := fmt.Sprintf("public/assets/images/%s%s", randID, ext)

	file, err := os.Create(imagePath)
	if err != nil {
		return "", err
	}

	defer file.Close()

	err = jpeg.Encode(file, img, nil)
	if err != nil {
		return "", err
	}

	return imagePath, nil
}
