package middleware

import (
	"context"
	"github.com/SawitProRecruitment/UserService/shared/jwt"
	"github.com/SawitProRecruitment/UserService/shared/util"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
	"strings"
)

func Auth(blacklistedUrl ...string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			if util.In(c.Request().Method+":"+c.Request().URL.Path, blacklistedUrl...) {
				return next(c)
			}

			authHeader := c.Request().Header.Get("Authorization")
			if !strings.HasPrefix(authHeader, "Bearer ") {
				return echo.NewHTTPError(http.StatusBadRequest, "invalid header type, should be bearer")
			}

			claims, err := jwt.GetSigner().ParseWithClaims(strings.TrimPrefix(authHeader, "Bearer "))
			if err != nil {
				return echo.NewHTTPError(http.StatusBadRequest, err.Error())
			}

			ctx := c.Request().Context()
			ctx = context.WithValue(ctx, "Claims", claims)
			ctx = context.WithValue(ctx, "UserID", claims.UserId)
			c.SetRequest(c.Request().WithContext(ctx))

			return next(c)
		}
	}
}

func Logger() echo.MiddlewareFunc    { return middleware.Logger() }
func RequestID() echo.MiddlewareFunc { return middleware.RequestID() }
