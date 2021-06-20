package server

import (
	"embed"
	"encoding/base64"
	"fmt"
	"io"
	"io/fs"
	"net/http"

	"github.com/DataDog/mmh3"
	"github.com/labstack/echo/v4"
	"github.com/xvello/weblock/utils"
)

//go:embed assets
var assetFiles embed.FS

/* wrappedAssets handles serving static assets. It:
 *  - wraps the fs.Open call to forbid access to folders
 *  - pre-computes a hash of the assets to use for cache management
 *  - exposes a echo handler
 */
type wrappedAssets struct {
	root   fs.FS
	isDir  map[string]bool
	hash   string
	eTag   string
	server http.Handler
}

func loadAssets() *wrappedAssets {
	hash := computeAssetsHash()
	assets := &wrappedAssets{
		root:  assetFiles,
		isDir: make(map[string]bool),
		hash:  hash,
		eTag:  fmt.Sprintf("\"%s\"", hash),
	}
	assets.server = http.FileServer(http.FS(assets))

	return assets
}

func (w *wrappedAssets) serve(c echo.Context) error {
	if c.Request().Header.Get("If-None-Match") == w.eTag {
		return c.NoContent(http.StatusNotModified)
	}
	c.Response().Before(func() {
		c.Response().Header().Set("Vary", "Accept-Encoding")
		c.Response().Header().Set("Cache-Control", "public, max-age=86400")
		c.Response().Header().Set("ETag", w.eTag)
	})
	w.server.ServeHTTP(c.Response(), c.Request())
	return nil
}

// Open implements http.Filesystem while denying access to directory listings
func (w wrappedAssets) Open(name string) (fs.File, error) {
	isDir, found := w.isDir[name]
	if isDir {
		return nil, fs.ErrNotExist
	}
	file, err := w.root.Open(name)
	if err != nil {
		return nil, err
	}
	if !found {
		stat, err := file.Stat()
		if err != nil {
			return nil, err
		}
		w.isDir[name] = stat.IsDir()
		if stat.IsDir() {
			return nil, fs.ErrNotExist
		}
	}
	return file, nil
}

// computeAssetsHash walks the assets filesystem to compute a hash of all files.
func computeAssetsHash() string {
	hash := &mmh3.HashWriter128{}
	err := utils.Walk(assetFiles, "", func(name string, reader io.Reader) error {
		hash.AddString(name)
		_, e := io.Copy(hash, reader)
		return e
	})
	if err != nil {
		return ""
	}
	return base64.RawURLEncoding.EncodeToString(hash.Sum128().Bytes())
}
