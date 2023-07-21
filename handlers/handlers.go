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

func sendError(c *fiber.Ctx, status int, details string) error {
	var message string
	switch status {
	case fiber.StatusInternalServerError:
		message = "Internal server error"
	case fiber.StatusNotFound:
		message = "Page not found"
	default:
		message = "Something went wrong"
	}

	c.Set("HX-Redirect", fmt.Sprintf("/error?status=%d&message=%s&details=%s", status, message, details))
	return c.SendStatus(status)
}

func sendTrigger(c *fiber.Ctx, status int, trigger string, message interface{}) error {
	alert, err := json.Marshal(fiber.Map{trigger: message})
	if err != nil {
		return sendError(c, fiber.StatusInternalServerError, "Could not marshal alert")
	}

	c.Set("HX-Trigger", string(alert))
	return c.SendStatus(status)
}

func showAlert(c *fiber.Ctx, status int, message string) error {
	value := fiber.Map{
		"status":  status,
		"message": message,
	}

	return sendTrigger(c, status, "showAlert", value)
}
