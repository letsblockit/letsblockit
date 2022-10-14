package auth

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// Proxy is used to run the server being an authenticating proxy.
// It expects an HTTP header to be set by the proxy, and uses its value as the unique user ID.
type Proxy struct {
	headerName string
}

func NewProxy(headerName string) *Proxy {
	return &Proxy{headerName: headerName}
}

func (e *Proxy) BuildMiddleware() echo.MiddlewareFunc {
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

func (e *Proxy) RegisterRoutes(_ EchoRouter) {}
