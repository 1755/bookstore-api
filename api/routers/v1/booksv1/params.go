package booksv1

import (
	"fmt"
	"net/url"

	"github.com/1755/bookstore-api/internal/book"
)

type GetBooksParams struct {
	PageLimit   uint   `form:"page[limit]" binding:"required,min=1,max=100"`
	PageOffset  uint   `form:"page[offset]" binding:"omitempty"`
	Sort        string `form:"sort" binding:"omitempty,oneof=title -title publishedYear -publishedYear"`
	FilterTitle string `form:"filter[title]" binding:"omitempty,min=2,max=50"`
}

var mapSortField = map[string]string{
	"name":      "name",
	"updatedAt": "updated_at",
}

func (p *GetBooksParams) ToGetManyParams() *book.GetManyParams {
	var sort *book.GetManyParamsSort = nil
	if p.Sort != "" {
		sort = &book.GetManyParamsSort{}
		direction := "asc"
		value := p.Sort
		if p.Sort[0] == '-' {
			direction = "desc"
			value = p.Sort[1:]
		}

		field, ok := mapSortField[value]
		if ok {
			sort = &book.GetManyParamsSort{
				Field:     field,
				Direction: direction,
			}
		}
	}

	return &book.GetManyParams{
		Offset:      p.PageOffset,
		Limit:       p.PageLimit,
		Sort:        sort,
		FilterTitle: p.FilterTitle,
	}
}

func (p *GetBooksParams) SelfQuery() url.Values {
	values := url.Values{}
	values.Add("page[limit]", fmt.Sprintf("%d", p.PageLimit))
	if p.PageOffset > 0 {
		values.Add("page[offset]", fmt.Sprintf("%d", p.PageOffset))
	}
	if p.Sort != "" {
		values.Add("sort", p.Sort)
	}
	if p.FilterTitle != "" {
		values.Add("filter[title]", p.FilterTitle)
	}
	return values
}

func (p *GetBooksParams) FirstQuery() url.Values {
	values := url.Values{}
	if p.PageOffset == 0 {
		return nil
	}

	values.Add("page[limit]", fmt.Sprintf("%d", p.PageLimit))
	if p.Sort != "" {
		values.Add("sort", p.Sort)
	}
	if p.FilterTitle != "" {
		values.Add("filter[title]", p.FilterTitle)
	}
	return values
}

func (p *GetBooksParams) NextQuery(models []*book.Model) url.Values {
	values := url.Values{}
	if len(models) == 0 {
		return nil
	}

	values.Add("page[limit]", fmt.Sprintf("%d", p.PageLimit))
	values.Add("page[offset]", fmt.Sprintf("%d", p.PageOffset+p.PageLimit))

	if p.Sort != "" {
		values.Add("sort", p.Sort)
	}
	if p.FilterTitle != "" {
		values.Add("filter[title]", p.FilterTitle)
	}
	return values
}

func (p *GetBooksParams) PrevQuery() url.Values {
	values := url.Values{}
	if p.PageOffset == 0 {
		return nil
	}

	values.Add("page[limit]", fmt.Sprintf("%d", p.PageLimit))
	if p.PageOffset > p.PageLimit {
		offset := uint(0)
		if p.PageOffset > p.PageLimit {
			offset = p.PageOffset - p.PageLimit
		}
		values.Add("page[offset]", fmt.Sprintf("%d", offset))
	}
	if p.Sort != "" {
		values.Add("sort", p.Sort)
	}
	if p.FilterTitle != "" {
		values.Add("filter[title]", p.FilterTitle)
	}
	return values
}

type GetBookParams struct {
	ID int32 `uri:"id" binding:"required,min=1"`
}

func (p *GetBookParams) SelfQuery() url.Values {
	values := url.Values{}
	return values
}

type DeleteBookParams struct {
	ID int32 `uri:"id" binding:"required,min=1"`
}

func (p *DeleteBookParams) SelfQuery() url.Values {
	values := url.Values{}
	return values
}

type UpdateBookParams struct {
	ID int32 `uri:"id" binding:"required,min=1"`
}

func (p *UpdateBookParams) SelfQuery() url.Values {
	values := url.Values{}
	return values
}
