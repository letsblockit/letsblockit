package server

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/xvello/weblock/filters"
)

func SetupRouter() (*echo.Echo, error) {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Pre(middleware.RemoveTrailingSlash())

	rawAssets, err := openAssets()
	if err != nil {
		return nil, err
	}
	pages, err := loadTemplates(buildHelpers(e))
	if err != nil {
		return nil, err
	}
	f, err := filters.LoadFilters()
	if err != nil {
		return nil, err
	}

	e.GET("/assets/*", echo.WrapHandler(http.FileServer(http.FS(rawAssets))))

	e.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "pong")
	})

	e.GET("/filters", func(c echo.Context) error {
		return pages.render(c, "list_filters", f.GetFilters())
	}).Name = "list_filters"

	e.GET("/filters/:name", func(c echo.Context) error {
		if filter, err := f.GetFilter(c.Param("name")); err == nil {
			return pages.render(c, "view_filter", filter)
		} else {
			return echo.NewHTTPError(http.StatusNotFound)
		}
	}).Name = "view_filter"
	return e, nil
}
