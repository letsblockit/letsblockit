package server

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/letsblockit/letsblockit/src/pages"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func (s *ServerTestSuite) TestUserAccount_Ok() {
	token, err := s.store.CreateListForUser(context.Background(), s.user)
	require.NoError(s.T(), err)
	s.markListDownloaded()
	require.NoError(s.T(), s.server.upsertFilterParams(s.c, s.user, "one", nil))
	require.NoError(s.T(), s.server.upsertFilterParams(s.c, s.user, "two", nil))

	req := httptest.NewRequest(http.MethodGet, "/user/account", nil)
	req.AddCookie(verifiedCookie)

	s.expectRender("user-account", pages.ContextData{
		"filter_count":    int64(2),
		"list_token":      token.String(),
		"list_downloaded": true,
	})
	s.runRequest(req, func(t *testing.T, rec *httptest.ResponseRecorder) {
		assert.Equal(t, 200, rec.Code, rec.Body)
		assert.Len(t, rec.Result().Cookies(), 2)
		cookie := rec.Result().Cookies()[0]
		assert.Equal(t, "has_account", cookie.Name)
		assert.Equal(t, "true", cookie.Value)
		assert.Equal(t, "/", cookie.Path)
	})
}

func (s *ServerTestSuite) TestUserAccount_CreateList() {
	req := httptest.NewRequest(http.MethodGet, "/user/account", nil)
	req.AddCookie(verifiedCookie)

	// First query
	s.expectP.Render(gomock.Any(), "user-account", gomock.Any())
	s.runRequest(req, assertOk)

	// Check that a list has been created for the user
	list, err := s.store.GetListForUser(context.Background(), s.user)
	require.NoError(s.T(), err)

	// Second query asserting the token value
	s.expectRender("user-account", pages.ContextData{
		"filter_count":    int64(0),
		"list_downloaded": false,
		"list_token":      list.Token.String(),
	})
	s.runRequest(req, assertOk)
}

func (s *ServerTestSuite) TestUserAccount_Anonymous() {
	req := httptest.NewRequest(http.MethodGet, "/user/account", nil)
	s.expectRender("user-account", nil)
	s.runRequest(req, assertOk)
}
