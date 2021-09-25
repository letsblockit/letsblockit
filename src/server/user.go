package server

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (s *Server) userLogin(c echo.Context) error {
	if getUser(c) != nil {
		return c.Redirect(http.StatusFound, "/user/account")
	}
	hc := s.buildHandlebarsContext(c, "Login")
	return s.pages.render(c, "user-login", hc)
}

func (s *Server) userAccount(c echo.Context) error {
	user := getUser(c)
	if user == nil {
		// Cannot find user session, redirect to login
		return c.Redirect(http.StatusFound, "/user/login")
	}

	hc := s.buildHandlebarsContext(c, "My account")
	if user.IsVerified() {
		hc["verified"] = true

	}
	return s.pages.render(c, "user-account", hc)
}
