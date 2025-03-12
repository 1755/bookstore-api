package authorsv1

import (
	"net/http"

	"github.com/1755/bookstore-api/api/routers"
	"github.com/1755/bookstore-api/api/schemas"
	"github.com/1755/bookstore-api/internal/author"
	"github.com/1755/bookstore-api/internal/lgr"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type CreateRouterBuilder struct {
	authorService author.Service
	config        *routers.Config
}

func NewCreateRouterBuilder(authorService author.Service, config *routers.Config) *CreateRouterBuilder {
	return &CreateRouterBuilder{authorService, config}
}

func (r *CreateRouterBuilder) Build(g *gin.RouterGroup) {
	g.POST("/", r.handler)
}

// CreateAuthor	 	Creating an author
//
//	@Summary		create author
//	@Description	Creates an author in the system.
//	@Tags			authors
//	@Accept			json
//	@Produce		json
//	@Param			body	body		schemas.Document[schemas.CreateResource[schemas.CreateAuthorRequest]]	true	"Author data in json api format"
//	@Success		201		{object}	schemas.Document[schemas.Resource[schemas.Author]]	"Author created successfully"
//	@Failure		400		{object}	schemas.DocumentError[schemas.Error]	"On validation error"
//	@Failure		500		{string}	schemas.DocumentError[schemas.Error]	"On internal server error"
//	@Router			/authors/ [post]
func (r *CreateRouterBuilder) handler(c *gin.Context) {
	ctx := c.Request.Context()
	logger := lgr.GetLogger(ctx)

	var requestDocument schemas.Document[schemas.CreateResource[schemas.CreateAuthorRequest]]
	if err := c.ShouldBindJSON(&requestDocument); err != nil {
		c.JSON(http.StatusBadRequest, schemas.DocumentError{
			Errors: schemas.NewValidationErrorsFromBindingError(err),
		})
		return
	}

	if requestDocument.Data.Type != "authors" {
		c.JSON(http.StatusBadRequest, schemas.DocumentError{
			Errors: []schemas.Error{
				{
					Title:  "Invalid body",
					Detail: "Type must be authors",
				},
			},
		})
		return
	}

	newModel := author.Model{
		Name: requestDocument.Data.Attributes.Name,
		Bio:  requestDocument.Data.Attributes.Bio,
	}

	model, err := r.authorService.Create(ctx, &newModel)
	if err != nil {
		logger.Error("failed to create author", zap.Error(err))

		c.JSON(http.StatusInternalServerError, schemas.DocumentError{
			Errors: []schemas.Error{schemas.InternalServerError},
		})
		return
	}

	document := schemas.NewAuthorDocumentFromModel(model, r.config.BaseURL)

	c.Header("Location", document.Links.Self)
	c.JSON(http.StatusCreated, document)
}
