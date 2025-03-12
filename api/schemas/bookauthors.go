package schemas

type CreateBookAuthorsRequest struct {
	AuthorID int32 `json:"author_id" binding:"required,min=1"`
}
