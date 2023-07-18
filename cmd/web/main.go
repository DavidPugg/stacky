package main

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
)

func main() {
	engine := html.New("./views", ".gotmpl")

	app := fiber.New(fiber.Config{
		Views:             engine,
		ViewsLayout:       "layouts/main",
		PassLocalsToViews: true,
	})

	app.Static("/public", "./public")

	app.Get("/", func(c *fiber.Ctx) error {
		c.Locals("Test", "Test Local")
		return c.Render("index", fiber.Map{
			"Title": "Hello, World!",
		})
	})

	app.Post("/clicked", func(c *fiber.Ctx) error {
		return renderPartial(c, "test", "Hello, World!")

	})

	fmt.Println("Server is running on port 3000")
	log.Fatal(app.Listen(":3000"))
}

func renderPartial(c *fiber.Ctx, view string, data interface{}) error {
	return c.Render(fmt.Sprintf("partials/%s", view), data, "layouts/empty")
}
