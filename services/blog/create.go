package blog

import (
	"github.com/dozheiny/it-captal-task/models/blog"
	"github.com/dozheiny/it-captal-task/models/user"
	"github.com/dozheiny/it-captal-task/pkg/mongo"
	"github.com/dozheiny/it-captal-task/serializers"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Create gets title and content form HTTP request and create new blog post.
func Create(c *fiber.Ctx) error {
	// Parse inputForm.
	inputForm := new(serializers.CreateBlog)

	if err := c.BodyParser(inputForm); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(serializers.NewState().
			SetCode(fiber.StatusBadRequest).
			SetStatus(false).SetMessage(errInvalidInputForm).SetData(err))
	}

	// Validate InputForm.
	if err := inputForm.Validation(); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(serializers.NewState().
			SetCode(fiber.StatusBadRequest).
			SetStatus(false).SetMessage(err.Error()))
	}

	b := blog.Model{
		ID:      primitive.NewObjectID(),
		Title:   inputForm.Title,
		Content: inputForm.Content,
	}

	// get user from context.
	u, ok := c.Locals("user").(*user.Model)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(serializers.NewState().
			SetStatus(false).
			SetCode(fiber.StatusInternalServerError).
			SetMessage(errInternalServerError))
	}

	b.User = *u

	// save document.
	if err := mongo.Insert(c.Context(), b); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(serializers.NewState().
			SetStatus(false).
			SetCode(fiber.StatusInternalServerError).
			SetMessage(errInternalServerError).SetDetails(err.Error()))
	}

	return c.Status(fiber.StatusCreated).JSON(serializers.NewState().
		SetCode(fiber.StatusCreated).SetStatus(true).SetMessage(successful))
}
