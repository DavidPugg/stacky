package middleware

import "github.com/gofiber/fiber/v2"

func MainAuthGuard(c *fiber.Ctx) error {
	if c.Locals("AuthUser").(*UserTokenData).ID == 0 {
		return c.Redirect("/login")
	}

	return c.Next()
}

func LoginAuthGuard(c *fiber.Ctx) error {
	if c.Locals("AuthUser").(*UserTokenData).ID != 0 {
		return c.Redirect("/")
	}

	return c.Next()
}
