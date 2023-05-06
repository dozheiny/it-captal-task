package user

import (
	"errors"
	"github.com/dozheiny/it-captal-task/models/user"
	"github.com/dozheiny/it-captal-task/pkg/mongo"
	"github.com/dozheiny/it-captal-task/serializers"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Create gets body HTTP request and serialize it as serializer.Login;
// Then check username, after checking username is not existed, creates that user
// @Description create user
// @Summary create user
// @Tags user
// @Router /user [POST]
// @Accept application/json
// @Produce json
// @Param inputForm body serializers.Login true "input forms"
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>) // @Success 200 {object} serializers.State
// @Failure 400 {object} serializers.State
// @Failure 500 {object} serializers.State
// @Failure 401 {object} serializers.State
func Create(c *fiber.Ctx) error {
	// Parse inputForm.
	inputForm := new(serializers.Login)

	if err := c.BodyParser(inputForm); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(serializers.NewState().
			SetCode(fiber.StatusBadRequest).
			SetStatus(false).SetMessage(errInvalidInputForm).SetDetails(err.Error()))
	}

	// Validate inputForm.
	if err := inputForm.Validation(); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(serializers.NewState().
			SetCode(fiber.StatusBadRequest).
			SetStatus(false).SetMessage(err.Error()))
	}

	// Find User by username
	// If user exists returns error as this username is existed.
	u := new(user.Model)

	if err := mongo.FindOne(c.Context(), u, bson.D{{Key: "username", Value: inputForm.Username}}); err == nil {
		return c.Status(fiber.StatusBadRequest).JSON(serializers.NewState().
			SetStatus(false).
			SetCode(fiber.StatusBadRequest).
			SetMessage(errThisUsernameExist))
	} else {
		if !errors.Is(err, mongo.NoDocuments) {
			return c.Status(fiber.StatusInternalServerError).JSON(serializers.NewState().
				SetStatus(false).
				SetCode(fiber.StatusInternalServerError).
				SetMessage(errInternalServerError).SetDetails(err.Error()))
		}
	}

	u.ID = primitive.NewObjectID()
	u.Username = inputForm.Username
	u.SetPassword(inputForm.Password)

	if err := mongo.Insert(c.Context(), u); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(serializers.NewState().
			SetStatus(false).
			SetCode(fiber.StatusInternalServerError).
			SetMessage(errInternalServerError).SetDetails(err.Error()))
	}

	return c.Status(fiber.StatusCreated).JSON(serializers.NewState().
		SetStatus(true).SetCode(fiber.StatusCreated).SetMessage(successful))
}
