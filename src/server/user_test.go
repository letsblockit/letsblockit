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
	s.runRequest(req, func(t *testing.T, rec *httptest.ResponseRecorder) {
		assert.Equal(t, 200, rec.Code, rec.Body)
		assert.Len(t, rec.Result().Cookies(), 1)
		cookie := rec.Result().Cookies()[0]
		assert.Equal(t, "has_account", cookie.Name)
		assert.Equal(t, "true", cookie.Value)
		assert.Equal(t, "/", cookie.Path)
	})
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
