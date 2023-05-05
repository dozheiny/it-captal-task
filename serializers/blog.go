package serializers

import "errors"

var (
	errTitleCannotEmpty   = errors.New("title cannot empty")
	errContentCannotEmpty = errors.New("content cannot empty")
)

type CreateBlog struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

func (i *CreateBlog) Validation() error {
	if len(i.Title) == 0 {
		return errTitleCannotEmpty
	}

	if len(i.Content) == 0 {
		return errContentCannotEmpty
	}

	return nil
}

type GetBlog struct {
	Title    string `query:"title"`
	Content  string `query:"content"`
	Username string `query:"username"`
	Page     int    `query:"page"`
	PerPage  int    `query:"per_page"`
}
