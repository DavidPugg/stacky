package data

import (
	"fmt"
	"mime/multipart"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func (d *Data) CreateMediaLocally(c *fiber.Ctx, image *multipart.FileHeader) (string, error) {
	randID := uuid.New().String()
	ext := filepath.Ext(image.Filename)

	serverURL := c.BaseURL()

	imagePath := fmt.Sprintf("public/assets/images/%s%s", randID, ext)

	err := c.SaveFile(image, imagePath)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s/%s", serverURL, imagePath), nil
}
