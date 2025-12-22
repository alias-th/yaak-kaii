package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (app *Application) guestMiddleware() gin.HandlerFunc {

	return func(c *gin.Context) {

		_, err := c.Cookie("guest_session")
		if err != nil {
			app.responseWithError(c, http.StatusUnauthorized, err)
			return
		}

		c.Next()
	}
}
