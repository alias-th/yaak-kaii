package utils

import (
	"fmt"
	"strings"
)

func GenerateSKU(category, color, size string) string {
	c := strings.ToUpper(category[:3])
	cl := strings.ToUpper(color[:2])
	sz := strings.ToUpper(size)

	return fmt.Sprintf("%s-%s-%s", c, cl, sz)
}
