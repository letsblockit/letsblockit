package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xvello/letsblockit/src/pages"
	"github.com/xvello/letsblockit/src/store"
)

func (s *ServerTestSuite) TestLogin_OK() {
	req := httptest.NewRequest(http.MethodGet, "/user/login", nil)
	s.expectRender("user-login", nil)
	s.runRequest(req, assertOk)
}

func (s *ServerTestSuite) TestLogin_RedirectRaw() {
	req := httptest.NewRequest(http.MethodGet, "/user/login", nil)
	s.login(true)
	s.runRequest(req, func(t *testing.T, rec *httptest.ResponseRecorder) {
		assert.Equal(t, 302, rec.Code)
		assert.Equal(t, "/user/account", rec.Header().Get("Location"))
	})
}

func (s *ServerTestSuite) TestLogin_RedirectHTMX() {
	req := httptest.NewRequest(http.MethodGet, "/user/login", nil)
	req.Header.Set("HX-Request", "true")
	s.login(true)
	s.runRequest(req, func(t *testing.T, rec *httptest.ResponseRecorder) {
		assert.Equal(t, 200, rec.Code)
		assert.Equal(t, "/user/account", rec.Header().Get("HX-Redirect"))
	})
}

func (s *ServerTestSuite) TestUserAccount_Verified() {
	list := &store.FilterList{
		UserID: s.user,
		Name:   "test",
	}
	req := httptest.NewRequest(http.MethodGet, "/user/account", nil)
	s.login(true)
	s.expectS.CountFilters(s.user).Return(int64(5), nil)
	s.expectS.GetOrCreateFilterList(s.user).Return(list, nil)
	s.expectRender("user-account", pages.ContextData{
		"filter_count": int64(5),
		"filter_list":  list,
	})
	s.runRequest(req, assertOk)
}

func (s *ServerTestSuite) TestUserAccount_NotVerified() {
	req := httptest.NewRequest(http.MethodGet, "/user/account", nil)
	s.login(false)
	s.expectRender("user-account", nil)
	s.runRequest(req, assertOk)
}

func (s *ServerTestSuite) TestUserAccount_Redirect() {
	req := httptest.NewRequest(http.MethodGet, "/user/account", nil)
	s.runRequest(req, func(t *testing.T, rec *httptest.ResponseRecorder) {
		assert.Equal(t, 302, rec.Code)
		assert.Equal(t, "/user/login", rec.Header().Get("Location"))
	})
}
