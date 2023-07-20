package handlers

import (
	"fmt"

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
	h.registerTodoRoutes(c)
}

func renderPartial(c *fiber.Ctx, view string, data interface{}) error {
	return c.Render(fmt.Sprintf("partials/%s", view), data, "layouts/empty")
}

func renderPage(c *fiber.Ctx, view string, data interface{}, layout ...string) error {
	var l string
	if len(layout) == 0 {
		l = "layouts/main"
	} else {
		l = layout[0]
	}

	if c.Locals("Error") != nil {
		l = "layouts/empty"
	}

	return c.Render(fmt.Sprintf("%s", view), data, l)
}

type Error struct {
	Status int
	Error  interface{}
}

func redirectWithError(c *fiber.Ctx, status int, err interface{}, handler func(c *fiber.Ctx) error) error {
	c.Status(status)
	c.Locals("Error", Error{Status: status, Error: err})
	return handler(c)
}
