package booksv1

import (
	"fmt"
	"net/http"

	"github.com/1755/bookstore-api/api/routers"
	"github.com/1755/bookstore-api/api/schemas"
	"github.com/1755/bookstore-api/internal/book"
	"github.com/1755/bookstore-api/internal/lgr"
	"github.com/gin-gonic/gin"
	"github.com/joomcode/errorx"
	"go.uber.org/zap"
)

type UpdateRouterBuilder struct {
	bookService book.Service
	config      *routers.Config
}

func NewUpdateRouterBuilder(bookService book.Service, config *routers.Config) *UpdateRouterBuilder {
	return &UpdateRouterBuilder{bookService, config}
}

func (r *UpdateRouterBuilder) Build(g *gin.RouterGroup) {
	g.PATCH("/:id", r.handler)
}

// UpdateBook	 	Updating the book by its ID
//
//	@Summary		update book by id
//	@Description	Updates a book in the system by its ID.
//	@Tags			books
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int																		true	"Book ID"
//	@Param			body	body		schemas.Document[schemas.UpdateResource[schemas.UpdateBookRequest]]	true	"Book data"
//	@Success		200		{object}	schemas.Document[schemas.Resource[schemas.Book]]						"Book updated"
//	@Failure		400		{object}	schemas.DocumentError[schemas.Error]									"On	validation	error"
//	@Failure		404		{object}	schemas.DocumentError[schemas.Error]									"On	book not found"
//	@Failure		500		{string}	schemas.DocumentError[schemas.Error]									"On	internal server	error"
//	@Router			/books/{id} [patch]
func (r *UpdateRouterBuilder) handler(c *gin.Context) {
	ctx := c.Request.Context()
	logger := lgr.GetLogger(ctx)

	var params UpdateBookParams
	if err := c.ShouldBindUri(&params); err != nil {
		c.JSON(http.StatusBadRequest, schemas.DocumentError{
			Errors: schemas.NewValidationErrorsFromBindingError(err),
		})
		return
	}

	var requestDocument schemas.Document[schemas.UpdateResource[schemas.UpdateBookRequest]]
	if err := c.ShouldBindJSON(&requestDocument); err != nil {
		c.JSON(http.StatusBadRequest, schemas.DocumentError{
			Errors: schemas.NewValidationErrorsFromBindingError(err),
		})
		return
	}

	if requestDocument.Data.Type != "books" {
		c.JSON(http.StatusBadRequest, schemas.DocumentError{
			Errors: []schemas.Error{
				{
					Title:  "Invalid body",
					Detail: "Type must be books",
				},
			},
		})
		return
	}

	if requestDocument.Data.ID != fmt.Sprintf("%d", params.ID) {
		c.JSON(http.StatusBadRequest, schemas.DocumentError{
			Errors: []schemas.Error{
				{
					Title:  "Invalid body",
					Detail: "Id in body must match id in path",
				},
			},
		})
		return
	}

	set := []book.UpdateField{}
	if requestDocument.Data.Attributes.Title != nil {
		set = append(set, book.UpdateTitleField(*requestDocument.Data.Attributes.Title))
	}
	if requestDocument.Data.Attributes.Summary != nil {
		set = append(set, book.UpdateSummaryField(*requestDocument.Data.Attributes.Summary))
	}
	if requestDocument.Data.Attributes.PublishedYear != nil {
		set = append(set, book.UpdatePublishedYearField(*requestDocument.Data.Attributes.PublishedYear))
	}

	if len(set) == 0 {
		c.JSON(http.StatusBadRequest, schemas.DocumentError{
			Errors: []schemas.Error{
				{
					Title:  "Invalid body",
					Detail: "No fields to update",
				},
			},
		})
		return
	}

	model, err := r.bookService.Update(c.Request.Context(), book.ID(params.ID), set...)
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

		logger.Error("failed to update book", zap.Error(err))

		c.JSON(http.StatusInternalServerError, schemas.DocumentError{
			Errors: []schemas.Error{schemas.InternalServerError},
		})
		return
	}

	document := schemas.NewBookDocumentFromModel(model, r.config.BaseURL)

	c.JSON(http.StatusOK, document)
}
