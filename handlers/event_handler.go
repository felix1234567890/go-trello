package handlers

import (
	"errors"
	"felix1234567890/go-trello/models"
	"felix1234567890/go-trello/service"
	"felix1234567890/go-trello/utils"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// EventHandler handles HTTP requests for events
type EventHandler struct {
	eventService service.EventService
}

// NewEventHandler creates a new EventHandler instance
func NewEventHandler(eventService service.EventService) *EventHandler {
	return &EventHandler{eventService: eventService}
}

// CreateEvent godoc
// @Summary      Create a new event
// @Description  Create a new event with the given details. Date can be YYYY-MM-DD or RFC3339.
// @Tags         events
// @Accept       json
// @Produce      json
// @Param        event body models.CreateEventRequest true "Event details"
// @Success      201  {object} models.Event
// @Failure      400  {object} fiber.Map "Invalid request body, validation error, or invalid date format"
// @Failure      500  {object} fiber.Map "Internal Server Error"
// @Router       /events [post]
func (h *EventHandler) CreateEvent(c *fiber.Ctx) error {
	var req models.CreateEventRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.HandleErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	if err := utils.ValidateRequest(req); err != nil {
		return utils.HandleErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	event, err := h.eventService.CreateEvent(&req)
	if err != nil {
		if err.Error() == "invalid date format" {
			return utils.HandleErrorResponse(c, fiber.StatusBadRequest, "Invalid date format. Use YYYY-MM-DD or RFC3339.")
		}
		return utils.HandleErrorResponse(c, fiber.StatusInternalServerError, "Failed to create event")
	}
	return utils.JsonResponse(c, fiber.StatusCreated, event)
}

// GetEventByID godoc
// @Summary      Get an event by ID
// @Description  Get details of a specific event by its ID.
// @Tags         events
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Event ID"
// @Success      200  {object}  models.Event
// @Failure      400  {object}  fiber.Map "Invalid event ID format"
// @Failure      404  {object}  fiber.Map "Event not found"
// @Failure      500  {object}  fiber.Map "Internal Server Error"
// @Router       /events/{id} [get]
func (h *EventHandler) GetEventByID(c *fiber.Ctx) error {
	idParam := c.Params("id")
	eventID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return utils.HandleErrorResponse(c, fiber.StatusBadRequest, "Invalid event ID format")
	}

	event, err := h.eventService.GetEventByID(uint(eventID))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utils.HandleErrorResponse(c, fiber.StatusNotFound, "Event not found")
		}
		return utils.HandleErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve event")
	}
	return utils.JsonResponse(c, fiber.StatusOK, event)
}

// GetAllEvents godoc
// @Summary      Get all events
// @Description  Get a list of all events.
// @Tags         events
// @Accept       json
// @Produce      json
// @Success      200  {array}  models.Event
// @Failure      500  {object} fiber.Map "Internal Server Error"
// @Router       /events [get]
func (h *EventHandler) GetAllEvents(c *fiber.Ctx) error {
	events, err := h.eventService.GetAllEvents()
	if err != nil {
		return utils.HandleErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve events")
	}
	return utils.JsonResponse(c, fiber.StatusOK, events)
}

// UpdateEvent godoc
// @Summary      Update an existing event
// @Description  Update an existing event with the given details. Date can be YYYY-MM-DD or RFC3339.
// @Tags         events
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Event ID"
// @Param        event body models.UpdateEventRequest true "Event details to update"
// @Success      200  {object} models.Event
// @Failure      400  {object} fiber.Map "Invalid request body, validation error, or invalid date format"
// @Failure      404  {object} fiber.Map "Event not found"
// @Failure      500  {object} fiber.Map "Internal Server Error"
// @Router       /events/{id} [put]
func (h *EventHandler) UpdateEvent(c *fiber.Ctx) error {
	idParam := c.Params("id")
	eventID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return utils.HandleErrorResponse(c, fiber.StatusBadRequest, "Invalid event ID format")
	}

	var req models.UpdateEventRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.HandleErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	if err := utils.ValidateRequest(req); err != nil {
		return utils.HandleErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	event, err := h.eventService.UpdateEvent(uint(eventID), &req)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utils.HandleErrorResponse(c, fiber.StatusNotFound, "Event not found")
		}
		if err.Error() == "invalid date format for update" {
			return utils.HandleErrorResponse(c, fiber.StatusBadRequest, "Invalid date format. Use YYYY-MM-DD or RFC3339.")
		}
		return utils.HandleErrorResponse(c, fiber.StatusInternalServerError, "Failed to update event")
	}
	return utils.JsonResponse(c, fiber.StatusOK, event)
}

// DeleteEvent godoc
// @Summary      Delete an event by ID
// @Description  Delete a specific event by its ID.
// @Tags         events
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Event ID"
// @Success      204  {object}  nil "Event deleted successfully"
// @Failure      400  {object}  fiber.Map "Invalid event ID format"
// @Failure      404  {object}  fiber.Map "Event not found"
// @Failure      500  {object}  fiber.Map "Internal Server Error"
// @Router       /events/{id} [delete]
func (h *EventHandler) DeleteEvent(c *fiber.Ctx) error {
	idParam := c.Params("id")
	eventID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return utils.HandleErrorResponse(c, fiber.StatusBadRequest, "Invalid event ID format")
	}

	err = h.eventService.DeleteEvent(uint(eventID))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utils.HandleErrorResponse(c, fiber.StatusNotFound, "Event not found")
		}
		return utils.HandleErrorResponse(c, fiber.StatusInternalServerError, "Failed to delete event")
	}
	return c.SendStatus(fiber.StatusNoContent)
}

// AddUserToEvent godoc
// @Summary      Add a user to an event
// @Description  Associate a user with an event.
// @Tags         events
// @Accept       json
// @Produce      json
// @Param        event_id path int true "Event ID"
// @Param        user_id  path int true "User ID"
// @Success      200 {object} fiber.Map "User added to event successfully"
// @Failure      400 {object} fiber.Map "Invalid event ID or user ID format"
// @Failure      404 {object} fiber.Map "Event or User not found"
// @Failure      500 {object} fiber.Map "Internal Server Error"
// @Router       /events/{event_id}/users/{user_id} [post]
func (h *EventHandler) AddUserToEvent(c *fiber.Ctx) error {
	eventIDParam := c.Params("event_id")
	eventID, err := strconv.ParseUint(eventIDParam, 10, 32)
	if err != nil {
		return utils.HandleErrorResponse(c, fiber.StatusBadRequest, "Invalid event ID format")
	}

	userIDParam := c.Params("user_id")
	userID, err := strconv.ParseUint(userIDParam, 10, 32)
	if err != nil {
		return utils.HandleErrorResponse(c, fiber.StatusBadRequest, "Invalid user ID format")
	}

	err = h.eventService.AddUserToEvent(uint(eventID), uint(userID))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utils.HandleErrorResponse(c, fiber.StatusNotFound, "Event or User not found")
		}
		return utils.HandleErrorResponse(c, fiber.StatusInternalServerError, "Failed to add user to event")
	}
	return utils.JsonResponse(c, fiber.StatusOK, fiber.Map{"message": "User added to event successfully"})
}

// RemoveUserFromEvent godoc
// @Summary      Remove a user from an event
// @Description  Disassociate a user from an event.
// @Tags         events
// @Accept       json
// @Produce      json
// @Param        event_id path int true "Event ID"
// @Param        user_id  path int true "User ID"
// @Success      200 {object} fiber.Map "User removed from event successfully"
// @Failure      400 {object} fiber.Map "Invalid event ID or user ID format"
// @Failure      404 {object} fiber.Map "Event or User not found, or user not in event"
// @Failure      500 {object} fiber.Map "Internal Server Error"
// @Router       /events/{event_id}/users/{user_id} [delete]
func (h *EventHandler) RemoveUserFromEvent(c *fiber.Ctx) error {
	eventIDParam := c.Params("event_id")
	eventID, err := strconv.ParseUint(eventIDParam, 10, 32)
	if err != nil {
		return utils.HandleErrorResponse(c, fiber.StatusBadRequest, "Invalid event ID format")
	}

	userIDParam := c.Params("user_id")
	userID, err := strconv.ParseUint(userIDParam, 10, 32)
	if err != nil {
		return utils.HandleErrorResponse(c, fiber.StatusBadRequest, "Invalid user ID format")
	}

	err = h.eventService.RemoveUserFromEvent(uint(eventID), uint(userID))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utils.HandleErrorResponse(c, fiber.StatusNotFound, "Event or User not found, or user not in event")
		}
		return utils.HandleErrorResponse(c, fiber.StatusInternalServerError, "Failed to remove user from event")
	}
	return utils.JsonResponse(c, fiber.StatusOK, fiber.Map{"message": "User removed from event successfully"})
}

// AddGroupToEvent godoc
// @Summary      Add a group to an event
// @Description  Associate a group with an event.
// @Tags         events
// @Accept       json
// @Produce      json
// @Param        event_id path int true "Event ID"
// @Param        group_id path int true "Group ID"
// @Success      200 {object} fiber.Map "Group added to event successfully"
// @Failure      400 {object} fiber.Map "Invalid event ID or group ID format"
// @Failure      404 {object} fiber.Map "Event or Group not found"
// @Failure      500 {object} fiber.Map "Internal Server Error"
// @Router       /events/{event_id}/groups/{group_id} [post]
func (h *EventHandler) AddGroupToEvent(c *fiber.Ctx) error {
	eventIDParam := c.Params("event_id")
	eventID, err := strconv.ParseUint(eventIDParam, 10, 32)
	if err != nil {
		return utils.HandleErrorResponse(c, fiber.StatusBadRequest, "Invalid event ID format")
	}

	groupIDParam := c.Params("group_id")
	groupID, err := strconv.ParseUint(groupIDParam, 10, 32)
	if err != nil {
		return utils.HandleErrorResponse(c, fiber.StatusBadRequest, "Invalid group ID format")
	}

	err = h.eventService.AddGroupToEvent(uint(eventID), uint(groupID))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utils.HandleErrorResponse(c, fiber.StatusNotFound, "Event or Group not found")
		}
		return utils.HandleErrorResponse(c, fiber.StatusInternalServerError, "Failed to add group to event")
	}
	return utils.JsonResponse(c, fiber.StatusOK, fiber.Map{"message": "Group added to event successfully"})
}

// RemoveGroupFromEvent godoc
// @Summary      Remove a group from an event
// @Description  Disassociate a group from an event.
// @Tags         events
// @Accept       json
// @Produce      json
// @Param        event_id path int true "Event ID"
// @Param        group_id path int true "Group ID"
// @Success      200 {object} fiber.Map "Group removed from event successfully"
// @Failure      400 {object} fiber.Map "Invalid event ID or group ID format"
// @Failure      404 {object} fiber.Map "Event or Group not found, or group not in event"
// @Failure      500 {object} fiber.Map "Internal Server Error"
// @Router       /events/{event_id}/groups/{group_id} [delete]
func (h *EventHandler) RemoveGroupFromEvent(c *fiber.Ctx) error {
	eventIDParam := c.Params("event_id")
	eventID, err := strconv.ParseUint(eventIDParam, 10, 32)
	if err != nil {
		return utils.HandleErrorResponse(c, fiber.StatusBadRequest, "Invalid event ID format")
	}

	groupIDParam := c.Params("group_id")
	groupID, err := strconv.ParseUint(groupIDParam, 10, 32)
	if err != nil {
		return utils.HandleErrorResponse(c, fiber.StatusBadRequest, "Invalid group ID format")
	}

	err = h.eventService.RemoveGroupFromEvent(uint(eventID), uint(groupID))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utils.HandleErrorResponse(c, fiber.StatusNotFound, "Event or Group not found, or group not in event")
		}
		return utils.HandleErrorResponse(c, fiber.StatusInternalServerError, "Failed to remove group from event")
	}
	return utils.JsonResponse(c, fiber.StatusOK, fiber.Map{"message": "Group removed from event successfully"})
}
