package api

import "github.com/gin-gonic/gin"

func (app *Application) createProduct(ctx *gin.Context) {
	// Implement the logic to handle product creation
	ctx.JSON(200, gin.H{"message": "createProduct endpoint"})
}
