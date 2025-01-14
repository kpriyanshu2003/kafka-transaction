package main

import (
	"log"
	"os"

	"example.com/kafka/src/config"
	"example.com/kafka/src/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
)

func init() {
	godotenv.Load()
	config.SetupDatabase()
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT is not found.")
	}

	app := fiber.New()
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{AllowOrigins: "*"}))
		app.Route("/", routes.Routes)
	app.Listen(":" + port)
}
