package server

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/jackc/pgtype"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/xvello/letsblockit/src/db"
	"github.com/xvello/letsblockit/src/filters"
	"github.com/xvello/letsblockit/src/pages"
)

var filter1 = &filters.Filter{
	Name: "filter1",
	Tags: []string{"tag1", "tag2"},
}

var filter2 = &filters.Filter{
	Name:  "filter2",
	Title: "TITLE2",
	Params: []filters.FilterParam{{
		Name:    "one",
		Type:    filters.StringParam,
		Default: "default",
	}, {
		Name:    "two",
		Type:    filters.BooleanParam,
		Default: true,
	}, {
		Name:    "three",
		Type:    filters.StringListParam,
		Default: []string{"a", "b"},
	}},
	Tags: []string{"tag2", "tag3"},
}

var filter2Defaults = map[string]interface{}{
	"one":   "default",
	"two":   true,
	"three": []string{"a", "b"},
}

var filter3 = &filters.Filter{
	Name: "filter3",
	Tags: []string{"tag3"},
}

func (s *ServerTestSuite) TestListFilters_OK() {
	req := httptest.NewRequest(http.MethodGet, "/filters", nil)
	req.AddCookie(verifiedCookie)

	tList := []string{"tag1", "tag2", "tag3"}
	s.expectF.GetTags().Return(tList)
	s.expectF.GetFilters().Return([]*filters.Filter{filter1, filter2, filter3})
	s.expectQ.GetActiveFiltersForUser(gomock.Any(), s.user).Return([]string{"filter2"}, nil)
	s.expectQ.HasUserDownloadedList(gomock.Any(), s.user).Return(true, nil)

	s.expectRender("list-filters", pages.ContextData{
		"filter_tags":       tList,
		"active_filters":    []*filters.Filter{filter2},
		"available_filters": []*filters.Filter{filter1, filter3},
		"list_downloaded":   true,
	})
	s.runRequest(req, assertOk)
}

func (s *ServerTestSuite) TestListFilters_ByTag() {
	req := httptest.NewRequest(http.MethodGet, "/filters/tag/tag2", nil)
	req.AddCookie(verifiedCookie)

	tList := []string{"tag1", "tag2", "tag3"}
	s.expectF.GetTags().Return(tList)
	s.expectF.GetFilters().Return([]*filters.Filter{filter1, filter2, filter3})
	s.expectQ.GetActiveFiltersForUser(gomock.Any(), s.user).Return([]string{"filter2"}, nil)
	s.expectQ.HasUserDownloadedList(gomock.Any(), s.user).Return(false, nil)

	s.expectRender("list-filters", pages.ContextData{
		"filter_tags":       tList,
		"tag_search":        "tag2",
		"active_filters":    []*filters.Filter{filter2},
		"available_filters": []*filters.Filter{filter1},
		"list_downloaded":   false,
	})
	s.runRequest(req, assertOk)
}

func (s *ServerTestSuite) TestViewFilter_Anonymous() {
	req := httptest.NewRequest(http.MethodGet, "/filters/filter2", nil)
	s.expectF.GetFilter("filter2").Return(filter2, nil)
	s.expectRenderFilter("filter2", filter2Defaults, "output")
	s.expectRender("view-filter", pages.ContextData{
		"filter":   filter2,
		"rendered": "output",
		"params":   filter2Defaults,
	})
	s.runRequest(req, assertOk)
}

func (s *ServerTestSuite) TestViewFilter_NoInstance() {
	req := httptest.NewRequest(http.MethodGet, "/filters/filter2", nil)
	req.AddCookie(verifiedCookie)
	s.expectF.GetFilter("filter2").Return(filter2, nil)
	s.expectQ.GetInstanceForUserAndFilter(gomock.Any(), db.GetInstanceForUserAndFilterParams{
		UserID:     s.user,
		FilterName: "filter2",
	}).Return(pgtype.JSONB{}, db.NotFound)
	s.expectRenderFilter("filter2", filter2Defaults, "output")
	s.expectRender("view-filter", pages.ContextData{
		"filter":   filter2,
		"rendered": "output",
		"params":   filter2Defaults,
	})
	s.runRequest(req, assertOk)
}

func (s *ServerTestSuite) TestViewFilter_HasInstance() {
	req := httptest.NewRequest(http.MethodGet, "/filters/filter2", nil)
	req.AddCookie(verifiedCookie)
	s.expectF.GetFilter("filter2").Return(filter2, nil)

	params := map[string]interface{}{"one": "1", "two": false}
	paramsB := pgtype.JSONB{}
	s.NoError(paramsB.Set(&params))
	s.expectQ.GetInstanceForUserAndFilter(gomock.Any(), db.GetInstanceForUserAndFilterParams{
		UserID:     s.user,
		FilterName: "filter2",
	}).Return(paramsB, nil)
	s.expectRenderFilter("filter2", params, "output")
	s.expectRender("view-filter", pages.ContextData{
		"filter":       filter2,
		"rendered":     "output",
		"params":       params,
		"has_instance": true,
	})
	s.runRequest(req, assertOk)
}

func (s *ServerTestSuite) TestViewFilter_Preview() {
	req := httptest.NewRequest(http.MethodPost, "/filters/filter2", strings.NewReader(buildFilter2FormBody().Encode()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	req.AddCookie(verifiedCookie)
	s.expectF.GetFilter("filter2").Return(filter2, nil)
	s.expectQ.GetInstanceForUserAndFilter(gomock.Any(), db.GetInstanceForUserAndFilterParams{
		UserID:     s.user,
		FilterName: "filter2",
	}).Return(pgtype.JSONB{}, db.NotFound)
	params := map[string]interface{}{
		"one":   "1",
		"two":   false,
		"three": []string{"option1", "option2"},
	}

	s.expectRenderFilter("filter2", params, "output")
	s.expectRender("view-filter", pages.ContextData{
		"filter":   filter2,
		"params":   params,
		"rendered": "output",
	})
	s.runRequest(req, assertOk)
}

func (s *ServerTestSuite) TestViewFilter_Create() {
	f := buildFilter2FormBody()
	f.Add("__save", "")
	req := httptest.NewRequest(http.MethodPost, "/filters/filter2", strings.NewReader(f.Encode()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	req.AddCookie(verifiedCookie)
	s.expectF.GetFilter("filter2").Return(filter2, nil)

	params := map[string]interface{}{
		"one":   "1",
		"two":   false,
		"three": []string{"option1", "option2"},
	}
	paramsB := pgtype.JSONB{}
	s.NoError(paramsB.Set(&params))
	s.expectQ.CountListsForUser(gomock.Any(), s.user).Return(int64(1), nil)
	s.expectQ.CountInstanceForUserAndFilter(gomock.Any(), db.CountInstanceForUserAndFilterParams{
		UserID:     s.user,
		FilterName: "filter2",
	}).Return(int64(0), nil)
	s.expectQ.CreateInstanceForUserAndFilter(gomock.Any(), db.CreateInstanceForUserAndFilterParams{
		UserID:     s.user,
		FilterName: "filter2",
		Params:     paramsB,
	}).Return(nil)
	s.expectRenderFilter("filter2", params, "output")
	s.expectRender("view-filter", pages.ContextData{
		"filter":       filter2,
		"params":       params,
		"rendered":     "output",
		"has_instance": true,
		"saved_ok":     true,
	})
	s.runRequest(req, assertOk)
}

func (s *ServerTestSuite) TestViewFilter_CreateEmptyParams() {
	f := buildFilter2FormBody() // Add params that will be ignored: filter1 does not have any
	f.Add("__save", "")
	req := httptest.NewRequest(http.MethodPost, "/filters/filter1", strings.NewReader(f.Encode()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	req.AddCookie(verifiedCookie)
	s.expectF.GetFilter("filter1").Return(filter1, nil)

	paramsB := pgtype.JSONB{}
	s.NoError(paramsB.Set(nil))
	s.expectQ.CountInstanceForUserAndFilter(gomock.Any(), db.CountInstanceForUserAndFilterParams{
		UserID:     s.user,
		FilterName: "filter1",
	}).Return(int64(0), nil)
	s.expectQ.CreateInstanceForUserAndFilter(gomock.Any(), db.CreateInstanceForUserAndFilterParams{
		UserID:     s.user,
		FilterName: "filter1",
		Params:     paramsB,
	}).Return(nil)
	s.expectQ.CountListsForUser(gomock.Any(), s.user).Return(int64(1), nil)
	s.expectRenderFilter("filter1", map[string]interface{}{}, "output")
	s.expectRender("view-filter", pages.ContextData{
		"filter":       filter1,
		"params":       map[string]interface{}{},
		"rendered":     "output",
		"has_instance": true,
		"saved_ok":     true,
	})
	s.runRequest(req, assertOk)
}

func (s *ServerTestSuite) TestViewFilter_CreateListToo() {
	f := buildFilter2FormBody()
	f.Add("__save", "")
	req := httptest.NewRequest(http.MethodPost, "/filters/filter2", strings.NewReader(f.Encode()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	req.AddCookie(verifiedCookie)
	s.expectF.GetFilter("filter2").Return(filter2, nil)

	params := map[string]interface{}{
		"one":   "1",
		"two":   false,
		"three": []string{"option1", "option2"},
	}
	paramsB := pgtype.JSONB{}
	s.NoError(paramsB.Set(&params))
	s.expectQ.CountListsForUser(gomock.Any(), s.user).Return(int64(0), nil)
	s.expectQ.CreateListForUser(gomock.Any(), s.user).Return(uuid.New(), nil)
	s.expectQ.CountInstanceForUserAndFilter(gomock.Any(), db.CountInstanceForUserAndFilterParams{
		UserID:     s.user,
		FilterName: "filter2",
	}).Return(int64(0), nil)
	s.expectQ.CreateInstanceForUserAndFilter(gomock.Any(), db.CreateInstanceForUserAndFilterParams{
		UserID:     s.user,
		FilterName: "filter2",
		Params:     paramsB,
	}).Return(nil)
	s.expectRenderFilter("filter2", params, "output")
	s.expectRender("view-filter", pages.ContextData{
		"filter":       filter2,
		"params":       params,
		"rendered":     "output",
		"has_instance": true,
		"saved_ok":     true,
	})
	s.runRequest(req, assertOk)
}

func (s *ServerTestSuite) TestViewFilter_Update() {
	f := buildFilter2FormBody()
	f.Add("__save", "")
	req := httptest.NewRequest(http.MethodPost, "/filters/filter2", strings.NewReader(f.Encode()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	req.AddCookie(verifiedCookie)
	s.expectF.GetFilter("filter2").Return(filter2, nil)

	params := map[string]interface{}{
		"one":   "1",
		"two":   false,
		"three": []string{"option1", "option2"},
	}
	paramsB := pgtype.JSONB{}
	s.NoError(paramsB.Set(&params))
	s.expectQ.CountInstanceForUserAndFilter(gomock.Any(), db.CountInstanceForUserAndFilterParams{
		UserID:     s.user,
		FilterName: "filter2",
	}).Return(int64(1), nil)
	s.expectQ.UpdateInstanceForUserAndFilter(gomock.Any(), db.UpdateInstanceForUserAndFilterParams{
		UserID:     s.user,
		FilterName: "filter2",
		Params:     paramsB,
	}).Return(nil)
	s.expectRenderFilter("filter2", params, "output")
	s.expectRender("view-filter", pages.ContextData{
		"filter":       filter2,
		"params":       params,
		"rendered":     "output",
		"has_instance": true,
		"saved_ok":     true,
	})
	s.runRequest(req, assertOk)
}

func (s *ServerTestSuite) TestViewFilter_Disable() {
	f := make(url.Values)
	f.Add("__disable", "")
	req := httptest.NewRequest(http.MethodPost, "/filters/filter2", strings.NewReader(f.Encode()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	req.AddCookie(verifiedCookie)
	s.expectF.GetFilter("filter2").Return(filter2, nil)
	s.expectQ.DeleteInstanceForUserAndFilter(gomock.Any(), db.DeleteInstanceForUserAndFilterParams{
		UserID:     s.user,
		FilterName: "filter2",
	}).Return(nil)

	s.runRequest(req, func(t *testing.T, rec *httptest.ResponseRecorder) {
		assert.Equal(t, 302, rec.Code)
		assert.Equal(t, "/filters", rec.Header().Get("Location"))
	})
}

func (s *ServerTestSuite) TestViewFilterRender_Defaults() {
	req := httptest.NewRequest(http.MethodPost, "/filters/filter2/render", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	req.AddCookie(verifiedCookie)
	s.expectF.GetFilter("filter2").Return(filter2, nil)

	s.expectRenderFilter("filter2", filter2Defaults, "output")
	s.expectRender("view-filter-render", pages.ContextData{
		"rendered": "output",
	})
	s.runRequest(req, assertOk)
}

func (s *ServerTestSuite) TestViewFilterRender_Params() {
	req := httptest.NewRequest(http.MethodPost, "/filters/filter2/render", strings.NewReader(buildFilter2FormBody().Encode()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	req.AddCookie(verifiedCookie)
	s.expectF.GetFilter("filter2").Return(filter2, nil)

	params := map[string]interface{}{
		"one":   "1",
		"two":   false,
		"three": []string{"option1", "option2"},
	}

	s.expectRenderFilter("filter2", params, "output")
	s.expectRender("view-filter-render", pages.ContextData{
		"rendered": "output",
	})
	s.runRequest(req, assertOk)
}

func buildFilter2FormBody() url.Values {
	f := make(url.Values)
	f.Add("one", "1")
	f.Add("two", "off")
	f.Add("three", "option1")
	f.Add("three", "option2")
	return f
}

func (s *ServerTestSuite) expectRenderFilter(name string, params interface{}, output string) {
	s.expectF.Render(gomock.Any(), name, params).
		DoAndReturn(func(w io.Writer, _ string, _ map[string]interface{}) error {
			_, err := w.Write([]byte(output))
			return err
		})
}
