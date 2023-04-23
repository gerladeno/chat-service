package middlewares

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func BlockMethods(methods ...string) echo.MiddlewareFunc {
	forbiddenMethods := make(map[string]struct{}, len(methods))
	for i := range methods {
		forbiddenMethods[methods[i]] = struct{}{}
	}
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if _, ok := forbiddenMethods[c.Request().Method]; ok {
				return c.JSON(http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
			}
			return next(c)
		}
	}
}
