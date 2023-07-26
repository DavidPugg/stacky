package middleware

import (
	"github.com/davidpugg/stacky/internal/utils"
	"github.com/gofiber/fiber/v2"
)

var titles = map[string]string{
	"/":     "Todos",
	"/test": "Test",
}

func UpdateTitle(c *fiber.Ctx) error {
	title := titles[c.Path()]

	c.Locals("PageTitle", title)

	utils.SetTrigger(c, utils.Trigger{
		Name: "updateTitle",
		Data: title,
	})

	return c.Next()
}
