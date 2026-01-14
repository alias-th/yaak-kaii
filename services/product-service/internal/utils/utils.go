package utils

import (
	"encoding/json"
	"fmt"
	"math"
	"strings"
	"yaak-kaii/services/product-service/pkg/types"

	"gorm.io/datatypes"
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

func ParseOptionsAsSet(raw datatypes.JSON) (map[string]struct{}, error) {
	if len(raw) == 0 {
		return map[string]struct{}{}, nil
	}
	var arr []string
	if err := json.Unmarshal([]byte(raw), &arr); err != nil {
		return nil, err
	}
	set := make(map[string]struct{}, len(arr))
	for _, v := range arr {
		v = strings.TrimSpace(v)
		if v == "" {
			continue
		}
		set[v] = struct{}{}
	}
	return set, nil
}

func BuildVariantKeyAndValidate(attrs map[string]string, axes []types.AxisDef) (string, error) {
	parts := make([]string, 0, len(axes))

	for _, ax := range axes {
		val, ok := attrs[ax.Key]
		val = strings.TrimSpace(val)

		if ax.Required && (!ok || val == "") {
			return "", fmt.Errorf("%s: %s", "missing required axis", ax.Key)
		}

		if !ok || val == "" {
			return "", fmt.Errorf("%s: %s", "missing required axis", ax.Key)
		}

		if len(ax.Options) > 0 {
			if _, exists := ax.Options[val]; !exists {
				return "", fmt.Errorf("%s: %s=%s", "invalid attribute value", ax.Key, val)
			}
		}

		parts = append(parts, fmt.Sprintf("%s=%s", ax.Key, NormValue(val)))
	}

	return strings.Join(parts, "|"), nil
}

func RejectUnknownKeys(attrs map[string]string, axes []types.AxisDef) error {
	allowedKeys := make(map[string]struct{}, len(axes))
	for _, ax := range axes {
		allowedKeys[ax.Key] = struct{}{}
	}
	for k := range attrs {
		if _, ok := allowedKeys[k]; !ok {
			return fmt.Errorf("%s: %s", "invalid attribute key", k)
		}
	}
	return nil
}

func NormValue(s string) string {
	s = strings.TrimSpace(s)
	s = strings.ToLower(s)
	s = strings.Join(strings.Fields(s), " ") // collapse whitespace
	return s
}
