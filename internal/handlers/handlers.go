package handlers

import (
	"github.com/davidpugg/stacky/internal/data"
	"github.com/gofiber/fiber/v2"
)

type Handlers struct {
	mediaEndpoint string
	data          *data.Data
}

func New(data *data.Data) *Handlers {
	return &Handlers{
		data:          data,
		mediaEndpoint: "/media",
	}
}

func (h *Handlers) RegisterRoutes(c *fiber.App) {
	h.registerAuthRoutes(c)
	h.registerPostRoutes(c)
	h.registerUserRoutes(c)
	h.registerMediaRoutes(c)

	h.registerViewRoutes(c) //Must be last
}
