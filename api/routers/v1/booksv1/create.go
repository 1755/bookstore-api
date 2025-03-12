package booksv1

import (
	"net/http"

	"github.com/1755/bookstore-api/api/routers"
	"github.com/1755/bookstore-api/api/schemas"
	"github.com/1755/bookstore-api/internal/book"
	"github.com/1755/bookstore-api/internal/lgr"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type CreateRouterBuilder struct {
	bookService book.Service
	config      *routers.Config
}

func NewCreateRouterBuilder(bookService book.Service, config *routers.Config) *CreateRouterBuilder {
	return &CreateRouterBuilder{bookService, config}
}

func (r *CreateRouterBuilder) Build(g *gin.RouterGroup) {
	g.POST("/", r.handler)
}

// CreateBook	 	Creating a book
//
//	@Summary		create book
//	@Description	Creates a book in the system.
//	@Tags			books
//	@Accept			json
//	@Produce		json
//	@Param			body	body		schemas.Document[schemas.CreateResource[schemas.CreateBookRequest]]	true	"Book data in json api format"
//	@Success		201		{object}	schemas.Document[schemas.Resource[schemas.Book]]	"Book created successfully"
//	@Failure		400		{object}	schemas.DocumentError[schemas.Error]	"On validation error"
//	@Failure		500		{string}	schemas.DocumentError[schemas.Error]	"On internal server error"
//	@Router			/books/ [post]
func (r *CreateRouterBuilder) handler(c *gin.Context) {
	ctx := c.Request.Context()
	logger := lgr.GetLogger(ctx)

	var requestDocument schemas.Document[schemas.CreateResource[schemas.CreateBookRequest]]
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

	newModel := book.Model{
		Title:         requestDocument.Data.Attributes.Title,
		Summary:       requestDocument.Data.Attributes.Summary,
		PublishedYear: requestDocument.Data.Attributes.PublishedYear,
	}

	model, err := r.bookService.Create(ctx, &newModel)
	if err != nil {
		logger.Error("failed to create book", zap.Error(err))

		c.JSON(http.StatusInternalServerError, schemas.DocumentError{
			Errors: []schemas.Error{schemas.InternalServerError},
		})
		return
	}

	document := schemas.NewBookDocumentFromModel(model, r.config.BaseURL)

	c.Header("Location", document.Links.Self)
	c.JSON(http.StatusCreated, document)
}
