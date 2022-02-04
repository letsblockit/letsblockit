package server

import (
	"encoding/base64"
	"io"
	"io/fs"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/suite"
)

const cssDir = "assets/dist"
const cssFile = cssDir + "/main.css"

type WrappedAssetsTestSuite struct {
	suite.Suite
	assets *wrappedAssets
}

func (s *WrappedAssetsTestSuite) TestEmbedIsSeekable() {
	file, err := s.assets.Open(cssFile)
	s.NoError(err)
	// embed.Open currently returns a seekable, ensure this is true in future go versions
	pos, err := file.Seek(0, io.SeekEnd)
	s.NoError(err)
	s.True(pos > 0)
}

func (s *WrappedAssetsTestSuite) TestOpenDir() {
	// Confirm we can open and stat a regular file
	f, err := s.assets.Open(cssFile)
	s.NoError(err)
	info, err := f.Stat()
	s.NoError(err)
	s.False(info.IsDir())

	// Confirm we deny opening the parent directory
	f, err = s.assets.Open(cssDir)
	s.Equal(fs.ErrNotExist, err)
	s.Nil(f)
}

func (s *WrappedAssetsTestSuite) TestHash() {
	out, err := base64.RawURLEncoding.DecodeString(s.assets.hash)
	s.NoError(err)
	s.Len(out, 16) // 128 bits hash expected
}

func (s *WrappedAssetsTestSuite) TestETag() {
	expectedETag := "\"" + s.assets.hash + "\""
	s.Equal(expectedETag, s.assets.eTag)
}

func (s *WrappedAssetsTestSuite) TestServe_OK() {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/"+cssFile, nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	s.NoError(s.assets.serve(c))
	s.Equal(http.StatusOK, rec.Result().StatusCode)
	s.Equal(s.assets.eTag, rec.Header().Get("etag"))
	s.Equal("text/css; charset=utf-8", rec.Header().Get("content-type"))
	s.Contains(rec.Header().Get("cache-control"), "max-age=")
}

func (s *WrappedAssetsTestSuite) TestServe_NotModified() {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/"+cssFile, nil)
	req.Header.Set("If-None-Match", s.assets.eTag)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	s.NoError(s.assets.serve(c))
	s.Equal(http.StatusNotModified, rec.Result().StatusCode)
}

func (s *WrappedAssetsTestSuite) TestServe_NotFound() {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/un/kn/own/fi/le", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	s.NoError(s.assets.serve(c))
	s.Equal(http.StatusNotFound, rec.Result().StatusCode)
}

func (s *WrappedAssetsTestSuite) TestServe_DenyFolder() {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/"+cssDir, nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	s.NoError(s.assets.serve(c))
	s.Equal(http.StatusNotFound, rec.Result().StatusCode)
}

func TestWrappedAssetsTestSuite(t *testing.T) {
	suite.Run(t, &WrappedAssetsTestSuite{
		assets: loadAssets(),
	})
}
