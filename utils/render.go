package utils

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func RenderPartial(c *fiber.Ctx, view string, data interface{}) error {
	return c.Render(fmt.Sprintf("partials/%s", view), data, "layouts/empty")
}
