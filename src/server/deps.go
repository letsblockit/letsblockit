package server

import (
	"io"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/xvello/letsblockit/src/db"
	"github.com/xvello/letsblockit/src/filters"
	"github.com/xvello/letsblockit/src/news"
	"github.com/xvello/letsblockit/src/pages"
)

type PageRenderer interface {
	RegisterHelpers(helpers map[string]interface{})
	Render(c echo.Context, name string, data *pages.Context) error
	RenderWithSidebar(c echo.Context, name, sidebar string, data *pages.Context) error
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
type UserPreferenceManager interface {
	Get(c echo.Context, user uuid.UUID) (*db.UserPreference, error)
	UpdateNewsCursor(c echo.Context, user uuid.UUID, at time.Time) error
}
