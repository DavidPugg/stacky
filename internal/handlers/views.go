package handlers

import (
	"github.com/davidpugg/stacky/internal/middleware"
	"github.com/davidpugg/stacky/internal/utils"
	"github.com/gofiber/fiber/v2"
)

func (h *Handlers) registerViewRoutes(c *fiber.App) {
	c.Get("/", middleware.UpdateTitle, h.renderMain)
	c.Get("/test", middleware.UpdateTitle, h.renderTest)
}

func (h *Handlers) renderMain(c *fiber.Ctx) error {
	todos, err := h.data.GetTodos()
	if err != nil {
		return utils.RenderError(c, fiber.StatusInternalServerError, "Error getting todos")
	}

	return utils.RenderPage(c, "index", fiber.Map{"Todos": todos})
}

func (h *Handlers) renderTest(c *fiber.Ctx) error {
	return utils.RenderPage(c, "test", fiber.Map{})
}
