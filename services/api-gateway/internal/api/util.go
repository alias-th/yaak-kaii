package api

import (
	"yaak-kaii/shared/contracts"

	"github.com/gin-gonic/gin"
)

func errorResponse(err error, code int) gin.H {
	return gin.H{
		"error": contracts.APIError{
			Code:    code,
			Message: err.Error(),
		},
	}
}
