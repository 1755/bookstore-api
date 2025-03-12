package schemas

type Resource[T any] struct {
	Type       string `json:"type"`
	ID         string `json:"id"`
	Attributes T      `json:"attributes"`
	Links      *Link  `json:"links,omitempty"`
}

type CreateResource[T any] struct {
	Type       string `json:"type" binding:"required"`
	Attributes T      `json:"attributes" binding:"required"`
}

type UpdateResource[T any] struct {
	Type       string `json:"type" binding:"required"`
	ID         string `json:"id" binding:"required"`
	Attributes T      `json:"attributes" binding:"required"`
}
