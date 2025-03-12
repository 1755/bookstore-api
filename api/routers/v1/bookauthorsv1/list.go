package bookauthorsv1

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

// ListBookAuthors	 	Getting the list of book authors
//
//	@Summary		get book authors
//	@Description	Returns a list of books. Returns 500 if there is a database error. Returns 200 and the book data on success.
//	@Tags			books
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"Book ID"
//	@Success		200				{object}	schemas.Document[[]schemas.Resource[schemas.Author]]
//	@Failure		400				{object}	schemas.DocumentError[schemas.Error]	"On validation error"
//	@Failure		500				{string}	schemas.DocumentError[schemas.Error]	"On internal server error"
//	@Router			/books/{id}/authors [get]
func (r *ListRouterBuilder) handler(c *gin.Context) {
	ctx := c.Request.Context()
	logger := lgr.GetLogger(ctx)

	var params ListBookAuthorsParams
	if err := c.ShouldBindUri(&params); err != nil {
		c.JSON(http.StatusBadRequest, schemas.DocumentError{
			Errors: schemas.NewValidationErrorsFromBindingError(err),
		})
		return
	}

	models, err := r.authorService.GetManyByBookID(c.Request.Context(), params.ID)
	if err != nil {
		logger.Error("failed to get book authors", zap.Error(err))

		c.JSON(http.StatusInternalServerError, schemas.InternalServerError)
		return
	}

	document := schemas.Document[[]schemas.Resource[schemas.Author]]{
		Links: &schemas.DocumentLink{
			Self: fmt.Sprintf("%s/v1/books/%d/authors/?%s", r.config.BaseURL, params.ID, params.SelfQuery().Encode()),
		},
		Data: make([]schemas.Resource[schemas.Author], len(models)),
	}

	for i, model := range models {
		document.Data[i] = schemas.NewAuthorDocumentFromModel(model, r.config.BaseURL).Data
	}

	c.JSON(http.StatusOK, document)
}
