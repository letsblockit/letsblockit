package server

import (
	"github.com/labstack/echo/v4"
)

func (s *Server) landingPageHandler(c echo.Context) error {
	hc := s.buildPageContext(c, "Let's Block It!")
	if hc.UserLoggedIn {
		c.Response().Header().Set("HX-Push", "/filters")
		return s.listFilters(c)
	}
	return s.pages.Render(c, "landing", hc)
}
