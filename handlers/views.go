package handlers

import "github.com/gofiber/fiber/v2"

func (h *Handlers) registerViewRoutes(c *fiber.App) {
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
	return renderPage(c, "test", nil)
}
