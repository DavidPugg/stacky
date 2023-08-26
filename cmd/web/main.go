package main

import (
	"fmt"
	"log"

	"github.com/davidpugg/stacky/internal/data"
	"github.com/davidpugg/stacky/internal/handlers"
	"github.com/davidpugg/stacky/internal/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/storage/sqlite3"
	"github.com/gofiber/template/html/v2"
	"github.com/spf13/viper"
)

func main() {
	//Config
	viper.SetConfigFile(".env")
	viper.SetConfigType("env")
	viper.AutomaticEnv()
	viper.ReadInConfig()

	//Fiber
	engine := html.New("./views", ".html")

	app := fiber.New(fiber.Config{
		Views:             engine,
		ViewsLayout:       "layouts/main",
		PassLocalsToViews: true,
		BodyLimit:         20 * 1024 * 1024,
	})

	app.Static("/public", "./public")
	app.Static("/cropperjs", "./node_modules/cropperjs/dist")
	app.Static("/alpinejs", "./node_modules/alpinejs/dist")

	//Data
	db := data.DBconnect()
	defer db.Close()

	ss := sqlite3.New()
	sessionStore := session.New(session.Config{
		Storage: ss,
	})

	dataInstance := data.New(db)

	//Middleware
	app.Use(cors.New(cors.Config{
		AllowCredentials: true,
		AllowOrigins:     "*",
	}))

	app.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed,
	}))

	app.Use(middleware.CheckSession(sessionStore))
	app.Use(middleware.SamePage)

	//Routes
	handlers.New(dataInstance, sessionStore).RegisterRoutes(app)

	//Server
	port := viper.GetInt("PORT")

	fmt.Printf("Server is running on port %d", port)
	log.Fatal(app.Listen(fmt.Sprintf(":%d", port)))
}
