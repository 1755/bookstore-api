package health

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (r *RouterBuilder) getHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.String(http.StatusOK, "healthy")
	}
}
