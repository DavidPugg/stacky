package handlers

import (
	"github.com/davidpugg/stacky/utils"
	"github.com/gofiber/fiber/v2"
)

func (h *Handlers) registerUtilRoutes(c *fiber.App) {
	r := c.Group("/utils")
	r.Get("/showAlert", h.showAlert)
}

func (h *Handlers) showAlert(c *fiber.Ctx) error {
	t := c.Query("type")
	m := c.Query("message")

	c.Set("HX-Retarget", "#alert")
	return utils.RenderPartial(c, "alert", fiber.Map{
		"Type":    t,
		"Message": m,
	})
}
