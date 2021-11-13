package server

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xvello/letsblockit/src/filters"
	"github.com/xvello/letsblockit/src/pages"
)

func (s *ServerTestSuite) TestHomepage_Anonymous() {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	tList := []string{"one", "two"}
	fList := []*filters.Filter{{
		Name: "f1",
	}}
	s.expectF.GetTags().Return(tList)
	s.expectF.GetFilters().Return(fList)
	s.expectRender("list-filters", pages.ContextData{
		"filter_tags":       tList,
		"available_filters": fList,
	})
	s.runRequest(req, assertOk)
}

func (s *ServerTestSuite) TestAbout_Anonymous() {
	req := httptest.NewRequest(http.MethodGet, "/about", nil)
	s.expectRenderWithContext("about", &pages.Context{
		CurrentSection:  "about",
		NavigationLinks: navigationLinks,
		Title:           "About: Let’s block it!",
		UserLoggedIn:    false,
	})
	s.runRequest(req, assertOk)
}

func (s *ServerTestSuite) TestAbout_LoggedVerified() {
	req := httptest.NewRequest(http.MethodGet, "/about", nil)
	req.AddCookie(verifiedCookie)
	s.expectRenderWithContext("about", &pages.Context{
		CurrentSection:  "about",
		NavigationLinks: navigationLinks,
		Title:           "About: Let’s block it!",
		UserID:          s.user,
		UserLoggedIn:    true,
		UserVerified:    true,
	})
	s.runRequest(req, assertOk)
}

func (s *ServerTestSuite) TestAbout_LoggedNoVerified() {
	req := httptest.NewRequest(http.MethodGet, "/about", nil)
	req.AddCookie(unverifiedCookie)
	s.expectRenderWithContext("about", &pages.Context{
		CurrentSection:  "about",
		NavigationLinks: navigationLinks,
		Title:           "About: Let’s block it!",
		UserID:          s.user,
		UserLoggedIn:    true,
		UserVerified:    false,
	})
	s.runRequest(req, assertOk)
}

func (s *ServerTestSuite) TestAbout_KratosDown() {
	s.oryServer.Close() // Kratos is unresponsive, continue anonymous
	req := httptest.NewRequest(http.MethodGet, "/about", nil)
	req.AddCookie(verifiedCookie)
	s.expectRenderWithContext("about", &pages.Context{
		CurrentSection:  "about",
		NavigationLinks: navigationLinks,
		Title:           "About: Let’s block it!",
	})
	s.runRequest(req, assertOk)
}

func (s *ServerTestSuite) TestAbout_InvalidKratosResponse() {
	req := httptest.NewRequest(http.MethodGet, "/about", nil)
	req.AddCookie(&http.Cookie{
		Name:  "ory_session_verified",
		Value: "invalid ]]]]", // JSON parsing error -> continue anonymous
	})
	s.expectRenderWithContext("about", &pages.Context{
		CurrentSection:  "about",
		NavigationLinks: navigationLinks,
		Title:           "About: Let’s block it!",
	})
	s.runRequest(req, assertOk)
}

func TestServerDryRun(t *testing.T) {
	// Try to use the unix socket, fallback to TCP on localhost
	pgHost := "/var/run/postgresql"
	if _, err := os.Stat(pgHost); err != nil {
		pgHost = "localhost"
	}

	server := NewServer(&Options{
		DryRun:       true,
		Migrations:   true,
		Reload:       true,
		Statsd:       "localhost:8125",
		DatabaseName: "lbi_tests",
		DatabaseHost: pgHost,
	})
	assert.Equal(t, ErrDryRunFinished, server.Start())
}
