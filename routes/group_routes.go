package routes

import (
	"felix1234567890/go-trello/handlers"
	"github.com/gofiber/fiber/v2"
)

func SetupGroupRoutes(router fiber.Router, handler *handlers.GroupHandler) {
	group := router.Group("/groups")
	group.Post("/", handler.CreateGroup)
	group.Get("/", handler.GetAllGroups)
	group.Get(":id", handler.GetGroup)
	group.Put(":id", handler.UpdateGroup)
	group.Delete(":id", handler.DeleteGroup)
	group.Post(":group_id/users/:user_id", handler.AddUserToGroup)
	group.Delete(":group_id/users/:user_id", handler.RemoveUserFromGroup)
}

