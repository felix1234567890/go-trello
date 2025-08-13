package routes

import (
	"felix1234567890/go-trello/handlers"
	"felix1234567890/go-trello/repository"
	"felix1234567890/go-trello/service"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// SetupEventRoutes configures the routes for event-related operations.
func SetupEventRoutes(router fiber.Router, db *gorm.DB) {
	// Initialize repository, service, and handler
	eventRepository := repository.NewEventRepository(db)
	eventService := service.NewEventService(eventRepository)
	eventHandler := handlers.NewEventHandler(eventService)

	// Create a new route group for events
	eventGroup := router.Group("/events")

	// Define routes for event operations
	eventGroup.Post("/", eventHandler.CreateEvent)
	eventGroup.Get("/", eventHandler.GetAllEvents)
	eventGroup.Get("/:id", eventHandler.GetEventByID)
	eventGroup.Put("/:id", eventHandler.UpdateEvent)
	eventGroup.Delete("/:id", eventHandler.DeleteEvent)

	// Define routes for managing event users
	eventGroup.Post("/:event_id/users/:user_id", eventHandler.AddUserToEvent)
	eventGroup.Delete("/:event_id/users/:user_id", eventHandler.RemoveUserFromEvent)

	// Define routes for managing event groups
	eventGroup.Post("/:event_id/groups/:group_id", eventHandler.AddGroupToEvent)
	eventGroup.Delete("/:event_id/groups/:group_id", eventHandler.RemoveGroupFromEvent)
}
