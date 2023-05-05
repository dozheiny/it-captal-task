package authentication

import (
	"github.com/dozheiny/it-captal-task/pkg/redis"
	"github.com/dozheiny/it-captal-task/serializers"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

// BlockToken save refresh token uuid as block token into redis.
func BlockToken(c *fiber.Ctx) error {
	// extract refresh uuid and refresh token expires time from jwt.MapClaims that gives from context.
	token, ok := c.Locals("token").(jwt.MapClaims)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(serializers.
			NewState().SetCode(fiber.StatusInternalServerError).SetMessage(errInternalServerError))
	}

	refreshUUID, ok := token["refresh_uuid"].(string)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(serializers.
			NewState().SetCode(fiber.StatusInternalServerError).SetMessage(errInternalServerError))
	}

	// check this refresh UUID exist in redis or not.
	refreshCheck, err := redis.GetValue(c.Context(), "block_"+refreshUUID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(serializers.
			NewState().SetCode(fiber.StatusInternalServerError).SetMessage(errInternalServerError))
	}

	if len(refreshCheck.Val()) != 0 {
		return c.Status(fiber.StatusUnauthorized).JSON(serializers.NewState().
			SetStatus(false).SetCode(fiber.StatusUnauthorized).SetMessage(errLoginAgain))
	}

	// save refresh uuid and access uuid value into redis as a block list.
	cmd, err := redis.SaveKey(c.Context(), "block_"+refreshUUID, true,
		time.Duration(token["exp"].(float64)))

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(serializers.
			NewState().SetCode(fiber.StatusInternalServerError).SetMessage(errInternalServerError).
			SetDetails(err.Error()))
	}

	if cmd.Err() != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(serializers.
			NewState().SetCode(fiber.StatusInternalServerError).SetMessage(errInternalServerError).
			SetDetails(cmd.Err().Error()))
	}

	return c.Status(fiber.StatusOK).JSON(serializers.NewState().
		SetCode(fiber.StatusOK).SetStatus(true).SetMessage(successful))
}
