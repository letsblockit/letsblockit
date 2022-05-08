package server

import (
	"encoding/xml"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/xvello/letsblockit/src/db"
	"github.com/xvello/letsblockit/src/news"
	"github.com/xvello/letsblockit/src/pages"
	"golang.org/x/tools/blog/atom"
)

var exampleReleases = []*news.Release{
	{CreatedAt: fixedNow.Add(time.Hour)},
	{CreatedAt: fixedNow.Add(time.Hour)},
	{CreatedAt: fixedNow},
	{CreatedAt: fixedNow.Add(-1 * time.Hour)},
}

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
		NewsCursor: fixedNow,
	}
	s.expectUP.UpdateNewsCursor(gomock.Any(), s.user, exampleReleases[0].CreatedAt)
	s.expectRender("news", pages.ContextData{
		"releases": exampleReleases,
		"newReleases": map[string]bool{
			"0": true,
			"1": true,
		},
	})
	s.runRequest(req, assertOk)
}

func (s *ServerTestSuite) TestNews_NoNews() {
	req := httptest.NewRequest(http.MethodGet, "/news", nil)
	req.AddCookie(verifiedCookie)
	s.releases = exampleReleases
	s.preferences = &db.UserPreference{
		UserID:     s.user,
		NewsCursor: exampleReleases[0].CreatedAt,
	}
	s.expectRender("news", pages.ContextData{
		"releases":    exampleReleases,
		"newReleases": map[string]bool{},
	})
	s.runRequest(req, assertOk)
}

func (s *ServerTestSuite) TestNewsAtom() {
	req := httptest.NewRequest(http.MethodGet, "/news.atom", nil)
	s.releases = exampleReleases
	s.runRequest(req, func(t *testing.T, rec *httptest.ResponseRecorder) {
		assert.Equal(t, http.StatusOK, rec.Code, rec.Body)
		var feed atom.Feed
		assert.NoError(t, xml.Unmarshal(rec.Body.Bytes(), &feed))
		assert.EqualValues(t, atom.Feed{}, feed)
	})
}
