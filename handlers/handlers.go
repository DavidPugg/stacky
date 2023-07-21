package handlers

import (
	"encoding/json"
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
	c.Get("/error", func(c *fiber.Ctx) error {
		return renderPage(c, "error", nil)
	})

	h.registerErrorRoutes(c)
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

	return c.Render(fmt.Sprintf("%s", view), data, l)
}

func sendTrigger(c *fiber.Ctx, status int, trigger string, message string) error {
	alert, err := json.Marshal(fiber.Map{trigger: message})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Something went wrong")
	}

	c.Status(status).Set("HX-Trigger", string(alert))
	return c.SendString(message)
}

func showAlert(c *fiber.Ctx, status int, message string) error {
	return sendTrigger(c, status, "showAlert", message)
}
