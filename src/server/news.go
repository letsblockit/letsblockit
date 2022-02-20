package server

import (
	"github.com/labstack/echo/v4"
)

func (s *Server) newsHandler(c echo.Context) error {
	hc := s.buildPageContext(c, "Recent changes")

	releases, err := s.releases.GetReleases()
	if err != nil {
		return err
	}
	hc.Add("releases", releases)
	return s.pages.Render(c, "news", hc)
}
