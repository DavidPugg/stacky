package middleware

import (
	"fmt"

	"github.com/davidpugg/stacky/internal/utils"
	"github.com/gofiber/fiber/v2"
)

var titles = map[string]string{
	"/":     "Todos",
	"/test": "Test",
}

func UpdateTitle(c *fiber.Ctx) error {
	title := titles[c.Path()]

	if title == "" {
		fmt.Println("No title found for path: ", c.Path(), " please add it to internal/middleware/title.go")
	}

	c.Locals("PageTitle", title)

	utils.SetTrigger(c, utils.Trigger{
		Name: "updateTitle",
		Data: title,
	})

	return c.Next()
}
