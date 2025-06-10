package main

import (
	"felix1234567890/go-trello/database"
	"felix1234567890/go-trello/handlers"
	"felix1234567890/go-trello/repository"
	"felix1234567890/go-trello/routes"
	"felix1234567890/go-trello/utils"
	"flag"
	"log"

	_ "felix1234567890/go-trello/docs"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/swagger"
	"github.com/joho/godotenv"
)

const defaultPort = "3000"

// @title			Go-Trello API
// @version		1.0
// @description	This is a sample swagger for Fiber
// @termsOfService	http://swagger.io/terms/
// @contact.name	API Support
// @contact.email	fiber@swagger.io
// @license.name	Apache 2.0
// @license.url	http://www.apache.org/licenses/LICENSE-2.0.html
// @host			localhost:3000
// @BasePath		/
func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}
	utils.InitSecretKey()
	db, err := database.ConnectToDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	port := flag.String("port", defaultPort, "server port")
	flag.Parse()
	app := fiber.New()
	app.Use(logger.New())
	app.Use(cors.New())
	app.Get("/swagger/*", swagger.HandlerDefault)
	globalPrefix := app.Group("/api")
	userRoutes := globalPrefix.Group("/users")
	routes.SetupUserRoutes(userRoutes, db)

	groupRepo := repository.NewGroupRepository(db)
	groupHandler := handlers.NewGroupHandler(groupRepo)
	routes.SetupGroupRoutes(globalPrefix, groupHandler)

	log.Fatal(app.Listen(":" + *port))
}
