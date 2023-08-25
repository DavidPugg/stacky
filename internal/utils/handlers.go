package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"text/template"

	"github.com/gofiber/fiber/v2"
)

type Trigger struct {
	Name string
	Data interface{}
}

type PageDetails struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

func SetPartial(c *fiber.Ctx, view string, data interface{}) error {
	content, err := readTemplateFile(generatePartialsURL(view))
	if err != nil {
		return err
	}

	tmpl, err := parseTemplate(content, data)
	if err != nil {
		return err
	}

	if c.Locals("Partials") == nil {
		c.Locals("Partials", tmpl+"\n")
	} else {
		c.Locals("Partials", c.Locals("Partials").(string)+tmpl+"\n")
	}

	return nil
}

func RenderPartial(c *fiber.Ctx, view string, data interface{}) error {
	err := c.Render(fmt.Sprintf("partials/%s", view), data, "layouts/empty")
	if err != nil {
		return err
	}

	return renderPartials(c)
}

func RenderPage(c *fiber.Ctx, view string, data interface{}, pd *PageDetails, layout ...string) error {
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

	c.Locals("PageDetails", pd)
	if err := SetTrigger(c, "default", Trigger{
		Name: "updatePageDetails",
		Data: pd,
	}); err != nil {
		return RenderError(c, fiber.StatusInternalServerError, "Error setting page details")
	}

	err := c.Render(view, data, l)
	if err != nil {
		return RenderError(c, fiber.StatusInternalServerError, "Error rendering page")
	}

	return renderPartials(c)
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
	}, &PageDetails{
		Title:       "Error",
		Description: message,
	})
}

func SetTrigger(c *fiber.Ctx, t string, triggers ...Trigger) error {
	if len(triggers) == 0 {
		return nil
	}

	if t == "settle" {
		t = "HX-Trigger-After-Settle"
	} else if t == "swap" {
		t = "HX-Trigger-After-Swap"
	} else if t == "default" {
		t = "HX-Trigger"
	}

	alertMap := make(map[string]interface{})
	for _, t := range triggers {
		alertMap[t.Name] = t.Data
	}

	alert, err := json.Marshal(alertMap)
	if err != nil {
		return RenderError(c, fiber.StatusInternalServerError, "Error setting trigger")
	}

	if hxTrigger := c.GetRespHeader(t); hxTrigger != "" {
		var hxTriggerMap map[string]interface{}
		if err := json.Unmarshal([]byte(hxTrigger), &hxTriggerMap); err != nil {
			return RenderError(c, fiber.StatusInternalServerError, "Error setting trigger")
		}

		for k, v := range hxTriggerMap {
			alertMap[k] = v
		}

		alert, err = json.Marshal(alertMap)
		if err != nil {
			return RenderError(c, fiber.StatusInternalServerError, "Error setting trigger")
		}
	}

	c.Set(t, string(alert))

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
		"Type":    t,
		"Message": message,
	}

	c.Status(status)
	return SetPartial(c, "alert", value)
}

func SendAlert(c *fiber.Ctx, status int, message string) error {
	if err := SetAlert(c, status, message); err != nil {
		return err
	}

	c.Set("HX-Reswap", "none")
	return c.SendString(c.Locals("Partials").(string))
}

func SetRedirect(c *fiber.Ctx, url string) {
	r, _ := json.Marshal(fiber.Map{
		"redirect": url,
	})

	c.Set("HX-Trigger-After-Settle", string(r))
}

func readTemplateFile(filePath string) (string, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	templateContent := string(content)

	return templateContent, nil
}

func parseTemplate(templateContent string, data interface{}) (string, error) {
	tmpl, err := template.New("").Parse(templateContent)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

func generatePartialsURL(name string) string {
	return fmt.Sprintf("views/partials/%s.html", name)
}

func renderPartials(c *fiber.Ctx) error {
	if c.Locals("Partials") == nil {
		return nil
	}

	_, err := c.WriteString(c.Locals("Partials").(string))
	if err != nil {
		return err
	}

	return nil
}
