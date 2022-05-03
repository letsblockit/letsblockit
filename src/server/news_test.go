package server

import (
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/xvello/letsblockit/src/db"
	"github.com/xvello/letsblockit/src/news"
	"github.com/xvello/letsblockit/src/pages"
)

var exampleReleases = []*news.Release{{
	CreatedAt: fixedNow.Add(time.Hour),
}, {
	CreatedAt: fixedNow.Add(time.Hour),
}, {
	CreatedAt: fixedNow,
}, {
	CreatedAt: fixedNow.Add(-1 * time.Hour),
}}

func (s *ServerTestSuite) TestNews_Anonymous() {
	req := httptest.NewRequest(http.MethodGet, "/news", nil)
	s.releases = exampleReleases
	s.expectRender("news", pages.ContextData{
		"releases":    exampleReleases,
		"newReleases": make(map[string]bool),
	})
	s.runRequest(req, assertOk)
}

func (s *ServerTestSuite) TestNews_LoggedIn() {
	req := httptest.NewRequest(http.MethodGet, "/news", nil)
	req.AddCookie(verifiedCookie)
	s.releases = exampleReleases
	s.preferences = &db.UserPreference{
		UserID:     s.user,
		LatestNews: fixedNow,
	}
	s.expectUP.BumpLatestNews(gomock.Any(), s.user)
	s.expectRender("news", pages.ContextData{
		"releases": exampleReleases,
		"newReleases": map[string]bool{
			"0": true,
			"1": true,
		},
	})
	s.runRequest(req, assertOk)
}
