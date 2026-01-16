package api

import (
	"fmt"
	"log"
	"net/http"
	"yaak-kaii/services/api-gateway/pkg/types"
	"yaak-kaii/shared/contracts"
	"yaak-kaii/shared/proto/product"

	"github.com/gin-gonic/gin"
)

func (app *Application) createCategories(c *gin.Context) {
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

func (app *Application) getCategories(c *gin.Context) {
	var body types.GetAllCategoriesRequest
	if err := c.ShouldBindQuery(&body); err != nil {
		app.responseWithError(c, http.StatusBadRequest, err)
		return
	}

	categoriesResp, err := app.GrpcClients.Product.Client.GetAllCategories(c, &product.GetAllCategoriesRequest{
		Limit: body.Limit,
		Page:  body.Page,
	})

	if err != nil {
		app.responseWithError(c, http.StatusBadRequest, err)
		return
	}

	log.Println(categoriesResp.Pagination, "categoriesResp")
	categories := make([]types.Category, 0, len(categoriesResp.Categories))
	for _, cat := range categoriesResp.Categories {
		attrs := make([]types.CategoryAttribute, 0, len(cat.Attributes))
		for _, attr := range cat.Attributes {
			opts := attr.Options
			if opts == nil {
				opts = []string{}
			}
			attrs = append(attrs, types.CategoryAttribute{
				Key:       attr.Key,
				Label:     attr.Label,
				Type:      attr.Type,
				Required:  attr.Required,
				Options:   opts,
				Scope:     attr.Scope,
				AxisOrder: attr.Order,
			})
		}

		categories = append(categories, types.Category{
			ID:          cat.CategoryId,
			Name:        cat.Name,
			Description: cat.Description,
			Attributes:  attrs,
		})
	}
	resp := types.GetAllCategoriesResponse{
		Categories: categories,
		Pagination: types.Pagination{
			Page:     categoriesResp.Pagination.Page,
			Limit:    categoriesResp.Pagination.Limit,
			Total:    categoriesResp.Pagination.Total,
			LastPage: categoriesResp.Pagination.LastPage,
			NextPage: categoriesResp.Pagination.NextPage,
			PrevPage: categoriesResp.Pagination.PrevPage,
			HasNext:  categoriesResp.Pagination.HasNext,
			HasPrev:  categoriesResp.Pagination.HasPrev,
		},
	}

	c.JSON(http.StatusOK, contracts.APIResponse{Data: resp})
}

func (app *Application) getCategory(c *gin.Context) {
	categoryID := c.Param("id")
	if categoryID == "" {
		app.responseWithError(c, http.StatusBadRequest, fmt.Errorf("invalid_category_id"))
		return
	}

	categoryResp, err := app.GrpcClients.Product.Client.GetCategory(c, &product.GetCategoryRequest{
		CategoryId: categoryID,
	})
	if err != nil {
		app.responseWithError(c, http.StatusBadRequest, err)
		return
	}

	attrs := make([]types.CategoryAttribute, 0, len(categoryResp.Category.Attributes))
	for _, attr := range categoryResp.Category.Attributes {
		opts := attr.Options
		if opts == nil {
			opts = []string{}
		}
		attrs = append(attrs, types.CategoryAttribute{
			Key:       attr.Key,
			Label:     attr.Label,
			Type:      attr.Type,
			Required:  attr.Required,
			Options:   opts,
			Scope:     attr.Scope,
			AxisOrder: attr.Order,
		})
	}

	resp := types.Category{
		ID:          categoryResp.Category.CategoryId,
		Name:        categoryResp.Category.Name,
		Description: categoryResp.Category.Description,
		Attributes:  attrs,
	}

	c.JSON(http.StatusOK, contracts.APIResponse{Data: resp})
}
