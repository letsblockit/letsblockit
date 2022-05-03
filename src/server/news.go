package server

import (
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/xvello/letsblockit/src/news"
)

type Release struct {
	*news.Release
	Fresh bool
}

func (s *Server) newsHandler(c echo.Context) error {
	releases, err := s.releases.GetReleases()
	if err != nil {
		return err
	}

	hc := s.buildPageContext(c, "Recent changes")
	if hc.UserLoggedIn {
		err := s.preferences.BumpLatestNews(c, hc.UserID)
		if err != nil {
			c.Logger().Warnf("failed to update latest news for %s: %s", hc.UserID, err)
		}
	}

	newReleases := make(map[string]bool) // handlebars lookup only supports string keys
	if hc.HasNews {
		// Shut down the menubar notification for this page
		hc.HasNews = false

		// Compute new releases to highlight
		if hc.Preferences != nil {
			for i, r := range releases {
				if r.CreatedAt.After(hc.Preferences.LatestNews) {
					newReleases[strconv.Itoa(i)] = true
				} else {
					break
				}
			}
		}
	}

	hc.Add("releases", releases)
	hc.Add("newReleases", newReleases)
	return s.pages.Render(c, "news", hc)
}
