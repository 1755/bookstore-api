package authorsv1

import (
	"net/http"

	"github.com/1755/bookstore-api/api/routers"
	"github.com/1755/bookstore-api/api/schemas"
	"github.com/1755/bookstore-api/internal/author"
	"github.com/1755/bookstore-api/internal/lgr"
	"github.com/gin-gonic/gin"
	"github.com/joomcode/errorx"
	"go.uber.org/zap"
)

type GetRouterBuilder struct {
	authorService author.Service
	config        *routers.Config
}

func NewGetRouterBuilder(authorService author.Service, config *routers.Config) *GetRouterBuilder {
	return &GetRouterBuilder{authorService, config}
}

func (r *GetRouterBuilder) Build(g *gin.RouterGroup) {
	g.GET("/:id", r.handler)
}

// GetAuthor	 	Getting the author by its ID
//
//	@Summary		get author by id
//	@Description	Gets an author from the database by its ID
//	@Tags			authors
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"Author ID"
//	@Success		200	{object}	schemas.Document[schemas.Resource[schemas.Author]]
//	@Failure		400	{object}	schemas.DocumentError[schemas.Error]	"On validation error"
//	@Failure		404	{object}	schemas.DocumentError[schemas.Error]	"On author not found"
//	@Failure		500	{string}	schemas.DocumentError[schemas.Error]	"On internal server error"
//	@Router			/authors/{id} [get]
func (r *GetRouterBuilder) handler(c *gin.Context) {
	ctx := c.Request.Context()
	logger := lgr.GetLogger(ctx)

	var params GetAuthorParams
	if err := c.ShouldBindUri(&params); err != nil {
		c.JSON(http.StatusBadRequest, schemas.DocumentError{
			Errors: schemas.NewValidationErrorsFromBindingError(err),
		})
		return
	}

	model, err := r.authorService.GetByID(c.Request.Context(), author.ID(params.ID))
	if err != nil {
		if errorx.IsNotFound(err) {
			c.JSON(http.StatusNotFound, schemas.DocumentError{
				Errors: []schemas.Error{
					{
						Title:  "Author not found",
						Detail: "The requested author was not found",
					},
				},
			})
			return
		}

		logger.Error("failed to get author", zap.Error(err))

		c.JSON(http.StatusInternalServerError, schemas.DocumentError{
			Errors: []schemas.Error{schemas.InternalServerError},
		})
		return
	}

	document := schemas.NewAuthorDocumentFromModel(model, r.config.BaseURL)

	c.JSON(http.StatusOK, document)
}
