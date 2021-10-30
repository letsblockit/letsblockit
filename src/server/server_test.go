package server

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
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
	})
	s.runRequest(req, assertOk)
}

func (s *ServerTestSuite) TestAbout_LoggedVerified() {
	req := httptest.NewRequest(http.MethodGet, "/about", nil)
	s.login(true)
	s.expectRenderWithContext("about", &pages.Context{
		CurrentSection:  "about",
		NavigationLinks: navigationLinks,
		Title:           "About: Let’s block it!",
		UserID:          s.user,
		UserVerified:    true,
	})
	s.runRequest(req, assertOk)
}

// Checks that the server can instantiate all its components.
// This includes parsing all pages and filters and creating a sqlite file.
func TestServerDryRun(t *testing.T) {
	dir, err := ioutil.TempDir("", "lbi")
	require.NoError(t, err)
	defer os.RemoveAll(dir)
	server := NewServer(&Options{
		DryRun:     true,
		Migrations: true,
		DataFolder: dir,
	})
	require.Equal(t, ErrDryRunFinished, server.Start())
}
