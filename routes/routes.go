package routes

import (
	"github.com/dozheiny/it-captal-task/pkg/jwt"
	"github.com/dozheiny/it-captal-task/services/authentication"
	"github.com/dozheiny/it-captal-task/services/blog"
	"github.com/gofiber/fiber/v2"
)

func RegisterAll(route fiber.Router) {

	// ---------------- AUTHENTICATION ROUTES ---------------- \\
	route.Post("/login", authentication.FindUserByUserName, authentication.SendAccessToken)
	route.Patch("/refresh", authentication.FindUserByRefreshToken, authentication.SendAccessToken)
	route.Patch("/logout", authentication.FindUserByRefreshToken, authentication.BlockToken)

	// ---------------- BLOG ROUTES ---------------- \\
	route.Post("/blog", jwt.AuthMiddleware, blog.Create)
}
