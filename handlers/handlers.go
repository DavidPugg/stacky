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

func renderError(c *fiber.Ctx, status int, details string, nonRender ...bool) error {
	var message string
	switch status {
	case fiber.StatusInternalServerError:
		message = "Internal server error"
	case fiber.StatusNotFound:
		message = "Page not found"
	default:
		message = "Something went wrong"
	}

	if len(nonRender) > 0 && nonRender[0] {
		c.Set("HX-Redirect", fmt.Sprintf("/error?status=%d&message=%s&details=%s", status, message, details))
		return c.SendStatus(status)
	}

	return c.Redirect(fmt.Sprintf("/error?status=%d&message=%s&details=%s", status, message, details))
}

func setTrigger(c *fiber.Ctx, trigger string, message interface{}) error {
	alert, err := json.Marshal(fiber.Map{trigger: message})
	if err != nil {
		return renderError(c, fiber.StatusInternalServerError, "Could not marshal alert", true)
	}

	c.Set("HX-Trigger", string(alert))
	return nil
}

func setAlert(c *fiber.Ctx, status int, message string) error {
	value := fiber.Map{
		"status":  status,
		"message": message,
	}

	c.Status(status)
	err := setTrigger(c, "showAlert", value)
	if err != nil {
		return err
	}
	return nil
}

func sendAlert(c *fiber.Ctx, status int, message string) error {
	if err := setAlert(c, fiber.StatusBadRequest, message); err != nil {
		return err
	}
	return c.SendStatus(fiber.StatusBadRequest)
}
