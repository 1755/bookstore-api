package booksv1

import (
	"net/http"

	"github.com/1755/bookstore-api/api/schemas"
	"github.com/1755/bookstore-api/internal/book"
	"github.com/1755/bookstore-api/internal/lgr"
	"github.com/gin-gonic/gin"
	"github.com/joomcode/errorx"
	"go.uber.org/zap"
)

type DeleteRouterBuilder struct {
	bookService book.Service
}

func NewDeleteRouterBuilder(bookService book.Service) *DeleteRouterBuilder {
	return &DeleteRouterBuilder{bookService}
}

func (r *DeleteRouterBuilder) Build(g *gin.RouterGroup) {
	g.DELETE("/:id", r.handler)
}

// DeleteBook	 	Deleting the book by its ID
//
//	@Summary		delete book by id
//	@Description	Deletes a book from the system by its ID.
//	@Tags			books
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"Book ID"
//	@Success		204	{object}	nil
//	@Failure		400	{object}	schemas.DocumentError[schemas.Error]	"On validation error"
//	@Failure		404	{object}	schemas.DocumentError[schemas.Error]	"On book not found"
//	@Failure		500	{string}	schemas.DocumentError[schemas.Error]	"On internal server error"
//	@Router			/books/{id} [delete]
func (r *DeleteRouterBuilder) handler(c *gin.Context) {
	ctx := c.Request.Context()
	logger := lgr.GetLogger(ctx)

	var params DeleteBookParams
	if err := c.ShouldBindUri(&params); err != nil {
		c.JSON(http.StatusBadRequest, schemas.DocumentError{
			Errors: schemas.NewValidationErrorsFromBindingError(err),
		})
		return
	}

	if err := r.bookService.Delete(c.Request.Context(), book.ID(params.ID)); err != nil {
		if errorx.IsNotFound(err) {
			c.JSON(http.StatusNotFound, schemas.DocumentError{
				Errors: []schemas.Error{
					{
						Title:  "Book not found",
						Detail: "The requested book was not found",
					},
				},
			})
			return
		}

		logger.Error("failed to delete book", zap.Error(err))

		c.JSON(http.StatusInternalServerError, schemas.DocumentError{
			Errors: []schemas.Error{schemas.InternalServerError},
		})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
