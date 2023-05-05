package authentication

import (
	"github.com/dozheiny/it-captal-task/models/user"
	"github.com/dozheiny/it-captal-task/pkg/jwt"
	"github.com/dozheiny/it-captal-task/serializers"
	"github.com/gofiber/fiber/v2"
)

// SendAccessToken returns access token to client.
func SendAccessToken(c *fiber.Ctx) error {
	u, ok := c.Locals("user").(*user.Model)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(serializers.NewState().
			SetStatus(false).SetCode(fiber.StatusInternalServerError).SetMessage(errInternalServerError))
	}

	// Generate new JWT token.
	token := jwt.NewToken()
	if err := token.Initialize(u.ID, true); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(serializers.NewState().
			SetStatus(false).
			SetCode(fiber.StatusInternalServerError).
			SetMessage(errInternalServerError).SetDetails(err.Error()))
	}

	return c.Status(fiber.StatusOK).JSON(serializers.NewState().
		SetStatus(true).SetCode(fiber.StatusOK).SetMessage(successful).SetData(serializers.LoginResponse{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
	}))
}
