package server

import (
	"io"

	"github.com/labstack/echo/v4"
	"github.com/xvello/letsblockit/src/filters"
	"github.com/xvello/letsblockit/src/pages"
	"github.com/xvello/letsblockit/src/store"
)

type PageRenderer interface {
	RegisterHelpers(helpers map[string]interface{})
	Render(c echo.Context, name string, data *pages.Context) error
}

type FilterRepository interface {
	GetFilter(name string) (*filters.Filter, error)
	GetFilters() []*filters.Filter
	GetTags() []string
	Render(w io.Writer, name string, data map[string]interface{}) error
}

type DataStore interface {
	CountFilters(user string) (int64, error)
	GetActiveFilterNames(user string) map[string]bool
	GetOrCreateFilterList(user string) (*store.FilterList, error)
	GetListForToken(token string) (*store.FilterList, error)
	UpsertFilterInstance(user, filterName string, params store.JSONMap) error
	DropFilterInstance(user string, filterName string) error
	GetFilterInstance(user string, filterName string) (*store.FilterInstance, error)
}
