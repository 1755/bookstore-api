package schemas

import (
	"encoding/json"
	"fmt"

	"github.com/go-playground/validator/v10"
)

type Error struct {
	Title  string `json:"title"`
	Detail string `json:"detail"`
}

var InternalServerError = Error{
	Title:  "Internal server error",
	Detail: "An internal server error occurred",
}

var NewValidationErrorsFromBindingError = func(err error) []Error {
	var errors []Error

	switch e := err.(type) {
	case *json.SyntaxError:
		errors = append(errors, Error{
			Title:  "Invalid JSON syntax",
			Detail: fmt.Sprintf("JSON syntax error at position %d: %s", e.Offset, e.Error()),
		})
	case *json.UnmarshalTypeError:
		errors = append(errors, Error{
			Title:  "Invalid type",
			Detail: fmt.Sprintf("Expected %s for field '%s' but got %s", e.Type, e.Field, e.Value),
		})
	case validator.ValidationErrors:
		for _, ve := range e {
			switch ve.Tag() {
			case "required":
				errors = append(errors, Error{
					Title:  "Invalid body",
					Detail: fmt.Sprintf("%s is required", ve.Field()),
				})
			case "max":
				errors = append(errors, Error{
					Title:  "Invalid body",
					Detail: fmt.Sprintf("%s exceeds maximum length", ve.Field()),
				})
			case "min":
				errors = append(errors, Error{
					Title:  "Invalid body",
					Detail: fmt.Sprintf("%s below minimum length", ve.Field()),
				})
			default:
				errors = append(errors, Error{
					Title:  "Invalid body",
					Detail: fmt.Sprintf("%s is invalid by tag %s", ve.Field(), ve.Tag()),
				})
			}
		}
	default:
		errors = append(errors, Error{
			Title:  "Invalid request body",
			Detail: "The request body could not be parsed. Please check the format and try again.",
		})
	}

	return errors
}
