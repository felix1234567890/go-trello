package main

import (
	"felix1234567890/go-trello/database"
	"felix1234567890/go-trello/routes"
	"flag"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
)

const defaultPort = "3000"

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}
	database.ConnectToDB()
	// utils.FakeUserFactory()
	port := flag.String("port", defaultPort, "server port")
	flag.Parse()
	app := fiber.New()
	app.Use(logger.New())
	app.Use(cors.New())
	globalPrefix := app.Group("/api")
	userRoutes := globalPrefix.Group("/users")
	routes.SetupUserRoutes(userRoutes)
	log.Fatal(app.Listen(":" + *port))
}
