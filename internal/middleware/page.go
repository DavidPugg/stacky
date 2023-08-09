package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
)

func SamePage(c *fiber.Ctx) error {
	if c.Get("HX-Boosted") == "true" {
		r := "/" + strings.Split(c.Get("Referer"), "/")[3]
		if r == c.Path() {
			c.Set("HX-Reswap", "none show:no-scroll")
			return c.SendString("Already on page")
		}
	}

	return c.Next()
}

func MainAuthGuard(c *fiber.Ctx) error {
	if c.Locals("AuthUser").(*UserTokenData).ID == 0 {
		return c.Redirect("/discover")
	}

	return c.Next()
}

func LoginAuthGuard(c *fiber.Ctx) error {
	if c.Locals("AuthUser").(*UserTokenData).ID != 0 {
		return c.Redirect("/")
	}

	return c.Next()
}
