package middleware

import (
	"github.com/davidpugg/stacky/internal/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

type UserTokenData struct {
	ID            int    `json:"id"`
	Avatar        string `json:"avatar"`
	Username      string `json:"username"`
	Email         string `json:"email"`
	Authenticated bool   `json:"authenticated"`
}

func NewUserTokenData(id int, avatar, username, email string) *UserTokenData {
	return &UserTokenData{
		ID:            id,
		Avatar:        utils.CreateImagePath(avatar),
		Username:      username,
		Email:         email,
		Authenticated: true,
	}
}

func CheckSession(s *session.Store) fiber.Handler {
	return func(c *fiber.Ctx) error {
		session, err := s.Get(c)
		if err != nil {
			return utils.SendAlert(c, fiber.StatusInternalServerError, "Error getting session")
		}

		if session.Get("id") == nil {
			c.Locals("AuthUser", &UserTokenData{})
			return c.Next()
		}

		data := NewUserTokenData(
			session.Get("id").(int),
			session.Get("avatar").(string),
			session.Get("name").(string),
			session.Get("email").(string),
		)

		c.Locals("AuthUser", data)

		return c.Next()
	}
}

func Authenticate(c *fiber.Ctx) error {
	if c.Locals("AuthUser").(*UserTokenData).ID == 0 {
		return utils.SendAlert(c, 401, "You must be logged in to do that.")
	}

	return c.Next()
}
