package helpers

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func GetPagination(c *gin.Context) (page, limit, offset int) {
	page = 1
	limit = 10

	if p := c.Query("page"); p != "" {
		fmt.Sscan(p, &page)
	}
	if l := c.Query("limit"); l != "" {
		fmt.Sscan(l, &limit)
	}

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	offset = (page - 1) * limit
	return
}
