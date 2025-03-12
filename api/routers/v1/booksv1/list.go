package booksv1

import (
	"fmt"
	"net/http"

	"github.com/1755/bookstore-api/api/routers"
	"github.com/1755/bookstore-api/api/schemas"
	"github.com/1755/bookstore-api/internal/book"
	"github.com/1755/bookstore-api/internal/lgr"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type ListRouterBuilder struct {
	bookService book.Service
	config      *routers.Config
}

func NewListRouterBuilder(bookService book.Service, config *routers.Config) *ListRouterBuilder {
	return &ListRouterBuilder{bookService, config}
}

func (r *ListRouterBuilder) Build(g *gin.RouterGroup) {
	g.GET("/", r.handler)
}

// ListBooks	 	Getting the list of books
//
//	@Summary		get books
//	@Description	Returns a list of books. Returns 500 if there is a database error. Returns 200 and the book data on success.
//	@Tags			books
//	@Accept			json
//	@Produce		json
//	@Param			page[size]		query		int		false	"Limit (default: 25, min: 1, max: 100)"
//	@Param			page[cursor]	query		string	false	"Cursor (min: 1, max: 4)"
//	@Param			sort			query		string	false	"Sort (default: -updatedAt, allowed: -updatedAt, updatedAt, -title, title)"
//	@Param			filter[title]	query		string	false	"Filter title (default: empty, min: 2, max: 50)"
//	@Success		200				{object}	schemas.Document[[]schemas.Resource[schemas.Book]]
//	@Failure		400				{object}	schemas.DocumentError[schemas.Error]	"On validation error"
//	@Failure		500				{string}	schemas.DocumentError[schemas.Error]	"On internal server error"
//	@Router			/books/ [get]
func (r *ListRouterBuilder) handler(c *gin.Context) {
	ctx := c.Request.Context()
	logger := lgr.GetLogger(ctx)

	params := GetBooksParams{
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

	models, err := r.bookService.GetMany(c.Request.Context(), params.ToGetManyParams())
	if err != nil {
		logger.Error("failed to get books", zap.Error(err))

		c.JSON(http.StatusInternalServerError, schemas.InternalServerError)
		return
	}

	document := schemas.Document[[]schemas.Resource[schemas.Book]]{
		Links: &schemas.DocumentLink{
			Self: fmt.Sprintf("%s/v1/books/?%s", r.config.BaseURL, params.SelfQuery().Encode()),
		},
		Data: make([]schemas.Resource[schemas.Book], len(models)),
	}
	if params.FirstQuery() != nil {
		document.Links.First = fmt.Sprintf("%s/v1/books/?%s", r.config.BaseURL, params.FirstQuery().Encode())
	}
	if params.PrevQuery() != nil {
		document.Links.Prev = fmt.Sprintf("%s/v1/books/?%s", r.config.BaseURL, params.PrevQuery().Encode())
	}
	if params.NextQuery(models) != nil {
		document.Links.Next = fmt.Sprintf("%s/v1/books/?%s", r.config.BaseURL, params.NextQuery(models).Encode())
	}

	for i, model := range models {
		document.Data[i] = schemas.NewBookDocumentFromModel(model, r.config.BaseURL).Data
	}

	c.JSON(http.StatusOK, document)
}
