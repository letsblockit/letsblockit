package server

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func (s *ServerTestSuite) TestRenderList_NotFound() {
	token := uuid.New()
	req := httptest.NewRequest(http.MethodGet, "/list/"+token.String(), nil)
	s.runRequest(req, func(t *testing.T, rec *httptest.ResponseRecorder) {
		assert.Equal(t, 404, rec.Code)
	})
}

func (s *ServerTestSuite) TestRenderList_OK() {
	token, err := s.store.CreateListForUser(context.Background(), s.user)
	require.NoError(s.T(), err)

	require.NoError(s.T(), s.server.upsertFilterParams(s.c, s.user, "filter2", filter2Custom))
	require.NoError(s.T(), s.server.upsertFilterParams(s.c, s.user, "filter1", nil))
	require.NoError(s.T(), s.server.upsertFilterParams(s.c, s.user, "custom-rules", nil))

	req := httptest.NewRequest(http.MethodGet, "/list/"+token.String(), nil)
	rec := httptest.NewRecorder()
	s.server.echo.ServeHTTP(rec, req)
	s.Equal(200, rec.Code)
	s.Equal(`! Title: letsblock.it - My filters
! Expires: 12 hours
! Homepage: https://letsblock.it
! License: https://github.com/letsblockit/letsblockit/blob/main/LICENSE.txt

! filter1
hello from one
! filter2
hello one blep
hello two blep

! custom-rules
custom`, rec.Body.String())

	list, err := s.store.GetListForUser(context.Background(), s.user)
	require.NoError(s.T(), err)
	require.True(s.T(), list.Downloaded)
}

func (s *ServerTestSuite) TestRenderList_WithReferer() {
	token, err := s.store.CreateListForUser(context.Background(), s.user)
	require.NoError(s.T(), err)
	req := httptest.NewRequest(http.MethodGet, "/list/"+token.String(), nil)
	req.Header.Set("Referer", "https://letsblock.it/user/account")
	rec := httptest.NewRecorder()
	s.server.echo.ServeHTTP(rec, req)
	s.Equal(200, rec.Code)

	list, err := s.store.GetListForUser(context.Background(), s.user)
	require.NoError(s.T(), err)
	require.False(s.T(), list.Downloaded)
}

func (s *ServerTestSuite) TestRenderList_BannedUser() {
	s.setUserBanned()
	token, err := s.store.CreateListForUser(context.Background(), s.user)
	require.NoError(s.T(), err)

	req := httptest.NewRequest(http.MethodGet, "/list/"+token.String(), nil)
	req.Header.Set("Referer", "https://letsblock.it/user/account")
	rec := httptest.NewRecorder()
	s.server.echo.ServeHTTP(rec, req)
	s.Equal(403, rec.Code)
}

func (s *ServerTestSuite) TestExportList_OK() {
	params := map[string]any{
		"one":   "blep",
		"two":   false,
		"three": []any{"one", "two"},
	}
	require.NoError(s.T(), s.server.upsertFilterParams(s.c, s.user, "filter2", params))
	require.NoError(s.T(), s.server.upsertFilterParams(s.c, s.user, "filter1", nil))
	require.NoError(s.T(), s.server.upsertFilterParams(s.c, s.user, "custom-rules", nil))

	list, err := s.store.GetListForUser(context.Background(), s.user)
	require.NoError(s.T(), err)

	req := httptest.NewRequest(http.MethodGet, "/export/"+list.Token.String(), nil)

	rec := httptest.NewRecorder()
	s.server.echo.ServeHTTP(rec, req)
	s.Equal(200, rec.Code)
	fmt.Println(rec.Body.String())
	s.Equal(fmt.Sprintf(`# letsblock.it filter list export
#
# List token: %s
# Export date: 2020-06-02
#
# You can edit this file and render it locally, check out instructions at:
# https://github.com/letsblockit/letsblockit/tree/main/cmd/render/README.md

title: My filters
instances:
- filter: filter1
- filter: filter2
  params:
    one: blep
    three:
    - one
    - two
    two: false
- filter: custom-rules
`, list.Token), rec.Body.String())
}

func (s *ServerTestSuite) TestExportList_BadUser() {
	otherUser := uuid.New().String()
	for {
		if otherUser != s.user {
			break
		}
		otherUser = uuid.New().String()
	}

	token, err := s.store.CreateListForUser(context.Background(), otherUser)
	require.NoError(s.T(), err)
	req := httptest.NewRequest(http.MethodGet, "/export/"+token.String(), nil)

	rec := httptest.NewRecorder()
	s.server.echo.ServeHTTP(rec, req)
	s.Equal(403, rec.Code)
}
