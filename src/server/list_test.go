package server

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/jackc/pgtype"
	"github.com/stretchr/testify/assert"
	"github.com/xvello/letsblockit/src/db"
)

func (s *ServerTestSuite) TestRenderList_NotFound() {
	token := uuid.New()
	req := httptest.NewRequest(http.MethodGet, "/list/"+token.String(), nil)
	s.expectQ.GetListForToken(gomock.Any(), token).Return(db.GetListForTokenRow{}, db.NotFound)
	s.runRequest(req, func(t *testing.T, rec *httptest.ResponseRecorder) {
		assert.Equal(t, 404, rec.Code)
	})
}

func (s *ServerTestSuite) TestRenderList_OK() {
	token := uuid.New()
	req := httptest.NewRequest(http.MethodGet, "/list/"+token.String(), nil)
	s.expectQ.GetListForToken(gomock.Any(), token).Return(db.GetListForTokenRow{
		ID:         int32(10),
		Downloaded: true,
	}, nil)
	s.expectQ.MarkListDownloaded(gomock.Any(), int32(10)).Return(nil)

	params := map[string]interface{}{"a": "1", "b": "2"}
	paramsB := pgtype.JSONB{}
	s.NoError(paramsB.Set(&params))
	s.expectQ.GetInstancesForList(gomock.Any(), int32(10)).Return([]db.GetInstancesForListRow{{
		FilterName: "one",
	}, {
		FilterName: "custom-rules",
	}, {
		FilterName: "two",
		Params:     paramsB,
	}}, nil)
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
! License: https://github.com/xvello/letsblockit/blob/main/LICENSE.txt

! one
content1
! two
content2
multiline
! custom-rules
custom`)
}

func (s *ServerTestSuite) TestRenderList_WithReferer() {
	token := uuid.New()
	req := httptest.NewRequest(http.MethodGet, "/list/"+token.String(), nil)
	req.Header.Set("Referer", "https://letsblock.it/user/account")
	s.expectQ.GetListForToken(gomock.Any(), token).Return(db.GetListForTokenRow{
		ID:         int32(10),
		Downloaded: false,
	}, nil)
	s.expectQ.GetInstancesForList(gomock.Any(), int32(10)).Return([]db.GetInstancesForListRow{}, nil)
	rec := httptest.NewRecorder()
	s.server.echo.ServeHTTP(rec, req)
	s.Equal(200, rec.Code)
}

func (s *ServerTestSuite) TestRenderList_BannedUser() {
	s.setUserBanned()
	token := uuid.New()
	req := httptest.NewRequest(http.MethodGet, "/list/"+token.String(), nil)
	req.Header.Set("Referer", "https://letsblock.it/user/account")
	s.expectQ.GetListForToken(gomock.Any(), token).Return(db.GetListForTokenRow{
		ID:         int32(10),
		UserID:     s.user,
		Downloaded: false,
	}, nil)
	rec := httptest.NewRecorder()
	s.server.echo.ServeHTTP(rec, req)
	s.Equal(403, rec.Code)
}

func (s *ServerTestSuite) TestExportList_OK() {
	token, err := uuid.Parse("59dc36b7-aca5-46a1-a742-cf46dd2cac10")
	s.NoError(err)
	req := httptest.NewRequest(http.MethodGet, "/list/"+token.String()+"/export", nil)
	req.AddCookie(verifiedCookie)

	s.expectQ.GetListForToken(gomock.Any(), token).Return(db.GetListForTokenRow{
		ID:         int32(10),
		UserID:     s.user,
		Downloaded: true,
	}, nil)

	params := map[string]interface{}{"a": "1", "b": "2"}
	paramsB := pgtype.JSONB{}
	s.NoError(paramsB.Set(&params))
	s.expectQ.GetInstancesForList(gomock.Any(), int32(10)).Return([]db.GetInstancesForListRow{{
		FilterName: "one",
	}, {
		FilterName: "custom-rules",
	}, {
		FilterName: "two",
		Params:     paramsB,
	}}, nil)
	rec := httptest.NewRecorder()
	s.server.echo.ServeHTTP(rec, req)
	s.Equal(200, rec.Code)
	fmt.Println(rec.Body.String())
	s.Equal(rec.Body.String(), `# letsblock.it filter list export
#
# List token: 59dc36b7-aca5-46a1-a742-cf46dd2cac10
# Export date: 2020-06-02
#
# You can edit this file and render it locally, check out instructions at:
# https://github.com/xvello/letsblockit/tree/main/cmd/render/README.md

title: My filters
instances:
- filter: one
- filter: two
  params:
    a: "1"
    b: "2"
- filter: custom-rules
`)
}

func (s *ServerTestSuite) TestExportList_BadUser() {
	token, otherUser := uuid.New(), uuid.New()
	for {
		if otherUser != s.user {
			break
		}
		otherUser = uuid.New()
	}
	req := httptest.NewRequest(http.MethodGet, "/list/"+token.String()+"/export", nil)
	req.AddCookie(verifiedCookie)

	s.expectQ.GetListForToken(gomock.Any(), token).Return(db.GetListForTokenRow{
		ID:         int32(10),
		UserID:     otherUser,
		Downloaded: true,
	}, nil)
	rec := httptest.NewRecorder()
	s.server.echo.ServeHTTP(rec, req)
	s.Equal(403, rec.Code)
}
