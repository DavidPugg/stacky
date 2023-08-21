package handlers

import (
	"os"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
)

func (h *Handlers) registerMediaRoutes(c *fiber.App) {
	r := c.Group(h.mediaEndpoint)
	r.Get("/:id", h.mediaShow)
}

func (h *Handlers) mediaShow(c *fiber.Ctx) error {
	id := c.Params("id")

	if id == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Missing image ID")
	}

	imagePath := filepath.Join("uploads", id)

	if _, err := os.Stat(imagePath); err != nil {
		if os.IsNotExist(err) {
			return fiber.NewError(fiber.StatusNotFound, "Image not found")
		}
		return fiber.NewError(fiber.StatusInternalServerError, "Internal server error")
	}

	contentType := "application/octet-stream"
	extension := filepath.Ext(imagePath)
	switch extension {
	case ".jpg", ".jpeg":
		contentType = "image/jpeg"
	case ".png":
		contentType = "image/png"
	case ".gif":
		contentType = "image/gif"
	}

	err := c.SendFile(imagePath, true)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Internal server error")
	}

	c.Set(fiber.HeaderContentType, contentType)
	c.Set(fiber.HeaderCacheControl, "public, max-age=3600")

	return nil
}
