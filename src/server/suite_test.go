package server

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/xvello/letsblockit/src/pages"
)

var (
	verifiedCookie = &http.Cookie{
		Name:  "ory_session_verified",
		Value: "true",
	}
	unverifiedCookie = &http.Cookie{
		Name:  "ory_session_verified",
		Value: "false",
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
	server    *Server
	expectF   *MockFilterRepositoryMockRecorder
	expectP   *MockPageRendererMockRecorder
	expectS   *MockDataStoreMockRecorder
	oryServer *httptest.Server
	user      string
}

func (s *ServerTestSuite) SetupTest() {
	c := gomock.NewController(s.T())
	fm := NewMockFilterRepository(c)
	pm := NewMockPageRenderer(c)
	sm := NewMockDataStore(c)
	s.expectF = fm.EXPECT()
	s.expectP = pm.EXPECT()
	s.expectS = sm.EXPECT()

	oryClientRetries = 0
	s.oryServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case oryLogoutInfoPath:
			fmt.Fprint(w, `{"logout_token":"token"}`)
		case oryWhoamiPath:
			cookie, _ := r.Cookie("ory_session_verified")
			fmt.Fprintf(w, whoAmiPattern, s.user, cookie.Value)
		default:
			http.NotFound(w, r)
		}
	}))

	s.user = uuid.NewString()
	s.server = &Server{
		assets: nil,
		echo:   echo.New(),
		options: &Options{
			OryUrl: s.oryServer.URL,
			silent: true,
		},
		store:   sm,
		filters: fm,
		pages:   pm,
	}
	s.server.setupRouter()
}

func (s *ServerTestSuite) expectRender(page string, data pages.ContextData) *gomock.Call {
	return s.expectP.Render(gomock.Any(), page, &pageDataMatcher{
		t:    s.T(),
		data: data,
	})
}

func (s *ServerTestSuite) expectRenderWithContext(page string, ctx *pages.Context) *gomock.Call {
	return s.expectP.Render(gomock.Any(), page, gomock.Eq(ctx))
}

func assertOk(t *testing.T, rec *httptest.ResponseRecorder) {
	assert.Equal(t, http.StatusOK, rec.Code, rec.Body)
}

func (s *ServerTestSuite) runRequest(req *http.Request, checks func(*testing.T, *httptest.ResponseRecorder)) {
	rec := httptest.NewRecorder()
	s.server.echo.ServeHTTP(rec, req)
	checks(s.T(), rec)
}

func TestServerTestSuite(t *testing.T) {
	suite.Run(t, new(ServerTestSuite))
}
