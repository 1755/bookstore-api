package schemas

import (
	"fmt"
	"time"

	"github.com/1755/bookstore-api/internal/book"
)

type CreateBookRequest struct {
	Title         string `json:"title" validate:"required,min=1,max=1024"`
	Summary       string `json:"summary" validate:"max=5000"`
	PublishedYear int32  `json:"publishedYear"`
}

type UpdateBookRequest struct {
	Title         *string `json:"title,omitempty" validate:"omitempty,min=1,max=1024"`
	Summary       *string `json:"summary,omitempty" validate:"omitempty,max=5000"`
	PublishedYear *int32  `json:"publishedYear,omitempty"`
}

type Book struct {
	Title         string    `json:"title"`
	Summary       string    `json:"summary"`
	PublishedYear int32     `json:"publishedYear"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

func NewBookResourceFromModel(model *book.Model) Resource[Book] {
	return Resource[Book]{
		ID:   fmt.Sprintf("%d", model.ID),
		Type: "books",
		Attributes: Book{
			Title:         model.Title,
			Summary:       model.Summary,
			PublishedYear: model.PublishedYear,
			CreatedAt:     model.CreatedAt,
			UpdatedAt:     model.UpdatedAt,
		},
	}
}

func NewBookDocumentFromModel(model *book.Model, baseURL string) Document[Resource[Book]] {
	return Document[Resource[Book]]{
		Links: &DocumentLink{
			Self: fmt.Sprintf("%s/v1/authors/%d", baseURL, model.ID),
		},
		Data: NewBookResourceFromModel(model),
	}
}
