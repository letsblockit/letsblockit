package server

import (
	"net/http"
	"net/http/httptest"

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
	req.AddCookie(verifiedCookie)
	s.expectRenderWithContext("about", &pages.Context{
		CurrentSection:  "about",
		NavigationLinks: navigationLinks,
		Title:           "About: Let’s block it!",
		UserID:          s.user,
		UserVerified:    true,
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
