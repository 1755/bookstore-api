package author

import (
	"github.com/joomcode/errorx"
)

var author = errorx.NewNamespace("author")

var NotFoundError = author.NewType("NotFound", errorx.NotFound())
var InternalError = author.NewType("Internal")
