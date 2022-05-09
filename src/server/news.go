package server

import (
	"fmt"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/xvello/letsblockit/src/news"
	"golang.org/x/tools/blog/atom"
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

	hc := s.buildPageContext(c, "Release notes")

	newReleases := make(map[string]bool) // handlebars lookup only supports string keys
	if hc.HasNews {
		if hc.UserLoggedIn && len(releases) > 0 {
			err := s.preferences.UpdateNewsCursor(c, hc.UserID, releases[0].CreatedAt)
			if err != nil {
				c.Logger().Warnf("failed to update latest news for %s: %s", hc.UserID, err)
			}
		}

		// Shut down the menubar notification for this page
		hc.HasNews = false

		// Compute new releases to highlight
		if hc.Preferences != nil {
			for i, r := range releases {
				if r.CreatedAt.After(hc.Preferences.NewsCursor) {
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

func (s *Server) newsAtomHandler(c echo.Context) error {
	releases, err := s.releases.GetReleases()
	if err != nil {
		return err
	}

	feed := atom.Feed{
		Title: "Release notes from letsblock.it",
		ID:    s.absoluteReverse(c, "news-atom"),
		Link: []atom.Link{{
			Rel:  "self",
			Href: s.absoluteReverse(c, "news-atom"),
			Type: "application/atom+xml",
		}, {
			Rel:      "alternate",
			Href:     s.absoluteReverse(c, "news"),
			Type:     "text/html",
			HrefLang: "en",
		}},
		Author: &atom.Person{
			Name: "Let's Block It contributors",
			URI:  "https://github.com/xvello/letsblockit",
		},
	}

	var latestUpdate time.Time
	for _, r := range releases {
		feed.Entry = append(feed.Entry, &atom.Entry{
			Title: fmt.Sprintf("Let's Block It: %s update", r.TagName),
			ID:    r.GithubUrl,
			Link: []atom.Link{{
				Rel:  "alternate",
				Href: r.GithubUrl,
				Type: "text/html",
			}},
			Published: atom.Time(r.CreatedAt),
			Updated:   atom.Time(r.PublishedAt),
			Author:    nil,
			Summary:   nil,
			Content: &atom.Text{
				Type: "html",
				Body: r.Description,
			},
		})
		if r.PublishedAt.After(latestUpdate) {
			latestUpdate = r.PublishedAt
		}
	}

	feed.Updated = atom.Time(latestUpdate)
	c.Response().Header().Set(echo.HeaderContentType, "application/atom+xml")
	return c.XMLPretty(200, &feed, "\t")
}
