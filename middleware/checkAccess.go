package middleware

import (
	"github.com/dozheiny/it-captal-task/models/user"
	"github.com/dozheiny/it-captal-task/serializers"
	"github.com/gofiber/fiber/v2"
)

const (
	errInternalServerError = "internal server error"
)

const (
	accessDenied = "access denied"
)

func AuthMiddleware(c *fiber.Ctx) error {
	// get user from context.
	u, ok := c.Locals("user").(*user.Model)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(serializers.NewState().
			SetStatus(false).SetCode(fiber.StatusInternalServerError).SetMessage(errInternalServerError))
	}

	// Check if the requested endpoint is accessible or not
	endpoint := c.Path()
	method := c.Method()

	for _, e := range u.Endpoints {
		if e.Endpoint == endpoint && e.Method == method {
			return c.Next()
		}
	}

	return c.Status(fiber.StatusForbidden).JSON(serializers.NewState().
		SetStatus(false).
		SetCode(fiber.StatusForbidden).SetMessage(accessDenied))
}
