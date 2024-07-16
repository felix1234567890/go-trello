package handlers

import (
	"errors"
	"felix1234567890/go-trello/models"
	"felix1234567890/go-trello/service"
	"felix1234567890/go-trello/utils"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type UserHandler struct {
	UserService service.UserService
}

func NewUserHandler(userService service.UserService) *UserHandler {
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
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(&fiber.Map{
				"message": "User with an id " + id + " was not found",
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"message": err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(&fiber.Map{
		"user": user,
	})
}

func (h *UserHandler) DeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")
	err := h.UserService.DeleteUser(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(&fiber.Map{
				"message": "User with an id " + id + " was not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"message": err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(&fiber.Map{
		"message": "User deleted successfully",
	})
}
func (h *UserHandler) UpdateUser(ctx *fiber.Ctx) error {
	var req *models.UpdateUserRequest
	id := ctx.Params("id")
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
		})
	}

	validationErrors := utils.ValidateRequest(req)
	if len(validationErrors) > 0 {
		// Return validation errors as a bad request
		errorMessages := make([]string, len(validationErrors))
		for i, validationErr := range validationErrors {
			errorMessages[i] = validationErr.Error()
		}
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"errors": errorMessages,
		})
	}
	err := h.UserService.UpdateUser(id, req)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ctx.Status(fiber.StatusNotFound).JSON(&fiber.Map{
				"message": "User with an id " + id + " could not be updated",
			})
		}
		return ctx.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"message": err.Error(),
		})
	}
	return ctx.Status(fiber.StatusOK).JSON(&fiber.Map{
		"message": "User updated successfully",
	})
}

func (h *UserHandler) CreateUser(ctx *fiber.Ctx) error {
	var req *models.CreateUserRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
		})
	}
	validationErrors := utils.ValidateRequest(req)

	if len(validationErrors) > 0 {
		errorMessages := make([]string, len(validationErrors))
		for i, validationErr := range validationErrors {
			errorMessages[i] = validationErr.Error()
		}
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"errors": errorMessages,
		})

	}
	user := req.ToUser()
	id, err := h.UserService.CreateUser(user)
	if err != nil {
		return utils.HandleErrorResponse(ctx, fiber.StatusInternalServerError, err.Error())
	}
	token, err := utils.CreateToken(id)
	return utils.JsonResponse(ctx, fiber.StatusCreated, fiber.Map{
		"token": token,
	})
}

func (h *UserHandler) Login(ctx *fiber.Ctx) error {
	var req *models.LoginUserRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
		})
	}
	validationErrors := utils.ValidateRequest(req)
	if len(validationErrors) > 0 {
		errorMessages := make([]string, len(validationErrors))
		for i, validationErr := range validationErrors {
			errorMessages[i] = validationErr.Error()
		}
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"errors": errorMessages,
		})
	}
	id, err := h.UserService.LoginUser(req)
	if err != nil {
		return utils.HandleErrorResponse(ctx, fiber.StatusInternalServerError, err.Error())
	}
	token, err := utils.CreateToken(id)
	return utils.JsonResponse(ctx, fiber.StatusCreated, fiber.Map{
		"token": token,
	})
}
func (h *UserHandler) GetMe(c *fiber.Ctx) error {
	user := c.Locals("user").(models.User)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": fiber.Map{"user": user}})
}
