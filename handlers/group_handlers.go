package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"felix1234567890/go-trello/models"
	"felix1234567890/go-trello/repository"
)

type GroupHandler struct {
	Repo *repository.GroupRepository
}

func NewGroupHandler(repo *repository.GroupRepository) *GroupHandler {
	return &GroupHandler{Repo: repo}
}

func (h *GroupHandler) CreateGroup(c *fiber.Ctx) error {
	var group models.Group
	if err := c.BodyParser(&group); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	if err := h.Repo.CreateGroup(&group); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(group)
}

func (h *GroupHandler) GetGroup(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid group id"})
	}
	group, err := h.Repo.GetGroupByID(uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "group not found"})
	}
	return c.Status(fiber.StatusOK).JSON(group)
}

func (h *GroupHandler) GetAllGroups(c *fiber.Ctx) error {
	groups, err := h.Repo.GetAllGroups()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(groups)
}

func (h *GroupHandler) UpdateGroup(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid group id"})
	}
	var group models.Group
	if err := c.BodyParser(&group); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	group.ID = uint(id)
	if err := h.Repo.UpdateGroup(&group); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(group)
}

func (h *GroupHandler) DeleteGroup(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid group id"})
	}
	if err := h.Repo.DeleteGroup(uint(id)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.SendStatus(fiber.StatusNoContent)
}

func (h *GroupHandler) AddUserToGroup(c *fiber.Ctx) error {
	groupID, err := strconv.Atoi(c.Params("group_id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid group id"})
	}
	userID, err := strconv.Atoi(c.Params("user_id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid user id"})
	}
	if err := h.Repo.AddUserToGroup(uint(groupID), uint(userID)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "user added to group"})
}

func (h *GroupHandler) RemoveUserFromGroup(c *fiber.Ctx) error {
	groupID, err := strconv.Atoi(c.Params("group_id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid group id"})
	}
	userID, err := strconv.Atoi(c.Params("user_id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid user id"})
	}
	if err := h.Repo.RemoveUserFromGroup(uint(groupID), uint(userID)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "user removed from group"})
}
