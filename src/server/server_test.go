package server

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/letsblockit/letsblockit/src/db"
	"github.com/letsblockit/letsblockit/src/news"
	"github.com/letsblockit/letsblockit/src/pages"
	"github.com/stretchr/testify/assert"
)

func TestServerDryRun(t *testing.T) {
	assert.Equal(t, ErrDryRunFinished, NewServer(&Options{
		AuthMethod:    "kratos",
		AuthKratosUrl: "http://localhost:4000/.ory",
		DatabaseUrl:   db.GetTestDatabaseURL(),
		StatsdTarget:  "localhost:8125",
		DryRun:        true,
		HotReload:     true,
	}).Start())
}

func (s *ServerTestSuite) TestHomepage_Anonymous() {
	s.user = ""
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	s.expectRender("landing", nil)
	s.runRequest(req, assertOk)
}

func (s *ServerTestSuite) TestAbout_Anonymous() {
	s.user = ""
	req := httptest.NewRequest(http.MethodGet, "/about", nil)
	s.expectRenderWithSidebarAndContext("help-about", "help-sidebar", &pages.Context{
		CurrentSection:  "help/about",
		NavigationLinks: navigationLinks,
		Title:           "About this project",
		Data: pages.ContextData{
			"page":          helpMenu[1].Pages[0],
			"menu_sections": helpMenu,
		},
		CSRFToken:    s.csrf,
		UserLoggedIn: false,
	})
	s.runRequest(req, assertOk)
}

func (s *ServerTestSuite) TestAbout_Logged() {
	pref, _ := s.server.preferences.Get(s.c, s.user)
	req := httptest.NewRequest(http.MethodGet, "/about", nil)
	s.expectRenderWithSidebarAndContext("help-about", "help-sidebar", &pages.Context{
		CurrentSection:  "help/about",
		NavigationLinks: navigationLinks,
		Title:           "About this project",
		Data: pages.ContextData{
			"page":          helpMenu[1].Pages[0],
			"menu_sections": helpMenu,
		},
		CSRFToken:    s.csrf,
		UserID:       s.user,
		UserLoggedIn: true,
		Preferences:  pref,
	})
	s.runRequest(req, assertOk)
}

func (s *ServerTestSuite) TestAbout_HasNews() {
	pref, _ := s.server.preferences.Get(s.c, s.user)
	s.releases = append(s.releases, &news.Release{
		CreatedAt: fixedNow.Add(time.Hour),
	})
	req := httptest.NewRequest(http.MethodGet, "/about", nil)
	s.expectRenderWithSidebarAndContext("help-about", "help-sidebar", &pages.Context{
		CurrentSection:  "help/about",
		NavigationLinks: navigationLinks,
		Title:           "About this project",
		Data: pages.ContextData{
			"page":          helpMenu[1].Pages[0],
			"menu_sections": helpMenu,
		},
		CSRFToken:    s.csrf,
		UserID:       s.user,
		UserLoggedIn: true,
		Preferences:  pref,
		HasNews:      true,
	})
	s.runRequest(req, assertOk)
}

func (s *ServerTestSuite) TestAbout_DarkMode() {
	s.NoError(s.server.preferences.UpdatePreferences(s.c, db.UpdateUserPreferencesParams{
		UserID:       s.user,
		ColorMode:    db.ColorModeDark,
		BetaFeatures: false,
	}))
	pref, _ := s.server.preferences.Get(s.c, s.user)
	req := httptest.NewRequest(http.MethodGet, "/about", nil)
	s.expectRenderWithSidebarAndContext("help-about", "help-sidebar", &pages.Context{
		CurrentSection:  "help/about",
		NavigationLinks: navigationLinks,
		Title:           "About this project",
		Data: pages.ContextData{
			"page":          helpMenu[1].Pages[0],
			"menu_sections": helpMenu,
		},
		CSRFToken:    s.csrf,
		UserID:       s.user,
		UserLoggedIn: true,
		Preferences:  pref,
	})
	s.runRequest(req, assertOk)
}

func (s *ServerTestSuite) TestAbout_BannedUser() {
	s.setUserBanned()
	req := httptest.NewRequest(http.MethodGet, "/about", nil)
	s.runRequest(req, func(t *testing.T, recorder *httptest.ResponseRecorder) {
		assert.Equal(t, 403, recorder.Result().StatusCode)
	})
}

func (s *ServerTestSuite) TestShouldReload_OK() {
	req := httptest.NewRequest(http.MethodGet, "http://localhost:4000/should-reload", nil)
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	req = req.WithContext(ctx)
	rec := httptest.NewRecorder()
	s.server.echo.ServeHTTP(rec, req)

	s.True(rec.Flushed)
	s.Equal("retry:1000\n", rec.Body.String())
}

func (s *ServerTestSuite) TestShouldReload_BadHost() {
	req := httptest.NewRequest(http.MethodGet, "http://unexpected/should-reload", nil)
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	req = req.WithContext(ctx)
	rec := httptest.NewRecorder()
	s.server.echo.ServeHTTP(rec, req)

	s.Equal(404, rec.Code, rec.Body)
}
