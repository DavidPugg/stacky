package handlers

import "github.com/gofiber/fiber/v2"

func (h *Handlers) registerErrorRoutes(c *fiber.App) {
	c.Get("/404", page404)
	c.Get("/500", page500)
}

func page404(c *fiber.Ctx) error {
	return renderPage(c, "error/404", nil)
}

func page500(c *fiber.Ctx) error {
	return renderPage(c, "error/500", nil)
}
