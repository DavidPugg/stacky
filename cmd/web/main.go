package main

import (
	"fmt"
	"log"

	"github.com/davidpugg/stacky/internal/data"
	"github.com/davidpugg/stacky/internal/handlers"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/template/html/v2"
	"github.com/spf13/viper"
)

func main() {
	//Config
	viper.SetConfigFile(".env")
	viper.SetConfigType("env")
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Error reading config file, %s", err)
	}

	//Fiber
	engine := html.New("./views", ".gotmpl")

	app := fiber.New(fiber.Config{
		Views:             engine,
		ViewsLayout:       "layouts/main",
		PassLocalsToViews: true,
	})

	app.Static("/public", "./public")

	//Middleware

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	app.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed,
	}))

	//Data
	db := data.DBconnect()
	defer db.Close()

	data := data.New(db)

	//Routes
	handlers.New(data).RegisterRoutes(app)

	//Server
	port := viper.GetInt("PORT")

	fmt.Printf("Server is running on port %d", port)
	log.Fatal(app.Listen(fmt.Sprintf(":%d", port)))
}
