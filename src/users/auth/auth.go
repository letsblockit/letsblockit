package auth

import "github.com/labstack/echo/v4"

const (
	hasAccountCookieName  = "has_account"
	hasAccountCookieValue = "true"
	userContextKey        = "_user"
	hasAuthContextKey     = "_has_auth"
	userActionRouteName   = "user-action"
)

type User interface {
	Id() string
}

type Backend interface {
	BuildMiddleware() echo.MiddlewareFunc
	RegisterRoutes(group *echo.Group)
}

func GetUserId(c echo.Context) string {
	if u, ok := c.Get(userContextKey).(User); ok {
		return u.Id()
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
