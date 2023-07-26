package handlers

import (
	"github.com/davidpugg/stacky/data"
	"github.com/gofiber/fiber/v2"
)

type Handlers struct {
	data *data.Data
}

func New(data *data.Data) *Handlers {
	return &Handlers{data: data}
}

func (h *Handlers) RegisterRoutes(c *fiber.App) {
	h.registerViewRoutes(c)
	h.registerUtilRoutes(c)
	h.registerTodoRoutes(c)
}
