package routes

import "github.com/gofiber/fiber/v2"

func SetupUserRoutes(app fiber.Router) {
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Hello, World!",
		})
	})
}
