package booksv1

import (
	"net/http"

	"github.com/1755/bookstore-api/api/routers"
	"github.com/1755/bookstore-api/api/schemas"
	"github.com/1755/bookstore-api/internal/book"
	"github.com/1755/bookstore-api/internal/lgr"
	"github.com/gin-gonic/gin"
	"github.com/joomcode/errorx"
	"go.uber.org/zap"
)

type GetRouterBuilder struct {
	bookService book.Service
	config      *routers.Config
}

func NewGetRouterBuilder(bookService book.Service, config *routers.Config) *GetRouterBuilder {
	return &GetRouterBuilder{bookService, config}
}

func (r *GetRouterBuilder) Build(g *gin.RouterGroup) {
	g.GET("/:id", r.handler)
}

// GetBook	 	Getting the book by its ID
//
//	@Summary		get book by id
//	@Description	Gets a book from the database by its ID
//	@Tags			books
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"Book ID"
//	@Success		200	{object}	schemas.Document[schemas.Resource[schemas.Book]]
//	@Failure		400	{object}	schemas.DocumentError[schemas.Error]	"On validation error"
//	@Failure		404	{object}	schemas.DocumentError[schemas.Error]	"On book not found"
//	@Failure		500	{string}	schemas.DocumentError[schemas.Error]	"On internal server error"
//	@Router			/books/{id} [get]
func (r *GetRouterBuilder) handler(c *gin.Context) {
	ctx := c.Request.Context()
	logger := lgr.GetLogger(ctx)

	var params GetBookParams
	if err := c.ShouldBindUri(&params); err != nil {
		c.JSON(http.StatusBadRequest, schemas.DocumentError{
			Errors: schemas.NewValidationErrorsFromBindingError(err),
		})
		return
	}

	model, err := r.bookService.GetByID(c.Request.Context(), book.ID(params.ID))
	if err != nil {
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

		logger.Error("failed to get book", zap.Error(err))

		c.JSON(http.StatusInternalServerError, schemas.DocumentError{
			Errors: []schemas.Error{schemas.InternalServerError},
		})
		return
	}

	document := schemas.NewBookDocumentFromModel(model, r.config.BaseURL)

	c.JSON(http.StatusOK, document)
}
