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

	c.Get("/test", func(c *fiber.Ctx) error {
		return renderPage(c, "test", nil)
	})

	c.Get("/showAlert", func(c *fiber.Ctx) error {
		t := c.Query("type")
		m := c.Query("message")

		c.Set("HX-Retarget", "#alert")
		return renderPartial(c, "alert", fiber.Map{
			"Type":    t,
			"Message": m,
		})
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

	if c.Get("HX-Request") == "true" {
		l = "layouts/empty"

		c.Set("HX-Push-Url", c.Path())
		c.Set("HX-Reswap", "innerHTML show:no-scroll")
		c.Set("HX-Retarget", "#content")
	}

	return c.Render(fmt.Sprintf("%s", view), data, l)
}

func renderError(c *fiber.Ctx, status int, details string) error {
	var message string
	switch status {
	case fiber.StatusInternalServerError:
		message = "Internal server error"
	case fiber.StatusNotFound:
		message = "Page not found"
	default:
		message = "Something went wrong"
	}

	c.Status(status)
	c.Set("HX-Push-Url", "/error")
	return renderPage(c, "error", fiber.Map{
		"Status":  status,
		"Message": message,
		"Details": details,
	})
}

func setTrigger(c *fiber.Ctx, trigger string, value interface{}) error {
	alert, err := json.Marshal(fiber.Map{trigger: value})
	if err != nil {
		return renderError(c, fiber.StatusInternalServerError, "Could not marshal alert")
	}

	c.Set("HX-Trigger", string(alert))
	return nil
}

func setAlert(c *fiber.Ctx, status int, message string) error {
	var t string
	switch status / 100 {
	case 1:
		t = "info"
	case 2:
		t = "success"
	case 3:
		t = "info"
	default:
		t = "error"
	}

	value := fiber.Map{
		"type":    t,
		"message": message,
	}

	c.Status(status)
	if err := setTrigger(c, "showAlert", value); err != nil {
		return err
	}

	return nil
}

func sendAlert(c *fiber.Ctx, status int, message string) error {
	if err := setAlert(c, status, message); err != nil {
		return err
	}

	c.Set("HX-Reswap", "none")
	return c.SendStatus(status)
}
