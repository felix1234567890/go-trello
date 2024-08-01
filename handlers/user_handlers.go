// Package handlers provides HTTP request handlers.
package handlers

import (
	"errors"
	"felix1234567890/go-trello/models"
	"felix1234567890/go-trello/service"
	"felix1234567890/go-trello/utils"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// UserHandler handles HTTP requests related to users.
type UserHandler struct {
	UserService service.UserService
}

// NewUserHandler creates a new UserHandler instance.
func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{
		UserService: userService,
	}
}

// GetUsers godoc
//
//	@Summary		Get all users
//	@Description	Get a list of all users
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Success		200	{array}		models.User	"List of users"
//	@Failure		500	{object}	fiber.Map	"Internal Server Error"
//	@Router			/users [get]
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

// GetUserById godoc
//
//	@Summary		Get a user by ID
//	@Description	Get details of a specific user
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string		true	"User ID"
//	@Success		200	{object}	models.User	"User details"
//	@Failure		404	{object}	fiber.Map	"User not found"
//	@Failure		500	{object}	fiber.Map	"Internal Server Error"
//	@Router			/users/{id} [get]
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

// DeleteUser godoc
//
//	@Summary		Delete a user
//	@Description	Delete a user by ID
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string		true	"User ID"
//	@Success		200	{object}	fiber.Map	"User deleted successfully"
//	@Failure		404	{object}	fiber.Map	"User not found"
//	@Failure		500	{object}	fiber.Map	"Internal Server Error"
//	@Router			/users/{id} [delete]
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

// UpdateUser godoc
//
//	@Summary		Update a user
//	@Description	Update a user's details
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string						true	"User ID"
//	@Param			user	body		models.UpdateUserRequest	true	"User details to update"
//	@Success		200		{object}	fiber.Map					"User updated successfully"
//	@Failure		400		{object}	fiber.Map					"Invalid request body or validation errors"
//	@Failure		404		{object}	fiber.Map					"User not found"
//	@Failure		500		{object}	fiber.Map					"Internal Server Error"
//	@Router			/users/{id} [put]
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

// CreateUser godoc
//
//	@Summary		Create a new user
//	@Description	Create a new user and return an authentication token
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			user	body		models.CreateUserRequest	true	"User details"
//	@Success		201		{object}	fiber.Map					"User created successfully, token returned"
//	@Failure		400		{object}	fiber.Map					"Invalid request body or validation errors"
//	@Failure		500		{object}	fiber.Map					"Internal Server Error"
//	@Router			/users [post]
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

// Login godoc
//
//	@Summary		User login
//	@Description	Authenticate a user and return a token
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			credentials	body		models.LoginUserRequest	true	"Login credentials"
//	@Success		200			{object}	fiber.Map				"Authentication successful, token returned"
//	@Failure		400			{object}	fiber.Map				"Invalid request body or validation errors"
//	@Failure		500			{object}	fiber.Map				"Internal Server Error"
//	@Router			/login [post]
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
	return utils.JsonResponse(ctx, fiber.StatusOK, fiber.Map{
		"token": token,
	})
}

// GetMe godoc
//
//	@Summary		Get current user
//	@Description	Get details of the currently authenticated user
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Success		200	{object}	fiber.Map	"Details of the authenticated user"
//	@Router			/me [get]
func (h *UserHandler) GetMe(c *fiber.Ctx) error {
	user := c.Locals("user").(models.User)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": fiber.Map{"user": user}})
}
