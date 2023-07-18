package main

import (
	"fmt"
	"log"

	"github.com/davidpugg/stacky/handlers"
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

	handlers.New().RegisterRoutes(app)

	fmt.Println("Server is running on port 3000")
	log.Fatal(app.Listen(":3000"))
}
