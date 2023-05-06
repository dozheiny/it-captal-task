package authentication

import (
	"errors"
	"github.com/dozheiny/it-captal-task/models/user"
	"github.com/dozheiny/it-captal-task/pkg/mongo"
	"github.com/dozheiny/it-captal-task/serializers"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

// FindUserByUserName parsers serializer from serializers.Login;
// After Validate save user into context.
// @Description auth
// @Summary login
// @Tags authentication
// @Router /login [POST]
// @Accept application/json
// @Produce json
// @Param inputForm body serializers.Login true "input forms"
// @Success 200 {object} serializers.State
// @Failure 400 {object} serializers.State
// @Failure 500 {object} serializers.State
// @Failure 401 {object} serializers.State
func FindUserByUserName(c *fiber.Ctx) error {
	// Parse inputForm.
	inputForm := new(serializers.Login)
	if err := c.BodyParser(inputForm); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(serializers.NewState().
			SetCode(fiber.StatusBadRequest).
			SetStatus(false).SetMessage(errInvalidInputForm).SetDetails(err.Error()))
	}

	// Validate InputForm.
	if err := inputForm.Validation(); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(serializers.NewState().
			SetCode(fiber.StatusBadRequest).
			SetStatus(false).SetMessage(err.Error()))
	}

	// Find user by username in a database.
	u := new(user.Model)

	if err := mongo.FindOne(c.Context(), u, bson.D{{Key: "username", Value: inputForm.Username}}); err != nil {
		if errors.Is(err, mongo.NoDocuments) {
			return c.Status(fiber.StatusNotFound).JSON(serializers.
				NewState().SetCode(fiber.StatusNotFound).SetStatus(false).SetMessage(errUserNotFound))
		}

		return c.Status(fiber.StatusInternalServerError).JSON(serializers.
			NewState().SetCode(fiber.StatusInternalServerError).SetMessage(errInternalServerError))
	}

	// check Password.
	if !u.VerifyPassword(inputForm.Password) {
		return c.Status(fiber.StatusForbidden).JSON(serializers.NewState().
			SetCode(fiber.StatusForbidden).SetStatus(false).SetMessage(errPasswordIsWrong))
	}

	c.Locals("user", u)
	return c.Next()
}
