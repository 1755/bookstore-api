package authorsv1

import (
	"fmt"
	"net/url"

	"github.com/1755/bookstore-api/internal/author"
)

type GetAuthorsParams struct {
	PageLimit  uint   `form:"page[limit]" binding:"required,min=1,max=100"`
	PageOffset uint   `form:"page[offset]" binding:"omitempty"`
	Sort       string `form:"sort" binding:"omitempty,oneof=name -name updatedAt -updatedAt"`
	FilterName string `form:"filter[name]" binding:"omitempty,min=2,max=50"`
}

var mapSortField = map[string]string{
	"name":      "name",
	"updatedAt": "updated_at",
}

func (p *GetAuthorsParams) ToGetManyParams() *author.GetManyParams {
	var sort *author.GetManyParamsSort = nil
	if p.Sort != "" {
		sort = &author.GetManyParamsSort{}
		direction := "asc"
		value := p.Sort
		if p.Sort[0] == '-' {
			direction = "desc"
			value = p.Sort[1:]
		}

		field, ok := mapSortField[value]
		if ok {
			sort = &author.GetManyParamsSort{
				Field:     field,
				Direction: direction,
			}
		}
	}

	return &author.GetManyParams{
		Offset:     p.PageOffset,
		Limit:      p.PageLimit,
		Sort:       sort,
		FilterName: p.FilterName,
	}
}

func (p *GetAuthorsParams) SelfQuery() url.Values {
	values := url.Values{}
	values.Add("page[limit]", fmt.Sprintf("%d", p.PageLimit))
	if p.PageOffset > 0 {
		values.Add("page[offset]", fmt.Sprintf("%d", p.PageOffset))
	}
	if p.Sort != "" {
		values.Add("sort", p.Sort)
	}
	if p.FilterName != "" {
		values.Add("filter[name]", p.FilterName)
	}
	return values
}

func (p *GetAuthorsParams) FirstQuery() url.Values {
	values := url.Values{}
	if p.PageOffset == 0 {
		return nil
	}

	values.Add("page[limit]", fmt.Sprintf("%d", p.PageLimit))
	if p.Sort != "" {
		values.Add("sort", p.Sort)
	}
	if p.FilterName != "" {
		values.Add("filter[name]", p.FilterName)
	}
	return values
}

func (p *GetAuthorsParams) NextQuery(models []*author.Model) url.Values {
	values := url.Values{}
	if len(models) == 0 {
		return nil
	}

	values.Add("page[limit]", fmt.Sprintf("%d", p.PageLimit))
	values.Add("page[offset]", fmt.Sprintf("%d", p.PageOffset+p.PageLimit))

	if p.Sort != "" {
		values.Add("sort", p.Sort)
	}
	if p.FilterName != "" {
		values.Add("filter[name]", p.FilterName)
	}
	return values
}

func (p *GetAuthorsParams) PrevQuery() url.Values {
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
	if p.FilterName != "" {
		values.Add("filter[name]", p.FilterName)
	}
	return values
}

type GetAuthorParams struct {
	ID int32 `uri:"id" binding:"required,min=1"`
}

func (p *GetAuthorParams) SelfQuery() url.Values {
	values := url.Values{}
	return values
}

type DeleteAuthorParams struct {
	ID int32 `uri:"id" binding:"required,min=1"`
}

func (p *DeleteAuthorParams) SelfQuery() url.Values {
	values := url.Values{}
	return values
}

type UpdateAuthorParams struct {
	ID int32 `uri:"id" binding:"required,min=1"`
}

func (p *UpdateAuthorParams) SelfQuery() url.Values {
	values := url.Values{}
	return values
}
