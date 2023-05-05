package authentication

import (
	"errors"
	"github.com/dozheiny/it-captal-task/models/user"
	"github.com/dozheiny/it-captal-task/pkg/jwt"
	"github.com/dozheiny/it-captal-task/pkg/mongo"
	"github.com/dozheiny/it-captal-task/serializers"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// FindUserByRefreshToken find user with refresh token;
// Then save user data into context.
func FindUserByRefreshToken(c *fiber.Ctx) error {
	// Parse inputForm.
	inputForm := new(serializers.RefreshToken)

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

	token, err := jwt.ExtractRefreshMetadata(inputForm.RefreshToken)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(serializers.NewState().
			SetCode(fiber.StatusUnauthorized).SetStatus(false).SetMessage(err.Error()))
	}

	userId, ok := token["userId"].(string)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(serializers.NewState().
			SetStatus(false).
			SetCode(fiber.StatusInternalServerError).
			SetMessage(errInternalServerError).SetDetails(err.Error()))
	}

	// convert userId from string to objectId.
	id, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(serializers.NewState().
			SetStatus(false).
			SetCode(fiber.StatusInternalServerError).
			SetMessage(errInternalServerError).SetDetails(err.Error()))
	}

	// Find user by username in a database.
	u := new(user.Model)

	if err := mongo.FindOne(c.Context(), u, bson.D{{Key: "_id", Value: id}}); err != nil {
		if errors.Is(err, mongo.NoDocuments) {
			return c.Status(fiber.StatusNotFound).JSON(serializers.
				NewState().SetCode(fiber.StatusNotFound).SetStatus(false).SetMessage(errUserNotFound))
		}

		return c.Status(fiber.StatusInternalServerError).JSON(serializers.
			NewState().SetCode(fiber.StatusInternalServerError).SetMessage(errInternalServerError))
	}

	c.Locals("user", u)
	c.Locals("token", token)
	return c.Next()
}
