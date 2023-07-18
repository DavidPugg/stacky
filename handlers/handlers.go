package handlers

import (
	"github.com/davidpugg/stacky/utils"
	"github.com/gofiber/fiber/v2"
)

type Handlers struct{}

func New() *Handlers {
	return &Handlers{}
}

func (h *Handlers) RegisterRoutes(c *fiber.App) {
	c.Get("/", h.index)
	c.Post("/clicked", h.clicked)
}

func (h *Handlers) index(c *fiber.Ctx) error {
	return c.Render("index", fiber.Map{
		"Title": "Hello, World!",
	})
}

func (h *Handlers) clicked(c *fiber.Ctx) error {
	return utils.RenderPartial(c, "test", "Hello, World!")
}
