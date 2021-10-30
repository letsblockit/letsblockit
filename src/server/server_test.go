package server

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/xvello/letsblockit/src/store"
)

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
	assert.Equal(t, ErrDryRunFinished, server.Start())
}

type ServerTestSuite struct {
	suite.Suite
	server *Server
	store  *store.Store
	fm     *MockfilterRepository
	pm     *MockpageRenderer
	user   string
}

func (s *ServerTestSuite) SetupTest() {
	c := gomock.NewController(s.T())
	s.fm = NewMockfilterRepository(c)
	s.pm = NewMockpageRenderer(c)
	s.user = uuid.NewString()

	var err error
	s.store, err = store.NewMemStore()
	s.Require().NoError(err)

	s.server = &Server{
		assets:  nil,
		echo:    echo.New(),
		options: &Options{},
		store:   s.store,
		filters: s.fm,
		pages:   s.pm,
	}
	s.server.setupRouter()
}

func (s *ServerTestSuite) fExpect() *MockfilterRepositoryMockRecorder {
	return s.fm.EXPECT()
}

func (s *ServerTestSuite) expectRender(page string, ctx map[string]interface{}) *gomock.Call {
	return s.pm.EXPECT().Render(gomock.Any(), page, gomock.Eq(ctx))
}

func (s *ServerTestSuite) login(verified bool) {
	s.server.echo.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set(
				userContextKey,
				&oryUser{
					Active: true,
					Identity: struct {
						Id        string
						Addresses []struct {
							Verified bool
						} `json:"verifiable_addresses"`
					}{
						Id: s.user,
						Addresses: []struct {
							Verified bool
						}{{
							Verified: verified,
						}},
					},
				},
			)
			return next(c)
		}
	})
}

func (s *ServerTestSuite) TestAbout_Anonymous() {
	req := httptest.NewRequest(http.MethodGet, "/about", nil)
	rec := httptest.NewRecorder()
	s.expectRender("about", map[string]interface{}{
		"navCurrent": "about",
		"navLinks":   navigationLinks,
		"title":      "About: Let’s block it!",
	})
	s.server.echo.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)
}

func (s *ServerTestSuite) TestAbout_LoggedVerified() {
	req := httptest.NewRequest(http.MethodGet, "/about", nil)
	rec := httptest.NewRecorder()
	s.login(true)
	s.expectRender("about", map[string]interface{}{
		"navCurrent": "about",
		"navLinks":   navigationLinks,
		"title":      "About: Let’s block it!",
		"logged":     true,
		"verified":   true,
	})
	s.server.echo.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)
}

func TestServerTestSuite(t *testing.T) {
	suite.Run(t, new(ServerTestSuite))
}
