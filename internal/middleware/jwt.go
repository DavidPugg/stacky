package middleware

import (
	"github.com/davidpugg/stacky/internal/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type UserTokenData struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

func ParseToken(c *fiber.Ctx) error {
	cookie := c.Cookies("jwt")

	if cookie == "" {
		return c.Next()
	}

	token, err := utils.ValidateToken(cookie)
	if err != nil {
		return c.Next()
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return c.Next()
	}

	data := UserTokenData{
		ID:       int(claims["id"].(float64)),
		Username: claims["username"].(string),
		Email:    claims["email"].(string),
	}

	utils.SetTrigger(c, utils.Trigger{
		Name: "setUser",
	})

	c.Locals("User", data)

	return c.Next()
}