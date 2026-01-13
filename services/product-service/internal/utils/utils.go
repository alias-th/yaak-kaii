package utils

import (
	"fmt"
	"math"
	"strings"
	"yaak-kaii/services/product-service/pkg/types"
)

func GenerateSKU(category, color, size string, variantNo int) string {
	c := strings.ToUpper(category[:3])
	cl := strings.ToUpper(color[:2])
	sz := strings.ToUpper(size)

	return fmt.Sprintf("%s-%s-%s-%02d", c, cl, sz, variantNo)
}

func NewPagination(page, pageSize int, total int64) *types.Pagination {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	lastPage := int(math.Ceil(float64(total) / float64(pageSize)))

	var (
		nextPage *int
		prevPage *int
	)

	if page < lastPage {
		n := page + 1
		nextPage = &n
	}

	if page > 1 && page <= lastPage {
		p := page - 1
		prevPage = &p
	}

	return &types.Pagination{
		Page:     page,
		PageSize: pageSize,
		Total:    total,
		LastPage: lastPage,
		NextPage: nextPage,
		PrevPage: prevPage,
		HasNext:  nextPage != nil,
		HasPrev:  prevPage != nil,
	}
}

func ProductOrder(sort string) string {
	switch sort {
	case "oldest":
		return "products.created_at ASC, products.id ASC"
	case "name_asc":
		return "products.name ASC, products.id ASC"
	case "name_desc":
		return "products.name DESC, products.id DESC"
	case "price_asc":
		return `
			(
				SELECT MIN(pv.price)
				FROM product_variants pv
				WHERE pv.product_id = products.id
			) ASC NULLS LAST,
			products.id ASC
		`
	case "price_desc":
		return `
			(
				SELECT MIN(pv.price)
				FROM product_variants pv
				WHERE pv.product_id = products.id
			) DESC NULLS LAST,
			products.id DESC
		`
	default:
		return "products.created_at DESC, products.id DESC"
	}
}
