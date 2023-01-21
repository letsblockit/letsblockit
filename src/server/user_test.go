package server

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/letsblockit/letsblockit/src/filters"
	"github.com/letsblockit/letsblockit/src/pages"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func (s *ServerTestSuite) TestUserAccount_Ok() {
	token, err := s.store.CreateListForUser(context.Background(), s.user)
	require.NoError(s.T(), err)
	s.markListDownloaded()
	require.NoError(s.T(), s.server.upsertFilterParams(s.c, s.user, &filters.Instance{Template: "one"}))
	require.NoError(s.T(), s.server.upsertFilterParams(s.c, s.user, &filters.Instance{Template: "two"}))

	req := httptest.NewRequest(http.MethodGet, "/user/account", nil)

	s.expectRender("user-account", pages.ContextData{
		"filter_count":    int64(2),
		"list_token":      token.String(),
		"list_downloaded": true,
	})
	s.runRequest(req, assertOk)
}

func (s *ServerTestSuite) TestUserAccount_CreateList() {
	req := httptest.NewRequest(http.MethodGet, "/user/account", nil)

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
	s.user = ""
	req := httptest.NewRequest(http.MethodGet, "/user/account", nil)
	s.expectRender("user-account", nil)
	s.runRequest(req, assertOk)
}

func (s *ServerTestSuite) TestRotateListToken_Ok() {
	token, err := s.store.CreateListForUser(context.Background(), s.user)
	require.NoError(s.T(), err)

	f := make(url.Values)
	f.Add("one", "blep")
	f.Add("token", token.String())
	f.Add("confirm", "on")
	f.Add(csrfLookup, s.csrf)
	req := httptest.NewRequest(http.MethodPost, "/user/rotate-token", strings.NewReader(f.Encode()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)

	s.expectP.Redirect(gomock.Any(), http.StatusSeeOther, "/user/account")
	s.runRequest(req, assertOk)

	list, err := s.store.GetListForUser(context.Background(), s.user)
	require.NoError(s.T(), err)
	require.NotEqual(s.T(), token, list.Token)
}

func (s *ServerTestSuite) TestRotateListToken_MissingCSRF() {
	token, err := s.store.CreateListForUser(context.Background(), s.user)
	require.NoError(s.T(), err)

	f := make(url.Values)
	f.Add("one", "blep")
	f.Add("token", token.String())
	f.Add("confirm", "on")
	req := httptest.NewRequest(http.MethodPost, "/user/rotate-token", strings.NewReader(f.Encode()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	s.runRequest(req, func(t *testing.T, recorder *httptest.ResponseRecorder) {
		assert.Equal(t, 400, recorder.Result().StatusCode)
	})
}
