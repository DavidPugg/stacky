package handlers

import (
	"github.com/davidpugg/stacky/internal/data"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

type Handlers struct {
	mediaEndpoint string
	data          *data.Data
	session       *session.Store
}

func New(data *data.Data, s *session.Store) *Handlers {
	return &Handlers{
		mediaEndpoint: "/media",
		data:          data,
		session:       s,
	}
}

func (h *Handlers) RegisterRoutes(c *fiber.App) {
	h.registerAuthRoutes(c)
	h.registerPostRoutes(c)
	h.registerUserRoutes(c)
	h.registerMediaRoutes(c)
	h.registerViewRoutes(c)
}
