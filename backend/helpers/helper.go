package helpers

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func UserIDFromContext(ctx context.Context) (int64, error) {
	id, ok := ctx.Value("user_id").(int64)
	if !ok {
		return 0, errors.New("unauthorized")
	}
	return id, nil
}

func ParseIntQuery(r *http.Request, key string, defaultVal int) int {
	valStr := r.URL.Query().Get(key)
	if valStr == "" {
		return defaultVal
	}

	val, err := strconv.Atoi(valStr)
	if err != nil || val < 0 {
		return defaultVal
	}

	return val
}

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
