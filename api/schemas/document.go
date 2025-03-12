package schemas

type Document[T any] struct {
	Links *DocumentLink `json:"links,omitempty"`
	Data  T             `json:"data"`
}

type DocumentError struct {
	Errors []Error `json:"errors"`
}
