package server

import (
	"net/http"
	"net/http/httptest"

	"github.com/letsblockit/letsblockit/src/filters"
	"github.com/letsblockit/letsblockit/src/pages"
)

func (s *ServerTestSuite) TestLanding_Anonymous() {
	s.user = ""
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	s.expectRender("landing", nil)
	s.runRequest(req, assertOk)
}

func (s *ServerTestSuite) TestLanding_LoggedIn() {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	s.expectRender("list-filters", pages.ContextData{
		"filter_tags":       filterTags,
		"available_filters": []*filters.Filter{filter1, filter2, filter3},
	})
	s.runRequest(req, assertOk)
}
