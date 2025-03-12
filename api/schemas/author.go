package schemas

import (
	"fmt"
	"time"

	"github.com/1755/bookstore-api/internal/author"
)

type CreateAuthorRequest struct {
	Name string `json:"name" binding:"required,min=1,max=150"`
	Bio  string `json:"bio" binding:"omitempty,max=5000"`
}

type UpdateAuthorRequest struct {
	Name *string `json:"name,omitempty" binding:"omitempty,min=1,max=150"`
	Bio  *string `json:"bio,omitempty" binding:"omitempty,max=5000"`
}

type Author struct {
	Name      string    `json:"name"`
	Bio       string    `json:"bio"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func NewAuthorResourceFromModel(model *author.Model) Resource[Author] {
	return Resource[Author]{
		ID:   fmt.Sprintf("%d", model.ID),
		Type: "authors",
		Attributes: Author{
			Name:      model.Name,
			Bio:       model.Bio,
			CreatedAt: model.CreatedAt,
			UpdatedAt: model.UpdatedAt,
		},
	}
}

func NewAuthorDocumentFromModel(model *author.Model, baseURL string) Document[Resource[Author]] {
	return Document[Resource[Author]]{
		Links: &DocumentLink{
			Self: fmt.Sprintf("%s/v1/authors/%d", baseURL, model.ID),
		},
		Data: NewAuthorResourceFromModel(model),
	}
}
