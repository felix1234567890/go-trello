package routes

import (
	"felix1234567890/go-trello/database"
	"felix1234567890/go-trello/handlers"
	"felix1234567890/go-trello/middlewares"
	"felix1234567890/go-trello/repository"
	"felix1234567890/go-trello/service"

	"github.com/gofiber/fiber/v2"
)

func SetupUserRoutes(app fiber.Router) {
	userRepository := repository.NewUserRepository(database.DB)
	userService := service.NewUserService(userRepository)
	userHandler := handlers.NewUserHandler(userService)

	userRoutes := app.Group("/users")
	userRoutes.Get("/", userHandler.GetUsers)
	userRoutes.Get("/me", middlewares.DeserializeUser, userHandler.GetMe)
	userRoutes.Get("/:id", userHandler.GetUserById)
	userRoutes.Delete("/:id", userHandler.DeleteUser)
	userRoutes.Put("/:id", userHandler.UpdateUser)
	userRoutes.Post("/", userHandler.CreateUser)
	userRoutes.Post("/login", userHandler.Login)

}
