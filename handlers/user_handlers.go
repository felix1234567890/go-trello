package handlers

import (
	"felix1234567890/go-trello/service"

	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	UserService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{
		UserService: userService,
	}
}

func (h *UserHandler) GetUsers(c *fiber.Ctx) error {
	users, err := h.UserService.GetUsers()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"message": err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(&fiber.Map{
		"users": users,
	})
}

func (h *UserHandler) GetUserById(c *fiber.Ctx) error {
	id := c.Params("id")
	user, err := h.UserService.GetUserById(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"message": err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(&fiber.Map{
		"user": user,
	})
}
