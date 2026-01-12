package api

import (
	"net/http"
	"yaak-kaii/services/api-gateway/pkg/types"
	"yaak-kaii/shared/contracts"
	"yaak-kaii/shared/proto/product"

	"github.com/gin-gonic/gin"
)

func (app *Application) createProduct(ctx *gin.Context) {
	userID := ctx.GetString("user_id")

	var reqBody types.CreateProductRequest

	err := ctx.ShouldBindJSON(&reqBody)
	if err != nil {
		app.responseWithError(ctx, 400, err)
		return
	}

	productRes, err := app.GrpcClients.Product.Client.CreateProduct(ctx, &product.CreateProductRequest{
		UserId:      userID,
		Name:        reqBody.Name,
		Description: reqBody.Description,
		Price:       reqBody.Price,
		CategoryId:  reqBody.CategoryID,
		Stock:       int32(reqBody.Stock),
		Attributes:  reqBody.Attributes,
	})
	if err != nil {
		app.responseWithError(ctx, http.StatusInternalServerError, err)
		return
	}

	res := contracts.APIResponse{
		Data: types.CreateProductResponse{
			ID: productRes.ProductId,
		},
	}

	ctx.JSON(http.StatusCreated, res)
}
