package routes

import (
	"github.com/dozheiny/it-captal-task/services/authentication"
	"github.com/gofiber/fiber/v2"
)

func RegisterAll(route fiber.Router) {

	route.Post("/login", authentication.FindUserByUserName, authentication.SendAccessToken)
}
