package server

import (
	"context"
	"io"

	"github.com/labstack/echo/v4"
	"github.com/xvello/letsblockit/src/filters"
)

type pageRenderer interface {
	RegisterHelpers(helpers map[string]interface{})
	Render(c echo.Context, name string, data map[string]interface{}) error
}

type filterRepository interface {
	GetFilter(name string) (*filters.Filter, error)
	GetFilters() []*filters.Filter
	GetTags() []string
	Render(ctx context.Context, w io.Writer, name string, data map[string]interface{}) error
}
