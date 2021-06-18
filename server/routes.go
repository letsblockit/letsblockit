package server

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/xvello/weblock/filters"
)

func SetupRouter() (*echo.Echo, error) {
	rawAssets, err := openAssets()
	if err != nil {
		return nil, err
	}
	pages, err := loadTemplates()
	if err != nil {
		return nil, err
	}
	f, err := filters.LoadFilters()
	if err != nil {
		return nil, err
	}

	e := echo.New()
	e.GET("/assets/*", echo.WrapHandler(http.FileServer(http.FS(rawAssets))))
	e.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "pong")
	})
	e.GET("/filter/:name", func(c echo.Context) error {
		if filter, err := f.GetFilter(c.Param("name")); err == nil {
			return pages.render(c, "view_filter", filter)
		} else {
			return echo.NewHTTPError(http.StatusNotFound)
		}
	})

	return e, nil
}
