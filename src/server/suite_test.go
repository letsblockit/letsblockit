package server

import (
	"context"
	"embed"
	"encoding/json"
	"fmt"
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

	fixedNow       = time.Date(2020, 06, 02, 17, 44, 22, 0, time.UTC)
	verifiedCookie = &http.Cookie{
		Name:  "ory_session_verified",
		Value: "true",
	}
	whoAmiPattern = `{
	  "id": "af9b460f-4ca0-453d-8bc7-cf68f30d4174",
	  "active": true,
	  "identity": {
		"id": "%s",
		"verifiable_addresses": [
		  {
			"verified": %s
		  }
		]
	  }
	}`

	// Filled from SetupTest
	filter1 = &filters.Filter{}
	filter2 = &filters.Filter{}
	filter3 = &filters.Filter{}
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
	c            echo.Context
	server       *Server
	expectP      *mocks.MockPageRendererMockRecorder
	expectR      *mocks.MockReleaseClientMockRecorder
	kratosServer *httptest.Server
	user         string
	csrf         string
	releases     []*news.Release
	store        db.Store
}

func (s *ServerTestSuite) SetupTest() {
	c := gomock.NewController(s.T())
	pm := mocks.NewMockPageRenderer(c)
	rm := mocks.NewMockReleaseClient(c)
	s.expectP = pm.EXPECT()
	s.expectR = rm.EXPECT()

	s.c = echo.New().NewContext(httptest.NewRequest(http.MethodGet, "/", nil), httptest.NewRecorder())
	s.store = db.NewTestStore(s.T())
	s.kratosServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var err error
		fmt.Println(r.URL.Path)
		switch r.URL.Path {
		case "/self-service/logout/browser":
			_, err = fmt.Fprint(w, `{"logout_url":"targetURL"}`)
		case "/sessions/whoami":
			cookie, _ := r.Cookie("ory_session_verified")
			_, err = fmt.Fprintf(w, whoAmiPattern, s.user, cookie.Value)
		case "/self-service/login/flows":
			switch r.URL.RawQuery {
			case "id=123456":
				_, err = fmt.Fprint(w, `{"ui":{"a": "1", "b": "2"},"return_to":"https://target"}`)
			case "id=666":
				_, err = fmt.Fprint(w, `{"invalid": true}`)
			}
		default:
			http.NotFound(w, r)
		}
		s.NoError(err)
	}))

	repo, err := filters.LoadFilters(testFilters)
	require.NoError(s.T(), err)
	require.NotEmpty(s.T(), repo.GetFilters())
	filter1, _ = repo.GetFilter("filter1")
	filter2, _ = repo.GetFilter("filter2")
	filter3, _ = repo.GetFilter("custom-rules")

	s.user = uuid.New().String()
	pref, err := users.NewPreferenceManager(s.store)
	require.NoError(s.T(), err)
	require.NoError(s.T(), pref.UpdateNewsCursor(s.c, s.user, fixedNow))
	s.csrf = random.String(32)
	s.server = &Server{
		auth:    auth.NewOryBackend(s.kratosServer.URL, pm, &statsd.NoOpClient{}),
		echo:    echo.New(),
		filters: repo,
		now:     func() time.Time { return fixedNow },
		options: &Options{
			AuthKratosUrl: s.kratosServer.URL,
			LogLevel:      "off",
		},
		pages:       pm,
		preferences: pref,
		releases:    rm,
		statsd:      &statsd.NoOpClient{},
		store:       s.store,
	}
	s.server.setupRouter()

	// TODO: Used by ory_test.go, remove after moving these tests to their own suite
	s.expectP.BuildPageContext(gomock.Any(), gomock.Any()).DoAndReturn(s.server.buildPageContext).
		MinTimes(0).MaxTimes(1)

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

func (s *ServerTestSuite) TearDownTest() {
	s.kratosServer.Close()
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
	suite.Run(t, new(ServerTestSuite))
}
