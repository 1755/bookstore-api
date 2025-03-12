package authorsv1

import (
	"net/http"

	"github.com/1755/bookstore-api/api/schemas"
	"github.com/1755/bookstore-api/internal/author"
	"github.com/1755/bookstore-api/internal/lgr"
	"github.com/gin-gonic/gin"
	"github.com/joomcode/errorx"
	"go.uber.org/zap"
)

type DeleteRouterBuilder struct {
	authorService author.Service
}

func NewDeleteRouterBuilder(authorService author.Service) *DeleteRouterBuilder {
	return &DeleteRouterBuilder{authorService}
}

func (r *DeleteRouterBuilder) Build(g *gin.RouterGroup) {
	g.DELETE("/:id", r.handler)
}

// DeleteAuthor	 	Deleting the author by its ID
//
//	@Summary		delete author by id
//	@Description	Deletes an author from the system by its ID.
//	@Tags			authors
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"Author ID"
//	@Success		204	{object}	nil
//	@Failure		400	{object}	schemas.DocumentError[schemas.Error]	"On validation error"
//	@Failure		404	{object}	schemas.DocumentError[schemas.Error]	"On author not found"
//	@Failure		500	{string}	schemas.DocumentError[schemas.Error]	"On internal server error"
//	@Router			/authors/{id} [delete]
func (r *DeleteRouterBuilder) handler(c *gin.Context) {
	ctx := c.Request.Context()
	logger := lgr.GetLogger(ctx)

	var params DeleteAuthorParams
	if err := c.ShouldBindUri(&params); err != nil {
		c.JSON(http.StatusBadRequest, schemas.DocumentError{
			Errors: schemas.NewValidationErrorsFromBindingError(err),
		})
		return
	}

	if err := r.authorService.Delete(c.Request.Context(), author.ID(params.ID)); err != nil {
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

		logger.Error("failed to delete author", zap.Error(err))

		c.JSON(http.StatusInternalServerError, schemas.DocumentError{
			Errors: []schemas.Error{schemas.InternalServerError},
		})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
