package server

import (
	"net/http"
	"net/http/httptest"

	"github.com/xvello/letsblockit/src/filters"
	"github.com/xvello/letsblockit/src/pages"
)

var filter1 = &filters.Filter{
	Name: "filter1",
	Tags: []string{"tag1", "tag2"},
}

var filter2 = &filters.Filter{
	Name: "filter2",
	Params: []filters.FilterParam{{
		Name:    "one",
		Type:    filters.StringParam,
		Default: "default",
	}, {
		Name:    "two",
		Type:    filters.BooleanParam,
		Default: true,
	}},
	Tags: []string{"tag2", "tag3"},
}
var filter3 = &filters.Filter{
	Name: "filter3",
	Tags: []string{"tag3"},
}

func (s *ServerTestSuite) TestListFilters_OK() {
	req := httptest.NewRequest(http.MethodGet, "/filters", nil)
	s.login(true)

	tList := []string{"tag1", "tag2", "tag3"}
	s.expectF.GetTags().Return(tList)
	s.expectF.GetFilters().Return([]*filters.Filter{filter1, filter2, filter3})
	s.expectS.GetActiveFilterNames(s.user).Return(map[string]bool{"filter2": true})

	s.expectRender("list-filters", pages.ContextData{
		"filter_tags":       tList,
		"active_filters":    []*filters.Filter{filter2},
		"available_filters": []*filters.Filter{filter1, filter3},
	})
	s.runRequest(req, assertOk)
}

func (s *ServerTestSuite) TestListFilters_ByTag() {
	req := httptest.NewRequest(http.MethodGet, "/filters/tag/tag2", nil)
	s.login(true)

	tList := []string{"tag1", "tag2", "tag3"}
	s.expectF.GetTags().Return(tList)
	s.expectF.GetFilters().Return([]*filters.Filter{filter1, filter2, filter3})
	s.expectS.GetActiveFilterNames(s.user).Return(map[string]bool{"filter2": true})

	s.expectRender("list-filters", pages.ContextData{
		"filter_tags":       tList,
		"tag_search":        "tag2",
		"active_filters":    []*filters.Filter{filter2},
		"available_filters": []*filters.Filter{filter1},
	})
	s.runRequest(req, assertOk)
}
