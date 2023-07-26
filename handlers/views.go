package handlers

import (
	"github.com/gofiber/fiber/v2"
)

var titles = map[string]string{
	"/":     "Todos",
	"/test": "Test",
}

func (h *Handlers) registerViewRoutes(c *fiber.App) {
	c.Use(func(c *fiber.Ctx) error {
		title := titles[c.Path()]

		c.Locals("PageTitle", title)
		setTrigger(c, "updateTitle", title)
		return c.Next()
	})

	c.Get("/", h.renderMain)
	c.Get("/test", h.renderTest)
}

func (h *Handlers) renderMain(c *fiber.Ctx) error {
	todos, err := h.data.GetTodos()
	if err != nil {
		return renderError(c, fiber.StatusInternalServerError, "Error getting todos")
	}

	return renderPage(c, "index", fiber.Map{"Todos": todos})
}

func (h *Handlers) renderTest(c *fiber.Ctx) error {
	return renderPage(c, "test", fiber.Map{})
}
