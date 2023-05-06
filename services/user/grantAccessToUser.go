package user

import (
	"errors"
	"github.com/dozheiny/it-captal-task/models/endpoints"
	user2 "github.com/dozheiny/it-captal-task/models/user"
	"github.com/dozheiny/it-captal-task/pkg/mongo"
	"github.com/dozheiny/it-captal-task/serializers"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// GrantAccessToUser grants access to one user/
// @Description grant access to user
// @Summary grant access to user
// @Tags user
// @Router /user/:user-id/grant-access [PUT]
// @Accept application/json
// @Produce json
// @Param inputForm body serializers.GrantAccess true "input forms"
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>) // @Success 200 {object} serializers.State
// @Failure 400 {object} serializers.State
// @Failure 500 {object} serializers.State
// @Failure 401 {object} serializers.State
func GrantAccessToUser(c *fiber.Ctx) error {
	// get userID from an url path.
	urlPath := &struct {
		ID primitive.ObjectID `params:"user"`
	}{}

	// Parse params.
	if err := c.ParamsParser(urlPath); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(serializers.NewState().
			SetStatus(false).
			SetCode(fiber.StatusBadRequest).
			SetMessage(errInvalidUserURLPath))
	}

	// find User by ID.
	user := new(user2.Model)

	if err := mongo.FindOne(c.Context(), user, bson.D{{Key: "_id", Value: urlPath.ID}}); err != nil {
		if errors.Is(err, mongo.NoDocuments) {
			return c.Status(fiber.StatusNotFound).JSON(serializers.
				NewState().SetCode(fiber.StatusNotFound).SetStatus(false).SetMessage(errUserNotFound))
		}

		return c.Status(fiber.StatusInternalServerError).JSON(serializers.
			NewState().SetCode(fiber.StatusInternalServerError).SetMessage(errInternalServerError))
	}

	// Parse inputForm.
	inputForm := new(serializers.GrantAccess)

	if err := c.BodyParser(inputForm); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(serializers.NewState().
			SetStatus(false).
			SetCode(fiber.StatusBadRequest).
			SetMessage(errInvalidInputForm))
	}

	for i := range inputForm.IDs {
		endpoint := new(endpoints.Model)
		if err := mongo.FindOne(c.Context(), endpoint, bson.D{{Key: "_id", Value: inputForm.IDs[i]}}); err != nil {
			if errors.Is(err, mongo.NoDocuments) {
				return c.Status(fiber.StatusNotFound).JSON(serializers.
					NewState().SetCode(fiber.StatusNotFound).SetStatus(false).SetMessage(errEndpointNotFound))
			}

			return c.Status(fiber.StatusInternalServerError).JSON(serializers.
				NewState().SetCode(fiber.StatusInternalServerError).SetMessage(errInternalServerError))
		}

		user.Endpoints = append(user.Endpoints, *endpoint)
	}

	// grant access.
	if err := mongo.UpdateByID(c, user, &user.ID, bson.D{{Key: "$set", Value: user}}); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(serializers.
			NewState().SetCode(fiber.StatusInternalServerError).SetMessage(errInternalServerError))
	}

	return c.Status(fiber.StatusAccepted).JSON(serializers.NewState().
		SetStatus(true).SetCode(fiber.StatusAccepted).SetMessage(successful))
}
