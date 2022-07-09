package server

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
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

	params := map[string]interface{}{"a": "1", "b": "2"}
	require.NoError(s.T(), s.server.upsertFilterParams(s.c, s.user, "two", params))
	require.NoError(s.T(), s.server.upsertFilterParams(s.c, s.user, "one", nil))
	require.NoError(s.T(), s.server.upsertFilterParams(s.c, s.user, "custom-rules", nil))

	req := httptest.NewRequest(http.MethodGet, "/list/"+token.String(), nil)
	rec := httptest.NewRecorder()
	s.expectF.Render(gomock.Any(), "one", nil).
		DoAndReturn(func(w io.Writer, _ string, _ map[string]interface{}) error {
			_, err := w.Write([]byte("content1"))
			return err
		})
	s.expectF.Render(gomock.Any(), "two", gomock.Eq(params)).
		DoAndReturn(func(w io.Writer, _ string, _ map[string]interface{}) error {
			_, err := w.Write([]byte("content2\nmultiline"))
			return err
		})
	s.expectF.Render(gomock.Any(), "custom-rules", nil).
		DoAndReturn(func(w io.Writer, _ string, _ map[string]interface{}) error {
			_, err := w.Write([]byte("custom"))
			return err
		})
	s.server.echo.ServeHTTP(rec, req)
	s.Equal(200, rec.Code)
	s.Equal(rec.Body.String(), `! Title: letsblock.it - My filters
! Expires: 12 hours
! Homepage: https://letsblock.it
! License: https://github.com/letsblockit/letsblockit/blob/main/LICENSE.txt

! one
content1
! two
content2
multiline
! custom-rules
custom`)

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
	params := map[string]interface{}{"a": "1", "b": "2"}
	require.NoError(s.T(), s.server.upsertFilterParams(s.c, s.user, "two", params))
	require.NoError(s.T(), s.server.upsertFilterParams(s.c, s.user, "one", nil))
	require.NoError(s.T(), s.server.upsertFilterParams(s.c, s.user, "custom-rules", nil))

	list, err := s.store.GetListForUser(context.Background(), s.user)
	require.NoError(s.T(), err)

	req := httptest.NewRequest(http.MethodGet, "/export/"+list.Token.String(), nil)
	req.AddCookie(verifiedCookie)

	rec := httptest.NewRecorder()
	s.server.echo.ServeHTTP(rec, req)
	s.Equal(200, rec.Code)
	fmt.Println(rec.Body.String())
	s.Equal(rec.Body.String(), fmt.Sprintf(`# letsblock.it filter list export
#
# List token: %s
# Export date: 2020-06-02
#
# You can edit this file and render it locally, check out instructions at:
# https://github.com/letsblockit/letsblockit/tree/main/cmd/render/README.md

title: My filters
instances:
- filter: one
- filter: two
  params:
    a: "1"
    b: "2"
- filter: custom-rules
`, list.Token))
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
	req.AddCookie(verifiedCookie)

	rec := httptest.NewRecorder()
	s.server.echo.ServeHTTP(rec, req)
	s.Equal(403, rec.Code)
}
