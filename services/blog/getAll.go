package blog

import (
	"errors"
	"github.com/dozheiny/it-captal-task/models/blog"
	"github.com/dozheiny/it-captal-task/pkg/mongo"
	"github.com/dozheiny/it-captal-task/serializers"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

// GetAll gets query parameters from a URL path and create a query.
// then returns blogs contain that query
// @Description get blogs
// @Summary get blogs
// @Tags blog
// @Router /blog [GET]
// @Accept application/json
// @Produce json
// @Param inputForm query serializers.GetBlog true "input queries"
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>) // @Success 200 {object} serializers.State
// @Failure 400 {object} serializers.State
// @Failure 500 {object} serializers.State
// @Failure 401 {object} serializers.State
func GetAll(c *fiber.Ctx) error {
	// Parse inputForm.
	inputForm := new(serializers.GetBlog)

	if err := c.QueryParser(inputForm); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(serializers.NewState().
			SetCode(fiber.StatusBadRequest).
			SetStatus(false).SetMessage(errInvalidInputForm).SetData(err))
	}

	// Handle Page and PerPage.
	if inputForm.Page <= 0 {
		inputForm.Page = 1
	}

	if inputForm.PerPage <= 0 {
		inputForm.PerPage = 5
	}

	// create query.
	query := bson.D{}

	if len(inputForm.Title) != 0 {
		query = append(query, bson.E{Key: "title", Value: bson.D{{Key: "$regex", Value: inputForm.Title}}})
	}

	if len(inputForm.Content) != 0 {
		query = append(query, bson.E{Key: "content", Value: bson.D{{Key: "$regex", Value: inputForm.Content}}})
	}

	if len(inputForm.Username) != 0 {
		query = append(query, bson.E{Key: "user.username", Value: bson.D{{Key: "$regex", Value: inputForm.Username}}})
	}

	cur, counters, err := mongo.FindAll(c, blog.Model{}, query, &inputForm.Page, &inputForm.PerPage)
	if err != nil {
		if errors.Is(err, mongo.NoDocuments) {
			return c.Status(fiber.StatusNotFound).JSON(serializers.
				NewState().SetCode(fiber.StatusNotFound).SetStatus(false).SetMessage(errBlogNotFound))
		}

		return c.Status(fiber.StatusInternalServerError).JSON(serializers.
			NewState().SetCode(fiber.StatusInternalServerError).
			SetMessage(errInternalServerError).SetDetails(err.Error()))
	}

	blogs := make([]blog.Model, 1)
	for cur.Next(c.Context()) {
		b := new(blog.Model)

		if err := cur.Decode(b); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(serializers.
				NewState().SetCode(fiber.StatusInternalServerError).
				SetMessage(errInternalServerError).SetDetails(err.Error()))
		}

		blogs = append(blogs, *b)
	}

	return c.Status(fiber.StatusOK).JSON(serializers.NewState().
		SetCode(fiber.StatusOK).SetStatus(true).SetMessage(successful).SetData(blogs).
		SetCounters(serializers.Counts{
			Page:      inputForm.Page,
			PerPage:   inputForm.PerPage,
			TotalPage: counters.TotalPage,
			Total:     counters.TotalDocument,
		}))
}
