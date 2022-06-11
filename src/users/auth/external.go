package auth

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// External is used to run the server being an authenticating proxy.
// It expects an HTTP header to be set by the proxy, and uses its value as the unique user ID.
type External struct {
	headerName string
}

func NewExternal(headerName string) *External {
	return &External{headerName: headerName}
}

func (e *External) BuildMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if id := c.Request().Header.Get(e.headerName); id != "" {
				setUserId(c, id)
				return next(c)
			}
			return c.String(http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
		}
	}
}

func (e *External) RegisterRoutes(_ *echo.Group) {}
