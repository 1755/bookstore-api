package authorsv1

import (
	"fmt"
	"net/http"

	"github.com/1755/bookstore-api/api/routers"
	"github.com/1755/bookstore-api/api/schemas"
	"github.com/1755/bookstore-api/internal/author"
	"github.com/1755/bookstore-api/internal/lgr"
	"github.com/gin-gonic/gin"
	"github.com/joomcode/errorx"
	"go.uber.org/zap"
)

type UpdateRouterBuilder struct {
	authorService author.Service
	config        *routers.Config
}

func NewUpdateRouterBuilder(authorService author.Service, config *routers.Config) *UpdateRouterBuilder {
	return &UpdateRouterBuilder{authorService, config}
}

func (r *UpdateRouterBuilder) Build(g *gin.RouterGroup) {
	g.PATCH("/:id", r.handler)
}

// UpdateAuthor	 	Updating the author by its ID
//
//	@Summary		update author by id
//	@Description	Updates an author in the system by its ID.
//	@Tags			authors
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int																		true	"Author ID"
//	@Param			body	body		schemas.Document[schemas.UpdateResource[schemas.UpdateAuthorRequest]]	true	"Author data"
//	@Success		200		{object}	schemas.Document[schemas.Resource[schemas.Author]]						"Author updated"
//	@Failure		400		{object}	schemas.DocumentError[schemas.Error]									"On	validation	error"
//	@Failure		404		{object}	schemas.DocumentError[schemas.Error]									"On	author not found"
//	@Failure		500		{string}	schemas.DocumentError[schemas.Error]									"On	internal server	error"
//	@Router			/authors/{id} [patch]
func (r *UpdateRouterBuilder) handler(c *gin.Context) {
	ctx := c.Request.Context()
	logger := lgr.GetLogger(ctx)

	var params UpdateAuthorParams
	if err := c.ShouldBindUri(&params); err != nil {
		c.JSON(http.StatusBadRequest, schemas.DocumentError{
			Errors: schemas.NewValidationErrorsFromBindingError(err),
		})
		return
	}

	var requestDocument schemas.Document[schemas.UpdateResource[schemas.UpdateAuthorRequest]]
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

	set := []author.UpdateField{}
	if requestDocument.Data.Attributes.Name != nil {
		set = append(set, author.UpdateNameField(*requestDocument.Data.Attributes.Name))
	}
	if requestDocument.Data.Attributes.Bio != nil {
		set = append(set, author.UpdateBioField(*requestDocument.Data.Attributes.Bio))
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

	model, err := r.authorService.Update(c.Request.Context(), author.ID(params.ID), set...)
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

		logger.Error("failed to update author", zap.Error(err))

		c.JSON(http.StatusInternalServerError, schemas.DocumentError{
			Errors: []schemas.Error{schemas.InternalServerError},
		})
		return
	}

	document := schemas.NewAuthorDocumentFromModel(model, r.config.BaseURL)

	c.JSON(http.StatusOK, document)
}
