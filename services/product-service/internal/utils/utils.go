package utils

import (
	"fmt"
	"math"
	"strings"
	"yaak-kaii/services/product-service/pkg/types"
)

func GenerateSKU(category, color, size string) string {
	c := strings.ToUpper(category[:3])
	cl := strings.ToUpper(color[:2])
	sz := strings.ToUpper(size)

	return fmt.Sprintf("%s-%s-%s", c, cl, sz)
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
