package util

import (
	"context"
	"github.com/labstack/echo/v4"
	"net/http"
)

func GetUserIDFromContext(ctx context.Context) (int64, error) {
	if ctx.Value("UserID") == nil {
		return 0, echo.NewHTTPError(http.StatusUnauthorized, "user are not logged in")
	}

	userId, ok := ctx.Value("UserID").(int64)
	if !ok {
		return 0, echo.NewHTTPError(http.StatusUnauthorized, "user are not logged in")
	}

	return userId, nil
}
