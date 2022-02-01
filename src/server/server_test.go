package server

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/xvello/letsblockit/src/pages"
)

func (s *ServerTestSuite) TestHomepage_Anonymous() {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	s.expectRender("landing", nil)
	s.runRequest(req, assertOk)
}

func (s *ServerTestSuite) TestAbout_Anonymous() {
	req := httptest.NewRequest(http.MethodGet, "/about", nil)
	s.expectRenderWithSidebarAndContext("help-about", "help-sidebar", &pages.Context{
		CurrentSection:  "help/about",
		NavigationLinks: navigationLinks,
		Title:           "About this project",
		Data: pages.ContextData{
			"page":          helpMenu[1].Pages[0],
			"menu_sections": helpMenu,
		},
		UserLoggedIn: false,
	})
	s.runRequest(req, assertOk)
}

func (s *ServerTestSuite) TestAbout_Logged() {
	req := httptest.NewRequest(http.MethodGet, "/about", nil)
	req.AddCookie(verifiedCookie)
	s.expectRenderWithSidebarAndContext("help-about", "help-sidebar", &pages.Context{
		CurrentSection:  "help/about",
		NavigationLinks: navigationLinks,
		Title:           "About this project",
		Data: pages.ContextData{
			"page":          helpMenu[1].Pages[0],
			"menu_sections": helpMenu,
		},
		UserID:       s.user,
		UserLoggedIn: true,
	})
	s.runRequest(req, assertOk)
}

func (s *ServerTestSuite) TestAbout_BannedUser() {
	s.setUserBanned()
	req := httptest.NewRequest(http.MethodGet, "/about", nil)
	req.AddCookie(verifiedCookie)
	s.runRequest(req, func(t *testing.T, recorder *httptest.ResponseRecorder) {
		assert.Equal(t, 403, recorder.Result().StatusCode)
	})
}

func (s *ServerTestSuite) TestAbout_KratosDown() {
	s.kratosServer.Close() // Kratos is unresponsive, continue anonymous
	req := httptest.NewRequest(http.MethodGet, "/about", nil)
	req.AddCookie(verifiedCookie)
	s.expectRenderWithSidebarAndContext("help-about", "help-sidebar", &pages.Context{
		CurrentSection:  "help/about",
		NavigationLinks: navigationLinks,
		Title:           "About this project",
		Data: pages.ContextData{
			"page":          helpMenu[1].Pages[0],
			"menu_sections": helpMenu,
		},
	})
	s.runRequest(req, assertOk)
}

func (s *ServerTestSuite) TestAbout_InvalidKratosResponse() {
	req := httptest.NewRequest(http.MethodGet, "/about", nil)
	req.AddCookie(&http.Cookie{
		Name:  "ory_session_verified",
		Value: "invalid ]]]]", // JSON parsing error -> continue anonymous
	})
	s.expectRenderWithSidebarAndContext("help-about", "help-sidebar", &pages.Context{
		CurrentSection:  "help/about",
		NavigationLinks: navigationLinks,
		Title:           "About this project",
		Data: pages.ContextData{
			"page":          helpMenu[1].Pages[0],
			"menu_sections": helpMenu,
		},
	})
	s.runRequest(req, assertOk)
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
