package api

import (
	"errors"
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func (app *Application) responseWithError(ctx *gin.Context, code int, err error) {
	log.Printf("Error: %v", err)
	ctx.JSON(code, gin.H{
		"error": err.Error(),
	})
}

func (app *Application) formatDate(timestamp *int64) (at string, in int64, err error) {
	if timestamp == nil {
		return "", 0, errors.New("expires_at timestamp is required")
	}

	expiresAt := time.Unix(*timestamp, 0)
	expiresIn := int64(time.Until(expiresAt).Seconds())

	return expiresAt.Format(time.RFC3339), expiresIn, nil
}
