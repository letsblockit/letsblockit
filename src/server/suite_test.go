package server

import (
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
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/xvello/letsblockit/src/db"
	"github.com/xvello/letsblockit/src/news"
	"github.com/xvello/letsblockit/src/pages"
	"github.com/xvello/letsblockit/src/server/mocks"
)

var (
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
)

type mockStore struct {
	*mocks.MockQuerier
}

func (m mockStore) RunTx(e echo.Context, f db.TxFunc) error {
	return f(e.Request().Context(), m)
}

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
	server       *Server
	expectF      *mocks.MockFilterRepositoryMockRecorder
	expectP      *mocks.MockPageRendererMockRecorder
	expectQ      *mocks.MockQuerierMockRecorder
	expectR      *mocks.MockReleaseClientMockRecorder
	expectUP     *mocks.MockUserPreferenceManagerMockRecorder
	kratosServer *httptest.Server
	user         uuid.UUID
	csrf         string
	releases     []*news.Release
	preferences  *db.UserPreference
}

func (s *ServerTestSuite) SetupTest() {
	c := gomock.NewController(s.T())
	fm := mocks.NewMockFilterRepository(c)
	pm := mocks.NewMockPageRenderer(c)
	qm := mocks.NewMockQuerier(c)
	rm := mocks.NewMockReleaseClient(c)
	upm := mocks.NewMockUserPreferenceManager(c)
	s.expectF = fm.EXPECT()
	s.expectP = pm.EXPECT()
	s.expectQ = qm.EXPECT()
	s.expectR = rm.EXPECT()
	s.expectUP = upm.EXPECT()

	s.kratosServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var err error
		fmt.Println(r.URL.Path)
		switch r.URL.Path {
		case oryLogoutInfoPath:
			_, err = fmt.Fprint(w, `{"logout_url":"targetURL"}`)
		case oryWhoamiPath:
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

	s.user = uuid.New()
	s.csrf = random.String(32)
	s.server = &Server{
		assets:  nil,
		echo:    echo.New(),
		filters: fm,
		now:     func() time.Time { return fixedNow },
		options: &Options{
			KratosURL: s.kratosServer.URL,
			silent:    true,
		},
		pages:       pm,
		preferences: upm,
		releases:    rm,
		statsd:      &statsd.NoOpClient{},
		store:       &mockStore{qm},
	}
	s.server.setupRouter()

	// Preferences and releases are called by buildPageContext for logged-in users.
	// Add catch-all expectations to avoid noise in the tests.
	// Values can be set by tests before running a query
	s.preferences = nil
	s.releases = nil
	s.expectUP.Get(gomock.Any(), s.user).DoAndReturn(func(c echo.Context, user uuid.UUID) (*db.UserPreference, error) {
		if user == s.user {
			return s.preferences, nil
		} else {
			return nil, db.NotFound
		}
	}).AnyTimes()
	s.expectR.GetLatestAt().DoAndReturn(func() (time.Time, error) {
		if len(s.releases) > 0 {
			return s.releases[0].CreatedAt, nil
		} else {
			return fixedNow, nil
		}
	}).AnyTimes()
	s.expectR.GetReleases().DoAndReturn(func() ([]*news.Release, error) {
		return s.releases, nil
	}).AnyTimes()
}

func (s *ServerTestSuite) setUserBanned() {
	if s.server.banned == nil {
		s.server.banned = make(map[uuid.UUID]struct{})
	}
	s.server.banned[s.user] = struct{}{}
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

func assertRedirect(target string) func(t *testing.T, rec *httptest.ResponseRecorder) {
	return func(t *testing.T, rec *httptest.ResponseRecorder) {
		assert.Equal(t, 302, rec.Code, rec.Body)
		assert.Equal(t, target, rec.Header().Get("Location"))
	}
}

func assertSeeOther(target string) func(t *testing.T, rec *httptest.ResponseRecorder) {
	return func(t *testing.T, rec *httptest.ResponseRecorder) {
		assert.Equal(t, 303, rec.Code, rec.Body)
		assert.Equal(t, target, rec.Header().Get("Location"))
	}
}

func (s *ServerTestSuite) runRequest(req *http.Request, checks func(*testing.T, *httptest.ResponseRecorder)) {
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

func TestServerTestSuite(t *testing.T) {
	suite.Run(t, new(ServerTestSuite))
}
