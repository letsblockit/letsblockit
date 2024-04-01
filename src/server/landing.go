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
	if s.options.Sunset {
		return s.pages.Render(c, "sunset", hc)
	}
	return s.pages.Render(c, "landing", hc)
}

func (s *Server) sunsetHandler(c echo.Context) error {
	hc := s.buildPageContext(c, "Project shutdown information")
	return s.pages.Render(c, "sunset", hc)
}
