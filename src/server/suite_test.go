package server

import (
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/DataDog/datadog-go/v5/statsd"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/random"
	"github.com/letsblockit/letsblockit/src/db"
	"github.com/letsblockit/letsblockit/src/filters"
	"github.com/letsblockit/letsblockit/src/news"
	"github.com/letsblockit/letsblockit/src/pages"
	"github.com/letsblockit/letsblockit/src/server/mocks"
	"github.com/letsblockit/letsblockit/src/users"
	"github.com/letsblockit/letsblockit/src/users/auth"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

var (
	//go:embed testdata/filters
	testFilters embed.FS
	fixedNow    = time.Date(2020, 06, 02, 17, 44, 22, 0, time.UTC)
)

type pageContextMatcher struct {
	t   *testing.T
	ctx *pages.Context
}

func (m *pageContextMatcher) Matches(x interface{}) bool {
	d, ok := x.(*pages.Context)
	if ok && m.ctx.RequestInfo == nil {
		m.ctx.RequestInfo = d.RequestInfo
	}
	return ok && assert.EqualValues(m.t, m.ctx, d)
}

func (m *pageContextMatcher) String() string {
	return fmt.Sprintf("is equal to %v (%T)", m.ctx, m.ctx)
}

type pageDataMatcher struct {
	t    *testing.T
	data pages.ContextData
}

func (m *pageDataMatcher) Matches(x interface{}) bool {
	d, ok := x.(*pages.Context)
	return ok && assert.EqualValues(m.t, m.data, d.Data)
}

func (m *pageDataMatcher) String() string {
	return fmt.Sprintf("is equal to %v (%T)", m.data, m.data)
}

type ServerTestSuite struct {
	suite.Suite
	c        echo.Context
	server   *Server
	expectP  *mocks.MockPageRendererMockRecorder
	expectR  *mocks.MockReleaseClientMockRecorder
	user     string
	csrf     string
	releases []*news.Release
	store    db.Store
}

// Implements auth.Backend: authenticates requests when s.user is not empty
func (s *ServerTestSuite) BuildMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if len(s.user) > 0 {
				c.Set("_user", s.user)
			}
			return next(c)
		}
	}
}

// Implements auth.Backend: do nothing
func (s *ServerTestSuite) RegisterRoutes(_ auth.EchoRouter) {}

func (s *ServerTestSuite) SetupTest() {
	c := gomock.NewController(s.T())
	pm := mocks.NewMockPageRenderer(c)
	rm := mocks.NewMockReleaseClient(c)
	s.expectP = pm.EXPECT()
	s.expectR = rm.EXPECT()

	s.c = echo.New().NewContext(httptest.NewRequest(http.MethodGet, "/", nil), httptest.NewRecorder())
	s.store = db.NewTestStore(s.T())
	s.csrf = random.String(32)
	s.user = uuid.New().String()
	pref, err := users.NewPreferenceManager(s.store)
	require.NoError(s.T(), err)
	require.NoError(s.T(), pref.UpdateNewsCursor(s.c, s.user, fixedNow))

	s.server = &Server{
		auth:    s,
		echo:    echo.New(),
		filters: filterRepo,
		now:     func() time.Time { return fixedNow },
		options: &Options{
			LogLevel: "off",
		},
		pages:       pm,
		preferences: pref,
		releases:    rm,
		statsd:      &statsd.NoOpClient{},
		store:       s.store,
	}
	s.server.setupRouter()

	// Preferences and releases are called by buildPageContext for logged-in users.
	// Add catch-all expectations to avoid noise in the tests.
	// Values can be set by tests before running a query
	s.releases = nil
	s.expectR.GetLatestAt().DoAndReturn(func() (time.Time, error) {
		if len(s.releases) > 0 {
			return s.releases[0].CreatedAt, nil
		} else {
			return fixedNow, nil
		}
	}).MinTimes(0)
	s.expectR.GetReleases().DoAndReturn(func() ([]*news.Release, error) {
		return s.releases, nil
	}).MinTimes(0)
}

func (s *ServerTestSuite) setUserBanned() {
	s.T().Helper()
	require.NoError(s.T(), s.store.AddUserBan(context.Background(), db.AddUserBanParams{
		UserID: s.user,
	}))
	s.server.bans, _ = users.LoadUserBans(s.server.store)
}

func (s *ServerTestSuite) markListDownloaded() {
	s.T().Helper()
	list, err := s.store.GetListForUser(context.Background(), s.user)
	require.NoError(s.T(), err)
	require.NoError(s.T(), s.store.MarkListDownloaded(context.Background(), list.Token))
}

func (s *ServerTestSuite) expectRender(page string, data pages.ContextData) *gomock.Call {
	return s.expectP.Render(gomock.Any(), page, &pageDataMatcher{
		t:    s.T(),
		data: data,
	})
}

func (s *ServerTestSuite) expectRenderWithSidebar(page, sidebar string, data pages.ContextData) *gomock.Call {
	return s.expectP.RenderWithSidebar(gomock.Any(), page, sidebar, &pageDataMatcher{
		t:    s.T(),
		data: data,
	})
}

func (s *ServerTestSuite) expectRenderWithSidebarAndContext(page, sidebar string, ctx *pages.Context) *gomock.Call {
	return s.expectP.RenderWithSidebar(gomock.Any(), page, sidebar, &pageContextMatcher{
		t:   s.T(),
		ctx: ctx,
	})
}

func (s *ServerTestSuite) expectRenderWithContext(page string, ctx *pages.Context) *gomock.Call {
	return s.expectP.Render(gomock.Any(), page, &pageContextMatcher{
		t:   s.T(),
		ctx: ctx,
	})
}

func assertOk(t *testing.T, rec *httptest.ResponseRecorder) {
	assert.Equal(t, http.StatusOK, rec.Code, rec.Body)
}

func (s *ServerTestSuite) runRequest(req *http.Request, checks func(*testing.T, *httptest.ResponseRecorder)) {
	s.T().Helper()
	rec := httptest.NewRecorder()
	if len(s.csrf) > 0 {
		req.AddCookie(&http.Cookie{
			Name:  csrfLookup,
			Value: s.csrf,
		})
	}
	s.server.echo.ServeHTTP(rec, req)
	checks(s.T(), rec)
}

func (s *ServerTestSuite) requireJSONEq(expected, actual any) {
	s.T().Helper()
	expectedJSON, err := json.Marshal(expected)
	require.NoError(s.T(), err)
	actualJSON, err := json.Marshal(actual)
	require.NoError(s.T(), err)
	require.JSONEq(s.T(), string(expectedJSON), string(actualJSON))
}

func (s *ServerTestSuite) requireInstanceCount(filter string, expected int64) {
	s.T().Helper()
	count, err := s.store.CountInstanceForUserAndFilter(context.Background(), db.CountInstanceForUserAndFilterParams{
		UserID:     s.user,
		FilterName: filter,
	})
	require.NoError(s.T(), err)
	require.EqualValues(s.T(), expected, count)
}

func TestServerTestSuite(t *testing.T) {
	// Load fixtures
	filterSources, err := fs.Sub(testFilters, "testdata/filters")
	require.NoError(t, err)
	filterRepo, err = filters.LoadFilters(filterSources)
	require.NoError(t, err)
	require.Len(t, filterRepo.GetFilters(), 3)
	filter1, _ = filterRepo.GetFilter("filter1")
	filter2, _ = filterRepo.GetFilter("filter2")
	filter3, _ = filterRepo.GetFilter("custom-rules")

	// Run test suite
	suite.Run(t, new(ServerTestSuite))
}
