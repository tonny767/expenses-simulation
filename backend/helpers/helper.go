package helpers

import (
	"context"
	"errors"
	"net/http"
	"strconv"
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
