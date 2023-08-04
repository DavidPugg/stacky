package middleware

import (
	"github.com/davidpugg/stacky/internal/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type UserTokenData struct {
	ID            int    `json:"id"`
	Username      string `json:"username"`
	Email         string `json:"email"`
	Authenticated bool   `json:"authenticated"`
}

func ParseToken(c *fiber.Ctx) error {
	c.Locals("AuthUser", &UserTokenData{})

	var t string
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		t = c.Cookies("jwt")
		if t == "" {
			return c.Next()
		}
	} else {
		t = authHeader[7:]
		if t == "" {
			return c.Next()
		}
	}

	token, err := utils.ValidateToken(t)
	if err != nil {
		return c.Next()
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return c.Next()
	}

	data := &UserTokenData{
		ID:            int(claims["id"].(float64)),
		Username:      claims["username"].(string),
		Email:         claims["email"].(string),
		Authenticated: true,
	}

	c.Locals("AuthUser", data)

	return c.Next()
}

func Authenticate(c *fiber.Ctx) error {
	if c.Locals("AuthUser").(*UserTokenData).ID == 0 {
		return utils.SendAlert(c, 401, "You must be logged in to do that.")
	}

	return c.Next()
}
