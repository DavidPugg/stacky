package main

import (
	"fmt"
	"log"

	"github.com/davidpugg/stacky/handlers"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	"github.com/spf13/viper"
)

func main() {

	//Config
	viper.SetConfigFile(".env")
	viper.ReadInConfig()

	//Fiber
	engine := html.New("./views", ".gotmpl")

	app := fiber.New(fiber.Config{
		Views:             engine,
		ViewsLayout:       "layouts/main",
		PassLocalsToViews: true,
	})

	app.Static("/public", "./public")

	//Routes
	handlers.New().RegisterRoutes(app)

	port := viper.GetInt("PORT")

	fmt.Printf("Server is running on port %d", port)
	log.Fatal(app.Listen(fmt.Sprintf(":%d", port)))
}
