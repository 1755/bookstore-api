package bookauthorsv1

import (
	"net/http"

	"github.com/1755/bookstore-api/api/routers"
	"github.com/1755/bookstore-api/api/schemas"
	"github.com/1755/bookstore-api/internal/author"
	"github.com/1755/bookstore-api/internal/book"
	"github.com/1755/bookstore-api/internal/lgr"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type CreateRouterBuilder struct {
	bookService   book.Service
	authorService author.Service
	config        *routers.Config
}

func NewCreateRouterBuilder(bookService book.Service, authorService author.Service, config *routers.Config) *CreateRouterBuilder {
	return &CreateRouterBuilder{bookService, authorService, config}
}

func (r *CreateRouterBuilder) Build(g *gin.RouterGroup) {
	g.POST("/", r.handler)
}

// CreateBookAuthors	 	Creating a book author
//
//	@Summary		create book author
//	@Description	Creates a book author in the system.
//	@Tags			books
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"Book ID"
//	@Param			body	body		schemas.Document[schemas.CreateResource[schemas.CreateBookAuthorsRequest]]	true	"Book data in json api format"
//	@Success		201		{object}	schemas.Document[schemas.Resource[schemas.Book]]	"Book author created successfully"
//	@Failure		400		{object}	schemas.DocumentError[schemas.Error]	"On validation error"
//	@Failure		500		{string}	schemas.DocumentError[schemas.Error]	"On internal server error"
//	@Router			/books/{id}/authors [post]
func (r *CreateRouterBuilder) handler(c *gin.Context) {
	ctx := c.Request.Context()
	logger := lgr.GetLogger(ctx)

	var params CreateBookAuthorsParams
	if err := c.ShouldBindUri(&params); err != nil {
		c.JSON(http.StatusBadRequest, schemas.DocumentError{
			Errors: schemas.NewValidationErrorsFromBindingError(err),
		})
		return
	}

	var requestDocument schemas.Document[schemas.CreateResource[schemas.CreateBookAuthorsRequest]]
	if err := c.ShouldBindJSON(&requestDocument); err != nil {
		c.JSON(http.StatusBadRequest, schemas.DocumentError{
			Errors: schemas.NewValidationErrorsFromBindingError(err),
		})
		return
	}

	if requestDocument.Data.Type != "book_authors" {
		c.JSON(http.StatusBadRequest, schemas.DocumentError{
			Errors: []schemas.Error{
				{
					Title:  "Invalid body",
					Detail: "Type must be book_authors",
				},
			},
		})
		return
	}

	err := r.bookService.LinkAuthor(ctx, book.ID(params.ID), requestDocument.Data.Attributes.AuthorID)
	if err != nil {
		logger.Error("failed to create book author", zap.Error(err))

		c.JSON(http.StatusInternalServerError, schemas.DocumentError{
			Errors: []schemas.Error{schemas.InternalServerError},
		})
		return
	}

	author, err := r.authorService.GetByID(ctx, author.ID(requestDocument.Data.Attributes.AuthorID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, schemas.DocumentError{
			Errors: []schemas.Error{schemas.InternalServerError},
		})
		return
	}

	document := schemas.NewAuthorDocumentFromModel(author, r.config.BaseURL)

	c.Header("Location", document.Links.Self)
	c.JSON(http.StatusCreated, document)
}
