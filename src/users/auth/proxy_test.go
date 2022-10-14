package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/random"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ProxyAuthSuite struct {
	suite.Suite
	echo       *echo.Echo
	user       string
	headerName string
}

func (s *ProxyAuthSuite) SetupTest() {
	s.user = random.String(32)
	s.headerName = random.String(16)
	s.echo = echo.New()
	s.echo.Use(NewProxy(s.headerName).BuildMiddleware())
	s.echo.GET("/", func(c echo.Context) error {
		return c.String(200, "Hello "+GetUserId(c))
	})
}

func (s *ProxyAuthSuite) Test_OK() {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(s.headerName, s.user)
	rec := httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	assert.Equal(s.T(), http.StatusOK, rec.Code, rec.Body)
	assert.Equal(s.T(), "Hello "+s.user, rec.Body.String())
}

func (s *ProxyAuthSuite) Test_Unauthorized() {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(s.user, s.user)
	rec := httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	assert.Equal(s.T(), http.StatusUnauthorized, rec.Code, rec.Body)
}

func TestProxyAuthSuite(t *testing.T) {
	suite.Run(t, new(ProxyAuthSuite))
}
