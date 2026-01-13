package api

import (
	"fmt"
	"net/http"
	"path"
	"strings"
	"yaak-kaii/services/api-gateway/pkg/types"
	"yaak-kaii/shared/contracts"
	"yaak-kaii/shared/proto/product"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

type UploadedImage struct {
	ProjectID string
	Urls      []string
}

func (app *Application) uploadProductImages(c *gin.Context) {
	productID := c.Param("id")

	form, err := c.MultipartForm()
	if err != nil {
		app.responseWithError(c, http.StatusBadRequest, err)
		return
	}

	files := form.File["images"]
	if len(files) == 0 {
		app.responseWithError(c, http.StatusBadRequest, fmt.Errorf("no images provided"))
		return
	}

	if strings.TrimSpace(productID) == "" {
		app.responseWithError(c, http.StatusBadRequest, fmt.Errorf("product ID is required"))
		return
	}

	_, err = app.GrpcClients.Product.Client.GetProductByID(c, &product.GetProductByIDRequest{
		ProductId: productID,
	})
	if err != nil {
		app.responseWithError(c, http.StatusBadRequest, err)
		return
	}

	urls := make([]string, 0, len(files))
	for _, fh := range files {
		// 5MB limit
		if fh.Size > 5<<20 {
			app.responseWithError(c, http.StatusBadRequest, fmt.Errorf("file too large: %s", fh.Filename))
			return
		}

		ct := strings.ToLower(strings.TrimSpace(fh.Header.Get("Content-Type")))
		if ct == "" {
			ct = "application/octet-stream"
		}
		if !isAllowedImageCT(ct) {
			app.responseWithError(c, http.StatusBadRequest, fmt.Errorf("unsupported content type: %s (%s)", fh.Filename, ct))
			return
		}

		f, err := fh.Open()
		if err != nil {
			app.responseWithError(c, http.StatusBadRequest, fmt.Errorf("cannot open file: %s", fh.Filename))
			return
		}

		ext := strings.ToLower(path.Ext(fh.Filename))
		if ext == "" {
			ext = guessExtFromContentType(ct)
		}

		key := fmt.Sprintf("products/%s/%s%s", productID, uuid.NewString(), ext)

		err = app.S3Uploader.Upload(f, key, ct)
		if err != nil {
			app.responseWithError(c, http.StatusInternalServerError, fmt.Errorf("failed to upload file: %s", fh.Filename))
			return
		}

		_ = f.Close()

		urls = append(urls, key)

	}

	_, err = app.GrpcClients.Product.Client.UploadProductImages(c, &product.UploadProductImagesRequest{
		ProductId: productID,
		ImageUrls: urls,
	})
	if err != nil {
		app.responseWithError(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":    "images uploaded successfully",
		"image_urls": urls,
	})

}

func isAllowedImageCT(ct string) bool {
	switch ct {
	case "image/jpeg", "image/png", "image/webp":
		return true
	default:
		return false
	}
}

func guessExtFromContentType(ct string) string {
	switch ct {
	case "image/jpeg":
		return ".jpg"
	case "image/png":
		return ".png"
	case "image/webp":
		return ".webp"
	default:
		return ""
	}
}
