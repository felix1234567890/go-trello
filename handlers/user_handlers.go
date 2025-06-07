// Package handlers provides HTTP request handlers.
package handlers

import (
	"errors"
	"felix1234567890/go-trello/models"
	"strconv"
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

	err := h.UserService.UpdateUser(id, req)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utils.HandleErrorResponse(ctx, fiber.StatusNotFound, "User with id "+id+" not found for update")
		}
		return utils.HandleErrorResponse(ctx, fiber.StatusInternalServerError, "Failed to update user")
	}
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
		// Use custom error types for better error handling
		if errors.Is(err, utils.ErrInvalidCredentials) {
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

// FollowUser godoc
//	@Summary		Follow a user
//	@Description	Follow another user by their ID
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			id	path		string		true	"Target User ID"
//	@Success		200	{object}	fiber.Map	"Successfully followed user (e.g., {\"message\": \"Successfully followed user\"})"
//	@Failure		400	{object}	fiber.Map	"Invalid target user ID or cannot follow yourself (e.g., {\"message\": \"Invalid target user ID\"} or {\"message\": \"Cannot follow yourself\"})"
//	@Failure		401	{object}	fiber.Map	"Unauthorized (e.g., {\"message\": \"Unauthorized\"})"
//	@Failure		404	{object}	fiber.Map	"User or target user not found (e.g., {\"message\": \"User or target user not found\"})"
//	@Failure		500	{object}	fiber.Map	"Internal Server Error (e.g., {\"message\": \"Failed to follow user\"})"
//	@Router			/users/{id}/follow [post]
func (h *UserHandler) FollowUser(c *fiber.Ctx) error {
	targetUserIDStr := c.Params("id")
	targetUserID, err := strconv.ParseUint(targetUserIDStr, 10, 32)
	if err != nil {
		return utils.HandleErrorResponse(c, fiber.StatusBadRequest, "Invalid target user ID")
	}

	user := c.Locals("user").(models.User)
	userID := user.ID

	if userID == uint(targetUserID) {
		return utils.HandleErrorResponse(c, fiber.StatusBadRequest, "Cannot follow yourself")
	}

	err = h.UserService.FollowUser(userID, uint(targetUserID))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utils.HandleErrorResponse(c, fiber.StatusNotFound, "User or target user not found")
		}
		return utils.HandleErrorResponse(c, fiber.StatusInternalServerError, "Failed to follow user")
	}

	return utils.JsonResponse(c, fiber.StatusOK, fiber.Map{"message": "Successfully followed user"})
}

// UnfollowUser godoc
//	@Summary		Unfollow a user
//	@Description	Unfollow another user by their ID
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			id	path		string		true	"Target User ID"
//	@Success		200	{object}	fiber.Map	"Successfully unfollowed user (e.g., {\"message\": \"Successfully unfollowed user\"})"
//	@Failure		400	{object}	fiber.Map	"Invalid target user ID or cannot unfollow yourself (e.g., {\"message\": \"Invalid target user ID\"} or {\"message\": \"Cannot unfollow yourself\"})"
//	@Failure		401	{object}	fiber.Map	"Unauthorized (e.g., {\"message\": \"Unauthorized\"})"
//	@Failure		404	{object}	fiber.Map	"User or target user not found, or not following (e.g., {\"message\": \"User or target user not found, or not following\"})"
//	@Failure		500	{object}	fiber.Map	"Internal Server Error (e.g., {\"message\": \"Failed to unfollow user\"})"
//	@Router			/users/{id}/unfollow [post]
func (h *UserHandler) UnfollowUser(c *fiber.Ctx) error {
	targetUserIDStr := c.Params("id")
	targetUserID, err := strconv.ParseUint(targetUserIDStr, 10, 32)
	if err != nil {
		return utils.HandleErrorResponse(c, fiber.StatusBadRequest, "Invalid target user ID")
	}

	user := c.Locals("user").(models.User)
	userID := user.ID

	if userID == uint(targetUserID) {
		return utils.HandleErrorResponse(c, fiber.StatusBadRequest, "Cannot unfollow yourself")
	}

	err = h.UserService.UnfollowUser(userID, uint(targetUserID))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utils.HandleErrorResponse(c, fiber.StatusNotFound, "User or target user not found, or not following")
		}
		return utils.HandleErrorResponse(c, fiber.StatusInternalServerError, "Failed to unfollow user")
	}

	return utils.JsonResponse(c, fiber.StatusOK, fiber.Map{"message": "Successfully unfollowed user"})
}
