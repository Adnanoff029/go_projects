package main

import (
	"fmt"
	"log"
	"os"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/Adnanoff029/url-shortener/api/routes"
)

func SetupRoutes(app *fiber.App) {
	app.Get("/:url", routes.ResolveURL)
	app.Post("/api/v1", routes.ShortenURL)
}

func main(){
	err := godotenv.Load()
	if err != nil{
		fmt.Println(err)
	}

	app := fiber.New()

	app.Use(logger.New())
	SetupRoutes(app)
	log.Fatal(app.Listen(os.Getenv("APP_PORT")))
}