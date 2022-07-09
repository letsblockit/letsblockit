package server

import (
	"io"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/letsblockit/letsblockit/src/filters"
	"github.com/letsblockit/letsblockit/src/news"
	"github.com/letsblockit/letsblockit/src/pages"
)

type PageRenderer interface {
	RegisterHelpers(helpers map[string]interface{})
	RegisterContextBuilder(b pages.ContextBuilder)
	BuildPageContext(c echo.Context, title string) *pages.Context
	Render(c echo.Context, name string, data *pages.Context) error
	RenderWithSidebar(c echo.Context, name, sidebar string, data *pages.Context) error
	RedirectToPage(c echo.Context, name string, params ...interface{}) error
	Redirect(c echo.Context, code int, target string) error
}

type FilterRepository interface {
	GetFilter(name string) (*filters.Filter, error)
	GetFilters() []*filters.Filter
	GetTags() []string
	Render(w io.Writer, name string, data map[string]interface{}) error
}

type ReleaseClient interface {
	GetReleases() ([]*news.Release, error)
	GetLatestAt() (time.Time, error)
}
