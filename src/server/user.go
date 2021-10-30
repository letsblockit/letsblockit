package server

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (s *Server) userLogin(c echo.Context) error {
	if getUser(c) != nil {
		return s.redirect(c, "user-account")
	}
	hc := s.buildHandlebarsContext(c, "Login")
	return s.pages.Render(c, "user-login", hc)
}

func (s *Server) userLogout(c echo.Context) error {
	if getUser(c) == nil {
		return s.redirect(c, "user-login")
	}
	logout, err := getLogoutUrl(s.options.OryProject, c)
	if err != nil {
		return err
	}
	return c.Redirect(http.StatusFound, logout)
}

func (s *Server) userAccount(c echo.Context) error {
	user := getUser(c)
	if user == nil {
		// Cannot find user session, redirect to login
		return s.redirect(c, "user-login")
	}

	var err error
	hc := s.buildHandlebarsContext(c, "My account")
	if user.IsVerified() {
		hc["filter_count"], err = s.store.CountFilters(user.Id())
		if err != nil {
			return err
		}
		hc["filter_list"], err = s.store.GetOrCreateFilterList(user.Id())
		if err != nil {
			return err
		}
	}
	return s.pages.Render(c, "user-account", hc)
}
