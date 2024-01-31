package util

import (
	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
	"net/http"
)

func TransformError(err error) error {
	switch errObj := err.(type) {
	case *pq.Error:
		if errObj.Code == "23505" {
			return echo.NewHTTPError(http.StatusConflict, "phone number already existed")
		}
	default:
		if err != nil && err.Error() == "sql: no rows in result set" {
			return nil
		}
	}

	return err
}
