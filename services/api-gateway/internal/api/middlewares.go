package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func guestMiddleware() gin.HandlerFunc {

	return func(c *gin.Context) {

		_, err := c.Cookie("guest_session")
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err, http.StatusUnauthorized))
			return
		}

		c.Next()
	}
}
