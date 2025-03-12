package authorsv1

import (
	"fmt"
	"net/http"

	"github.com/1755/bookstore-api/api/routers"
	"github.com/1755/bookstore-api/api/schemas"
	"github.com/1755/bookstore-api/internal/author"
	"github.com/1755/bookstore-api/internal/lgr"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type ListRouterBuilder struct {
	authorService author.Service
	config        *routers.Config
}

func NewListRouterBuilder(authorService author.Service, config *routers.Config) *ListRouterBuilder {
	return &ListRouterBuilder{authorService, config}
}

func (r *ListRouterBuilder) Build(g *gin.RouterGroup) {
	g.GET("/", r.handler)
}

// ListAuthors	 	Getting the list of authors
//
//	@Summary		get authors
//	@Description	Returns a list of authors. Returns 500 if there is a database error. Returns 200 and the author data on success.
//	@Tags			authors
//	@Accept			json
//	@Produce		json
//	@Param			page[size]		query		int		false	"Limit (default: 25, min: 1, max: 100)"
//	@Param			page[cursor]	query		string	false	"Cursor (min: 1, max: 4)"
//	@Param			sort			query		string	false	"Sort (default: -updatedAt, allowed: -updatedAt, updatedAt, -name, name)"
//	@Param			filter[name]	query		string	false	"Filter name (default: empty, min: 2, max: 50)"
//	@Success		200				{object}	schemas.Document[[]schemas.Resource[schemas.Author]]
//	@Failure		400				{object}	schemas.DocumentError[schemas.Error]	"On validation error"
//	@Failure		500				{string}	schemas.DocumentError[schemas.Error]	"On internal server error"
//	@Router			/authors/ [get]
func (r *ListRouterBuilder) handler(c *gin.Context) {
	ctx := c.Request.Context()
	logger := lgr.GetLogger(ctx)

	params := GetAuthorsParams{
		PageLimit:  25,
		PageOffset: 0,
		Sort:       "-updatedAt",
	}

	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, schemas.DocumentError{
			Errors: []schemas.Error{
				{
					Title:  "Invalid query parameters",
					Detail: err.Error(),
				},
			},
		})
		return
	}

	models, err := r.authorService.GetMany(c.Request.Context(), params.ToGetManyParams())
	if err != nil {
		logger.Error("failed to get authors", zap.Error(err))

		c.JSON(http.StatusInternalServerError, schemas.InternalServerError)
		return
	}

	document := schemas.Document[[]schemas.Resource[schemas.Author]]{
		Links: &schemas.DocumentLink{
			Self: fmt.Sprintf("%s/v1/authors/?%s", r.config.BaseURL, params.SelfQuery().Encode()),
		},
		Data: make([]schemas.Resource[schemas.Author], len(models)),
	}
	if params.FirstQuery() != nil {
		document.Links.First = fmt.Sprintf("%s/v1/authors/?%s", r.config.BaseURL, params.FirstQuery().Encode())
	}
	if params.PrevQuery() != nil {
		document.Links.Prev = fmt.Sprintf("%s/v1/authors/?%s", r.config.BaseURL, params.PrevQuery().Encode())
	}
	if params.NextQuery(models) != nil {
		document.Links.Next = fmt.Sprintf("%s/v1/authors/?%s", r.config.BaseURL, params.NextQuery(models).Encode())
	}

	for i, model := range models {
		document.Data[i] = schemas.NewAuthorDocumentFromModel(model, r.config.BaseURL).Data
	}

	c.JSON(http.StatusOK, document)
}
