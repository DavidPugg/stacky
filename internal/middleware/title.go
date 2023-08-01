package middleware

import (
	"fmt"

	"github.com/davidpugg/stacky/internal/utils"
	"github.com/gofiber/fiber/v2"
)

type PageDetails struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

var details = map[string]PageDetails{
	"/":         {"Home", "Home page description"},
	"/test":     {"Test", "Test page description"},
	"/login":    {"Login", "Login page description"},
	"/register": {"Register", "Register page description"},
}

func UpdatePageDetails(c *fiber.Ctx) error {
	pd := details[c.Path()]

	if pd == (PageDetails{}) {
		fmt.Println("No page details found for path: ", c.Path(), " please add it to internal/middleware/title.go")
	}

	c.Locals("PageDetails", pd)

	utils.SetTrigger(c, utils.Trigger{
		Name: "updatePageDetails",
		Data: pd,
	})

	return c.Next()
}
