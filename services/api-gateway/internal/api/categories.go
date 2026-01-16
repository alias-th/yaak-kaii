package api

import (
	"net/http"
	"yaak-kaii/services/api-gateway/pkg/types"
	"yaak-kaii/shared/contracts"
	"yaak-kaii/shared/proto/product"

	"github.com/gin-gonic/gin"
)

func (app *Application) createCategory(c *gin.Context) {
	var body types.CreateCategoryRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		app.responseWithError(c, http.StatusBadRequest, err)
		return
	}

	attributes := make([]*product.CategoryAttribute, 0, len(body.Attributes))
	for _, attr := range body.Attributes {
		attributes = append(attributes, &product.CategoryAttribute{
			Key:      attr.Key,
			Label:    attr.Label,
			Type:     attr.Type,
			Required: attr.Required,
			Options:  attr.Options,
			Scope:    attr.Scope,
			Order:    attr.AxisOrder,
		})

	}

	cat, err := app.GrpcClients.Product.Client.CreateCategory(c, &product.CreateCategoryRequest{
		Name:        body.Name,
		Description: body.Description,
		Attributes:  attributes,
	})
	if err != nil {
		app.responseWithError(c, http.StatusBadRequest, err)
		return
	}

	res := contracts.APIResponse{Data: types.CreateCategoryResponse{
		CategoryID: cat.CategoryId,
	}}

	c.JSON(http.StatusCreated, res)
}
