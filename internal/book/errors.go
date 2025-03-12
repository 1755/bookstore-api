package book

import (
	"github.com/joomcode/errorx"
)

var book = errorx.NewNamespace("book")

var NotFoundError = book.NewType("NotFound", errorx.NotFound())
var InternalError = book.NewType("Internal")
