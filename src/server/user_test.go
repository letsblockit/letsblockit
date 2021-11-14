package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/xvello/letsblockit/src/db"
	"github.com/xvello/letsblockit/src/pages"
)

func (s *ServerTestSuite) TestLogin_RedirectRegistration() {
	req := httptest.NewRequest(http.MethodGet, "/user/login", nil)
	s.runRequest(req, assertRedirect("/.ory/ui/registration"))
}

func (s *ServerTestSuite) TestLogin_RedirectLogin() {
	req := httptest.NewRequest(http.MethodGet, "/user/login", nil)
	req.AddCookie(&http.Cookie{
		Name:  "has_account",
		Value: "true",
	})
	s.runRequest(req, assertRedirect("/.ory/ui/login"))
}

func (s *ServerTestSuite) TestLogin_RedirectLoggedRaw() {
	req := httptest.NewRequest(http.MethodGet, "/user/login", nil)
	req.AddCookie(verifiedCookie)
	s.runRequest(req, assertRedirect("/user/account"))
}

func (s *ServerTestSuite) TestLogin_RedirectLoggedHTMX() {
	req := httptest.NewRequest(http.MethodGet, "/user/login", nil)
	req.Header.Set("HX-Request", "true")
	req.AddCookie(verifiedCookie)
	s.runRequest(req, func(t *testing.T, rec *httptest.ResponseRecorder) {
		assert.Equal(t, 200, rec.Code, rec.Body)
		assert.Equal(t, "/user/account", rec.Header().Get("HX-Redirect"))
	})
}

func (s *ServerTestSuite) TestUserAccount_Verified() {
	token := uuid.New()
	req := httptest.NewRequest(http.MethodGet, "/user/account", nil)
	req.AddCookie(verifiedCookie)
	s.expectQ.GetListForUser(gomock.Any(), s.user).Return(db.GetListForUserRow{
		Token:         token,
		InstanceCount: 5,
	}, nil)
	s.expectRender("user-account", pages.ContextData{
		"filter_count": int64(5),
		"list_token":   token.String(),
	})
	s.runRequest(req, assertOk)
}

func (s *ServerTestSuite) TestUserAccount_CreateList() {
	token := uuid.New()
	req := httptest.NewRequest(http.MethodGet, "/user/account", nil)
	req.AddCookie(verifiedCookie)
	s.expectQ.GetListForUser(gomock.Any(), s.user).Return(db.GetListForUserRow{}, db.NotFound)
	s.expectQ.CreateListForUser(gomock.Any(), s.user).Return(token, nil)
	s.expectRender("user-account", pages.ContextData{
		"filter_count": 0,
		"list_token":   token.String(),
	})
	s.runRequest(req, assertOk)
}

func (s *ServerTestSuite) TestUserAccount_NotVerified() {
	req := httptest.NewRequest(http.MethodGet, "/user/account", nil)
	req.AddCookie(unverifiedCookie)
	s.expectRender("user-account", nil)
	s.runRequest(req, assertOk)
}

func (s *ServerTestSuite) TestUserAccount_Redirect() {
	req := httptest.NewRequest(http.MethodGet, "/user/account", nil)
	s.runRequest(req, assertRedirect("/user/login"))
}

func (s *ServerTestSuite) TestUserLogout_OK() {
	req := httptest.NewRequest(http.MethodGet, "/user/logout", nil)
	req.AddCookie(verifiedCookie)
	s.runRequest(req, assertRedirect("targetURL"))
}

func (s *ServerTestSuite) TestUserLogout_Redirect() {
	req := httptest.NewRequest(http.MethodGet, "/user/logout", nil)
	s.runRequest(req, assertRedirect("/user/login"))
}
