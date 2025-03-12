package bookauthorsv1

import (
	"net/url"
)

type CreateBookAuthorsParams struct {
	ID int32 `uri:"id" binding:"required,min=1"`
}

func (p *CreateBookAuthorsParams) SelfQuery() url.Values {
	values := url.Values{}
	return values
}

type DeleteBookAuthorsParams struct {
	ID       int32 `uri:"id" binding:"required,min=1"`
	AuthorID int32 `uri:"author_id" binding:"required,min=1"`
}

func (p *DeleteBookAuthorsParams) SelfQuery() url.Values {
	values := url.Values{}
	return values
}

type ListBookAuthorsParams struct {
	ID int32 `uri:"id" binding:"required,min=1"`
}

func (p *ListBookAuthorsParams) SelfQuery() url.Values {
	values := url.Values{}
	return values
}
