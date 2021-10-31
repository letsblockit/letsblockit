package server

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/xvello/letsblockit/src/store"
)

func (s *ServerTestSuite) TestRenderList_NotFound() {
	req := httptest.NewRequest(http.MethodGet, "/list/invalid", nil)
	s.expectS.GetListForToken("invalid").Return(nil, store.ErrRecordNotFound)
	s.runRequest(req, func(t *testing.T, rec *httptest.ResponseRecorder) {
		assert.Equal(t, 404, rec.Code)
	})
}

func (s *ServerTestSuite) TestRenderList_OK() {
	req := httptest.NewRequest(http.MethodGet, "/list/mytoken", nil)
	s.expectS.GetListForToken("mytoken").Return(&store.FilterList{
		Name:  "List name",
		Token: "mytoken",
		FilterInstances: []*store.FilterInstance{{
			FilterName: "two",
			Params:     map[string]interface{}{"a": 1, "b": 2},
		}, {
			FilterName: "custom-rules",
		}, {
			FilterName: "one",
		}},
	}, nil)

	rec := httptest.NewRecorder()
	s.expectF.Render(gomock.Any(), "one", nil).
		DoAndReturn(func(w io.Writer, _ string, _ map[string]interface{}) error {
			_, err := w.Write([]byte("content1"))
			return err
		})
	s.expectF.Render(gomock.Any(), "two", map[string]interface{}{"a": 1, "b": 2}).
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
	s.Equal(rec.Body.String(), `! Title: letsblock.it - List name
! Expires: 1 hour
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
