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

func (app *Application) listProducts(ctx *gin.Context) {
	var listProducts types.ListProductRequest
	if err := ctx.ShouldBindQuery(&listProducts); err != nil {
		app.responseWithError(ctx, 400, err)
		return
	}
	var maxPrice int64
	var minPrice int64
	if listProducts.MaxPrice != nil {
		maxPrice = *listProducts.MaxPrice
	}
	if listProducts.MinPrice != nil {
		minPrice = *listProducts.MinPrice
	}

	res, err := app.GrpcClients.Product.Client.ListProducts(ctx, &product.ListProductsRequest{
		Q:           listProducts.Q,
		Limit:       listProducts.Limit,
		Page:        listProducts.Page,
		CategoryIds: listProducts.CategoryID,
		MinPrice:    minPrice,
		MaxPrice:    maxPrice,
		Sort:        listProducts.Sort,
	})
	if err != nil {
		app.responseWithError(ctx, http.StatusInternalServerError, err)
		return
	}

	products := []types.ProductResponse{}
	for _, p := range res.Products {
		variants := []types.ProductVariant{}
		for _, v := range p.Variants {
			variants = append(variants, types.ProductVariant{
				SKU:        v.Sku,
				Price:      v.Price,
				Stock:      v.Stock,
				Attributes: v.Attributes,
			})
		}
		products = append(products, types.ProductResponse{
			ID:          p.Id,
			Name:        p.Name,
			ShopID:      p.ShopId,
			Description: p.Description,
			Status:      p.Status,
			Variants:    variants,
			Images:      p.Images,
			Category: types.ProductCategory{
				CategoryID: p.Category.CategoryId,
				Name:       p.Category.Name,
			},
		})
	}

	pagination := types.Pagination{
		Page:     res.Pagination.Page,
		Limit:    res.Pagination.Limit,
		Total:    res.Pagination.Total,
		LastPage: res.Pagination.LastPage,
		NextPage: res.Pagination.NextPage,
		PrevPage: res.Pagination.PrevPage,
		HasNext:  res.Pagination.HasNext,
		HasPrev:  res.Pagination.HasPrev,
	}

	mapRes := contracts.APIResponse{
		Data: types.ListProductResponse{
			Products:   products,
			Pagination: pagination,
		},
	}

	ctx.JSON(200, mapRes)
}
