package routes

import (
	"felix1234567890/go-trello/handlers"
	"felix1234567890/go-trello/middlewares"
	"felix1234567890/go-trello/repository"
	"felix1234567890/go-trello/service"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func SetupUserRoutes(app fiber.Router, db *gorm.DB) {
	userRepository := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepository)
	userHandler := handlers.NewUserHandler(userService)

	app.Get("/", userHandler.GetUsers)
	app.Get("/me", middlewares.DeserializeUser(db), userHandler.GetMe)
	app.Get("/:id", userHandler.GetUserById)
	app.Delete("/:id", userHandler.DeleteUser)
	app.Put("/:id", userHandler.UpdateUser)
	app.Post("/", userHandler.CreateUser)
	app.Post("/login", userHandler.Login)

}
