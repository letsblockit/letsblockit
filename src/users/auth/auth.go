package auth

import "github.com/labstack/echo/v4"

const (
	hasAccountCookieName  = "has_account"
	hasAccountCookieValue = "true"
	userContextKey        = "_user"
	hasAuthContextKey     = "_has_auth"
	userActionRouteName   = "user-action"
)

type Backend interface {
	BuildMiddleware() echo.MiddlewareFunc
	RegisterRoutes(group *echo.Group)
}

func setUserId(c echo.Context, id string) {
	c.Set(userContextKey, id)
}

func GetUserId(c echo.Context) string {
	if u, ok := c.Get(userContextKey).(string); ok {
		return u
	}
	return ""
}

func HasAccount(c echo.Context) bool {
	_, err := c.Cookie(hasAccountCookieName)
	return err == nil
}

func HasAuth(c echo.Context) bool {
	return c.Get(hasAuthContextKey) != nil
}
