package routes

import (
	"felix1234567890/go-trello/database"
	"felix1234567890/go-trello/handlers"
	"felix1234567890/go-trello/repository"
	"felix1234567890/go-trello/service"

	"github.com/gofiber/fiber/v2"
)

func SetupUserRoutes(app fiber.Router) {
	userRepository := repository.NewUserRepository(database.DB)
	userService := service.NewUserService(userRepository)
	userHandler := handlers.NewUserHandler(userService)
	app.Get("/", userHandler.GetUsers)
	app.Get("/:id", userHandler.GetUserById)
	app.Delete("/:id", userHandler.DeleteUser)
	app.Put("/:id", userHandler.UpdateUser)
	app.Post("/", userHandler.CreateUser)
}
