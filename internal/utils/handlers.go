package utils

import (
	"encoding/json"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

type Trigger struct {
	Name string
	Data interface{}
}

func RenderPartial(c *fiber.Ctx, view string, data interface{}) error {
	return c.Render(fmt.Sprintf("partials/%s", view), data, "layouts/empty")
}

func RenderPage(c *fiber.Ctx, view string, data interface{}, layout ...string) error {
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

func RenderError(c *fiber.Ctx, status int, details string) error {
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
	return RenderPage(c, "error", fiber.Map{
		"Status":  status,
		"Message": message,
		"Details": details,
	})
}

func SetTrigger(c *fiber.Ctx, triggers ...Trigger) error {
	if len(triggers) == 0 {
		return nil
	}

	alertMap := make(map[string]interface{})
	for _, t := range triggers {
		alertMap[t.Name] = t.Data
	}

	alert, err := json.Marshal(alertMap)
	if err != nil {
		return RenderError(c, fiber.StatusInternalServerError, "Error setting trigger")
	}

	c.Set("HX-Trigger", string(alert))
	return nil
}

func SetAlert(c *fiber.Ctx, status int, message string) error {
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
	if err := SetTrigger(c, Trigger{
		Name: "showAlert",
		Data: value,
	}); err != nil {
		return err
	}

	return nil
}

func SendAlert(c *fiber.Ctx, status int, message string) error {
	if err := SetAlert(c, status, message); err != nil {
		return err
	}

	c.Set("HX-Reswap", "none")
	return c.SendStatus(status)
}