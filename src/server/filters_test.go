package server

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/jackc/pgtype"
	"github.com/labstack/echo/v4"
	"github.com/letsblockit/letsblockit/src/db"
	"github.com/letsblockit/letsblockit/src/filters"
	"github.com/letsblockit/letsblockit/src/pages"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	filterRepo *filters.Repository
	filter1    *filters.Filter
	filter2    *filters.Filter
	filter3    *filters.Filter

	filter2Defaults = map[string]any{
		"one":                    "default",
		"two":                    true,
		"three":                  []any{"a", "b"},
		"three---preset---dummy": false,
	}
	filter2DefaultOutput = "hello a default\nhello b default\n"
	filter2Custom        = map[string]any{
		"one":                    "blep",
		"two":                    true,
		"three":                  []string{"one", "two"},
		"three---preset---dummy": false,
	}
	filter2CustomOutput = "hello one blep\nhello two blep\n"
	filter2Preset       = map[string]any{
		"one":                    "blep",
		"two":                    true,
		"three":                  []string{"one"},
		"three---preset---dummy": true,
	}
	filter2PresetOutput = `hello one blep
!! filter2 with dummy preset
hello presetA blep
hello presetB blep
`
	filter2PresetTestOutput = `hello one blep:style(border: 2px dashed red !important)
!! filter2 with dummy preset
hello presetA blep:style(border: 2px dashed red !important)
hello presetB blep:style(border: 2px dashed red !important)
`
	filterTags = []string{"tag1", "tag2", "tag3"}
)

func (s *ServerTestSuite) TestListFilters_OK() {
	req := httptest.NewRequest(http.MethodGet, "/filters", nil)

	require.NoError(s.T(), s.server.upsertFilterParams(s.c, s.user, &filters.Instance{Filter: "filter1", TestMode: true}))
	require.NoError(s.T(), s.server.upsertFilterParams(s.c, s.user, &filters.Instance{Filter: "filter2"}))
	token := s.markListDownloaded()

	s.expectRender("list-filters", pages.ContextData{
		"filter_tags":       filterTags,
		"active_filters":    []*filters.Filter{filter1, filter2},
		"available_filters": []*filters.Filter{filter3},
		"testing_filters":   map[string]bool{"filter1": true},
		"list_downloaded":   true,
		"list_token":        token,
		"updated_filters":   map[string]bool{"filter2": true},
	})
	s.runRequest(req, assertOk)
}

func (s *ServerTestSuite) TestListFilters_ByTag() {
	req := httptest.NewRequest(http.MethodGet, "/filters/tag/tag2", nil)

	require.NoError(s.T(), s.server.upsertFilterParams(s.c, s.user, &filters.Instance{Filter: "filter1"}))
	list, err := s.store.GetListForUser(context.Background(), s.user)
	require.NoError(s.T(), err)

	s.expectRender("list-filters", pages.ContextData{
		"filter_tags":       filterTags,
		"tag_search":        "tag2",
		"active_filters":    []*filters.Filter{filter1},
		"available_filters": []*filters.Filter{filter2},
		"list_downloaded":   false,
		"list_token":        list.Token.String(),
	})
	s.runRequest(req, assertOk)
}

func (s *ServerTestSuite) TestViewFilter_Anonymous() {
	req := httptest.NewRequest(http.MethodGet, "/filters/filter2", nil)
	s.expectRender("view-filter", pages.ContextData{
		"filter":    filter2,
		"rendered":  filter2DefaultOutput,
		"params":    filter2Defaults,
		"test_mode": false,
	})
	s.runRequest(req, assertOk)
}

func (s *ServerTestSuite) TestViewFilter_NoInstance() {
	req := httptest.NewRequest(http.MethodGet, "/filters/filter2", nil)
	s.expectRender("view-filter", pages.ContextData{
		"filter":    filter2,
		"rendered":  filter2DefaultOutput,
		"params":    filter2Defaults,
		"test_mode": false,
	})
	s.runRequest(req, assertOk)
}

func (s *ServerTestSuite) TestViewFilter_HasInstance() {
	req := httptest.NewRequest(http.MethodGet, "/filters/filter2", nil)

	params := map[string]any{
		"two":   true,
		"three": []any{"one", "two"},
	}
	require.NoError(s.T(), s.server.upsertFilterParams(s.c, s.user, &filters.Instance{
		Filter: "filter2",
		Params: params,
	}))
	s.expectRender("view-filter", pages.ContextData{
		"filter":       filter2,
		"rendered":     "hello one \nhello two \n",
		"params":       params,
		"has_instance": true,
		"test_mode":    false,
		"new_params":   map[string]bool{"one": true, "three---preset---dummy": true},
	})
	s.runRequest(req, assertOk)
}

func (s *ServerTestSuite) TestViewFilter_HasTestInstance() {
	req := httptest.NewRequest(http.MethodGet, "/filters/filter2", nil)

	params := map[string]any{
		"two":   true,
		"three": []any{"one", "two"},
	}
	require.NoError(s.T(), s.server.upsertFilterParams(s.c, s.user, &filters.Instance{
		Filter:   "filter2",
		Params:   params,
		TestMode: true,
	}))
	s.expectRender("view-filter", pages.ContextData{
		"filter":       filter2,
		"rendered":     "hello one :style(border: 2px dashed red !important)\nhello two :style(border: 2px dashed red !important)\n",
		"params":       params,
		"has_instance": true,
		"test_mode":    true,
		"new_params":   map[string]bool{"one": true, "three---preset---dummy": true},
	})
	s.runRequest(req, assertOk)
}

func (s *ServerTestSuite) TestViewFilter_Preview() {
	f := buildFilter2CustomBody()
	f.Add(csrfLookup, s.csrf)
	req := httptest.NewRequest(http.MethodPost, "/filters/filter2", strings.NewReader(f.Encode()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)

	s.expectRender("view-filter", pages.ContextData{
		"filter":    filter2,
		"params":    filter2Custom,
		"rendered":  filter2CustomOutput,
		"test_mode": false,
	})
	s.runRequest(req, assertOk)
	s.requireInstanceCount("filter2", 0)
}

func (s *ServerTestSuite) TestViewFilter_Create() {
	f := buildFilter2CustomBody()
	f.Add(csrfLookup, s.csrf)
	f.Add("__save", "")
	req := httptest.NewRequest(http.MethodPost, "/filters/filter2", strings.NewReader(f.Encode()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)

	s.expectRender("view-filter", pages.ContextData{
		"filter":       filter2,
		"params":       filter2Custom,
		"rendered":     filter2CustomOutput,
		"has_instance": true,
		"saved_ok":     true,
		"test_mode":    false,
	})
	s.runRequest(req, assertOk)

	stored, err := s.store.GetInstanceForUserAndFilter(context.Background(), db.GetInstanceForUserAndFilterParams{
		UserID:     s.user,
		FilterName: "filter2",
	})
	require.NoError(s.T(), err)
	s.requireJSONEq(filter2Custom, stored.Params)
	s.requireInstanceCount("filter2", 1)
}

func (s *ServerTestSuite) TestViewFilter_CreateEmptyParams() {
	_, err := s.store.CreateListForUser(context.Background(), s.user)
	require.NoError(s.T(), err)

	f := buildFilter2CustomBody() // Add params that will be ignored: filter1 does not have any
	f.Add(csrfLookup, s.csrf)
	f.Add("__save", "")
	req := httptest.NewRequest(http.MethodPost, "/filters/filter1", strings.NewReader(f.Encode()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)

	s.expectRender("view-filter", pages.ContextData{
		"filter":       filter1,
		"params":       map[string]any{},
		"rendered":     "hello from one\n",
		"has_instance": true,
		"saved_ok":     true,
		"test_mode":    false,
	})
	s.runRequest(req, assertOk)

	stored, err := s.store.GetInstanceForUserAndFilter(context.Background(), db.GetInstanceForUserAndFilterParams{
		UserID:     s.user,
		FilterName: "filter1",
	})
	require.NoError(s.T(), err)
	require.Equal(s.T(), pgtype.Null, stored.Params.Status)
	s.requireInstanceCount("filter1", 1)
}

func (s *ServerTestSuite) TestViewFilter_Update() {
	require.NoError(s.T(), s.server.upsertFilterParams(s.c, s.user, &filters.Instance{Filter: "filter2"}))
	s.requireInstanceCount("filter2", 1)

	f := buildFilter2PresetBody()
	f.Add(csrfLookup, s.csrf)
	f.Add("__save", "")
	f.Add("__test_mode", "on")
	req := httptest.NewRequest(http.MethodPost, "/filters/filter2", strings.NewReader(f.Encode()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)

	s.expectRender("view-filter", pages.ContextData{
		"filter":       filter2,
		"params":       filter2Preset,
		"rendered":     filter2PresetTestOutput,
		"has_instance": true,
		"saved_ok":     true,
		"test_mode":    true,
	})
	s.runRequest(req, assertOk)

	stored, err := s.store.GetInstanceForUserAndFilter(context.Background(), db.GetInstanceForUserAndFilterParams{
		UserID:     s.user,
		FilterName: "filter2",
	})
	require.NoError(s.T(), err)
	s.True(stored.TestMode)
	s.requireJSONEq(filter2Preset, stored.Params)
	s.requireInstanceCount("filter2", 1)
}

func (s *ServerTestSuite) TestViewFilter_Disable() {
	require.NoError(s.T(), s.server.upsertFilterParams(s.c, s.user, &filters.Instance{Filter: "filter2"}))
	s.requireInstanceCount("filter2", 1)

	f := make(url.Values)
	f.Add(csrfLookup, s.csrf)
	f.Add("__disable", "")
	req := httptest.NewRequest(http.MethodPost, "/filters/filter2", strings.NewReader(f.Encode()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	s.expectP.RedirectToPage(gomock.Any(), "list-filters")
	s.runRequest(req, assertOk)
	s.requireInstanceCount("filter2", 0)
}

func (s *ServerTestSuite) TestViewFilter_MissingCSRF() {
	f := buildFilter2CustomBody()
	f.Add("__save", "")
	req := httptest.NewRequest(http.MethodPost, "/filters/filter2", strings.NewReader(f.Encode()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	s.runRequest(req, func(t *testing.T, recorder *httptest.ResponseRecorder) {
		assert.Equal(t, 400, recorder.Result().StatusCode)
	})
}

func (s *ServerTestSuite) TestViewFilter_NotFound() {
	req := httptest.NewRequest(http.MethodGet, "/filters/filter7", nil)
	s.runRequest(req, func(t *testing.T, rec *httptest.ResponseRecorder) {
		assert.Equal(t, http.StatusNotFound, rec.Code, rec.Body)
	})
}

func (s *ServerTestSuite) TestViewFilterRender_Defaults() {
	req := httptest.NewRequest(http.MethodPost, "/filters/filter2/render", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	s.expectRender("view-filter-render", pages.ContextData{
		"rendered": filter2DefaultOutput,
	})
	s.runRequest(req, assertOk)
}

func (s *ServerTestSuite) TestViewFilterRender_Custom() {
	req := httptest.NewRequest(http.MethodPost, "/filters/filter2/render", strings.NewReader(buildFilter2CustomBody().Encode()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)

	s.expectRender("view-filter-render", pages.ContextData{
		"rendered": filter2CustomOutput,
	})
	s.runRequest(req, assertOk)
}

func (s *ServerTestSuite) TestViewFilterRender_WithPreset() {
	req := httptest.NewRequest(http.MethodPost, "/filters/filter2/render", strings.NewReader(buildFilter2PresetBody().Encode()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)

	s.expectRender("view-filter-render", pages.ContextData{
		"rendered": filter2PresetOutput,
	})
	s.runRequest(req, assertOk)
}

func (s *ServerTestSuite) TestViewFilterRender_WithPresetAndTestMode() {
	f := buildFilter2PresetBody()
	f.Add("__test_mode", "on")
	req := httptest.NewRequest(http.MethodPost, "/filters/filter2/render", strings.NewReader(f.Encode()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)

	s.expectRender("view-filter-render", pages.ContextData{
		"rendered": filter2PresetTestOutput,
	})
	s.runRequest(req, assertOk)
}

func (s *ServerTestSuite) TestViewFilterRender_NotFound() {
	req := httptest.NewRequest(http.MethodPost, "/filters/filter7/render", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	s.runRequest(req, func(t *testing.T, rec *httptest.ResponseRecorder) {
		assert.Equal(t, http.StatusNotFound, rec.Code, rec.Body)
	})
}

func (s *ServerTestSuite) TestViewFilterRender_LoggedIn() {
	f := buildFilter2CustomBody()
	f.Add("__logged_in", "true")
	req := httptest.NewRequest(http.MethodPost, "/filters/filter2/render", strings.NewReader(f.Encode()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)

	s.expectRenderWithContext("view-filter-render", &pages.Context{
		NakedContent:    true,
		CurrentSection:  "filters",
		NavigationLinks: navigationLinks,
		UserLoggedIn:    true,
		Data: pages.ContextData{
			"rendered": filter2CustomOutput,
		},
	})
	s.runRequest(req, assertOk)
}

func buildFilter2CustomBody() url.Values {
	f := make(url.Values)
	f.Add("one", "blep")
	f.Add("two", "on")
	f.Add("three", "one")
	f.Add("three", "two")
	return f
}

func buildFilter2PresetBody() url.Values {
	f := make(url.Values)
	f.Add("one", "blep")
	f.Add("two", "on")
	f.Add("three", "one")
	f.Add("three---preset---dummy", "on")
	return f
}
