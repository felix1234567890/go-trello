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
//	@Success		200	{object}	fiber.Map	"List of users (e.g., {\"users\": [{\"id\": 1, ...}]})"
//	@Failure		500	{object}	fiber.Map	"Internal Server Error (e.g., {\"message\": \"Failed to retrieve users\"})"
//	@Router			/users [get]
func (h *UserHandler) GetUsers(c *fiber.Ctx) error {
	users, err := h.UserService.GetUsers()
	if err != nil {
		return utils.HandleErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve users") // Generic message
	}
	return utils.JsonResponse(c, fiber.StatusOK, fiber.Map{"users": users})
}

// GetUserById godoc
//
//	@Summary		Get a user by ID
//	@Description	Get details of a specific user
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string		true	"User ID"
//	@Success		200	{object}	fiber.Map	"User details (e.g., {\"user\": {\"id\": 1, ...}})"
//	@Failure		404	{object}	fiber.Map	"User not found (e.g., {\"message\": \"User with id X not found\"})"
//	@Failure		500	{object}	fiber.Map	"Internal Server Error (e.g., {\"message\": \"Failed to retrieve user\"})"
//	@Router			/users/{id} [get]
func (h *UserHandler) GetUserById(c *fiber.Ctx) error {
	id := c.Params("id")
	user, err := h.UserService.GetUserById(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utils.HandleErrorResponse(c, fiber.StatusNotFound, "User with id "+id+" not found")
		}
		return utils.HandleErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve user") // Generic message
	}
	return utils.JsonResponse(c, fiber.StatusOK, fiber.Map{"user": user})
}

// DeleteUser godoc
//
//	@Summary		Delete a user
//	@Description	Delete a user by ID
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string		true	"User ID"
//	@Success		200	{object}	fiber.Map	"User deleted successfully (e.g., {\"message\": \"User deleted successfully\"})"
//	@Failure		404	{object}	fiber.Map	"User not found (e.g., {\"message\": \"User with id X not found\"})"
//	@Failure		500	{object}	fiber.Map	"Internal Server Error (e.g., {\"message\": \"Failed to delete user\"})"
//	@Router			/users/{id} [delete]
func (h *UserHandler) DeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")
	err := h.UserService.DeleteUser(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utils.HandleErrorResponse(c, fiber.StatusNotFound, "User with id "+id+" not found")
		}
		return utils.HandleErrorResponse(c, fiber.StatusInternalServerError, "Failed to delete user") // Generic message
	}
	return utils.JsonResponse(c, fiber.StatusOK, fiber.Map{"message": "User deleted successfully"})
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
//	@Success		200		{object}	fiber.Map					"User updated successfully (e.g., {\"message\": \"User updated successfully\"})"
//	@Failure		400		{object}	fiber.Map					"Invalid request body (e.g., {\"message\": \"Invalid request body\"}) or validation errors (e.g., {\"errors\": {\"field\": \"error message\"}})"
//	@Failure		404		{object}	fiber.Map					"User not found (e.g., {\"message\": \"User with id X not found for update\"})"
//	@Failure		500		{object}	fiber.Map					"Internal Server Error (e.g., {\"message\": \"Failed to update user\"})"
//	@Router			/users/{id} [put]
func (h *UserHandler) UpdateUser(ctx *fiber.Ctx) error {
	var req *models.UpdateUserRequest
	id := ctx.Params("id")
	if err := ctx.BodyParser(&req); err != nil {
		return utils.HandleErrorResponse(ctx, fiber.StatusBadRequest, "Invalid request body")
	}

	validationErrorsMap := utils.ValidateRequest(req)
	if validationErrorsMap != nil {
		return utils.JsonResponse(ctx, fiber.StatusBadRequest, fiber.Map{"errors": validationErrorsMap})
	}

	err := h.UserService.UpdateUser(id, req) // Assuming UpdateUser returns the updated user or specific errors
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utils.HandleErrorResponse(ctx, fiber.StatusNotFound, "User with id "+id+" not found for update")
		}
		// Consider if UpdateUser can return other specific errors, e.g., validation from service layer
		return utils.HandleErrorResponse(ctx, fiber.StatusInternalServerError, "Failed to update user") // Generic message
	}
	// Assuming UpdateUser is successful, and we want to return the updated user model or a success message
	// For now, sticking to the original plan of returning a message.
	// If you want to return the user:
	// updatedUser, err := h.UserService.GetUserById(id) // You might need to fetch it again if UpdateUser doesn't return it
	// if err != nil { ... handle this ... }
	// return utils.JsonResponse(ctx, fiber.StatusOK, fiber.Map{"user": updatedUser})
	return utils.JsonResponse(ctx, fiber.StatusOK, fiber.Map{"message": "User updated successfully"})
}

// CreateUser godoc
//
//	@Summary		Create a new user
//	@Description	Create a new user and return an authentication token
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			user	body		models.CreateUserRequest	true	"User details"
//	@Success		201		{object}	fiber.Map					"User created successfully, token returned (e.g., {\"token\": \"jwt_token_here\"})"
//	@Failure		400		{object}	fiber.Map					"Invalid request body (e.g., {\"message\": \"Invalid request body\"}) or validation errors (e.g., {\"errors\": {\"field\": \"error message\"}})"
//	@Failure		500		{object}	fiber.Map					"Internal Server Error (e.g., {\"message\": \"Failed to create user\"} or {\"message\": \"Failed to create token\"})"
//	@Router			/users [post]
func (h *UserHandler) CreateUser(ctx *fiber.Ctx) error {
	var req *models.CreateUserRequest
	if err := ctx.BodyParser(&req); err != nil {
		return utils.HandleErrorResponse(ctx, fiber.StatusBadRequest, "Invalid request body")
	}
	validationErrorsMap := utils.ValidateRequest(req)
	if validationErrorsMap != nil {
		return utils.JsonResponse(ctx, fiber.StatusBadRequest, fiber.Map{"errors": validationErrorsMap})
	}

	user := req.ToUser()
	userId, err := h.UserService.CreateUser(user) // Renamed id to userId for clarity
	if err != nil {
		// Consider specific errors from CreateUser, e.g., email already exists
		return utils.HandleErrorResponse(ctx, fiber.StatusInternalServerError, "Failed to create user") // Generic message
	}
	token, err := utils.CreateToken(userId)
	if err != nil {
		// This error wasn't handled before for CreateToken
		return utils.HandleErrorResponse(ctx, fiber.StatusInternalServerError, "Failed to create token")
	}
	return utils.JsonResponse(ctx, fiber.StatusCreated, fiber.Map{"token": token})
}

// Login godoc
//
//	@Summary		User login
//	@Description	Authenticate a user and return a token
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			credentials	body		models.LoginUserRequest	true	"Login credentials"
//	@Success		200			{object}	fiber.Map				"Authentication successful, token returned (e.g., {\"token\": \"jwt_token_here\"})"
//	@Failure		400			{object}	fiber.Map				"Invalid request body (e.g., {\"message\": \"Invalid request body\"}) or validation errors (e.g., {\"errors\": {\"field\": \"error message\"}})"
//	@Failure		401			{object}	fiber.Map				"Unauthorized (e.g., {\"message\": \"Invalid username or password\"})"
//	@Failure		500			{object}	fiber.Map				"Internal Server Error (e.g., {\"message\": \"Login failed\"} or {\"message\": \"Failed to create token\"})"
//	@Router			/login [post]
func (h *UserHandler) Login(ctx *fiber.Ctx) error {
	var req *models.LoginUserRequest
	if err := ctx.BodyParser(&req); err != nil {
		return utils.HandleErrorResponse(ctx, fiber.StatusBadRequest, "Invalid request body")
	}
	validationErrorsMap := utils.ValidateRequest(req)
	if validationErrorsMap != nil {
		return utils.JsonResponse(ctx, fiber.StatusBadRequest, fiber.Map{"errors": validationErrorsMap})
	}

	userId, err := h.UserService.LoginUser(req) // Renamed id to userId
	if err != nil {
		// LoginUser might return specific errors like "invalid credentials"
		// For now, using a generic message or err.Error() if it's safe
		if errors.Is(err, gorm.ErrRecordNotFound) || err.Error() == "invalid credentials" { // Example check
			return utils.HandleErrorResponse(ctx, fiber.StatusUnauthorized, "Invalid username or password")
		}
		return utils.HandleErrorResponse(ctx, fiber.StatusInternalServerError, "Login failed") // Generic message
	}
	token, err := utils.CreateToken(userId)
	if err != nil {
		// This error wasn't handled before for CreateToken
		return utils.HandleErrorResponse(ctx, fiber.StatusInternalServerError, "Failed to create token")
	}
	return utils.JsonResponse(ctx, fiber.StatusOK, fiber.Map{"token": token})
}

// GetMe godoc
//
//	@Summary		Get current user
//	@Description	Get details of the currently authenticated user
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Success		200	{object}	fiber.Map	"Details of the authenticated user (e.g., {\"user\": {\"id\": 1, ...}})"
//	@Router			/me [get]
func (h *UserHandler) GetMe(c *fiber.Ctx) error {
	user := c.Locals("user").(models.User)
	return utils.JsonResponse(c, fiber.StatusOK, fiber.Map{"user": user})
}
