package handlers

import (
	"github.com/davidpugg/stacky/internal/middleware"
	"github.com/davidpugg/stacky/internal/utils"
	"github.com/gofiber/fiber/v2"
)

func (h *Handlers) registerViewRoutes(c *fiber.App) {
	c.Get("/", middleware.UpdatePageDetails, h.renderMain)
	c.Get("/login", middleware.UpdatePageDetails, h.renderLogin)
	c.Get("/register", middleware.UpdatePageDetails, h.renderRegister)
}

func (h *Handlers) renderMain(c *fiber.Ctx) error {
	posts, err := h.data.GetPosts()
	if err != nil {
		return utils.RenderError(c, fiber.StatusInternalServerError, "Error fetching posts")
	}

	return utils.RenderPage(c, "index", fiber.Map{"Posts": posts})
}

func (h *Handlers) renderLogin(c *fiber.Ctx) error {
	return utils.RenderPage(c, "login", fiber.Map{})
}

func (h *Handlers) renderRegister(c *fiber.Ctx) error {
	return utils.RenderPage(c, "register", fiber.Map{})
}
