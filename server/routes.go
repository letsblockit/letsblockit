package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xvello/weblock/filters"
)

func SetupRouter() (*gin.Engine, error) {
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

	r := gin.Default()
	r.StaticFS("/assets", http.FS(rawAssets))
	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})
	r.GET("/filter/:name", func(c *gin.Context) {
		if filter, err := f.GetFilter(c.Param("name")); err == nil {
			pages.render(c, "view_filter", filter)
		} else {
			c.AbortWithStatus(http.StatusNotFound)
		}
	})
	return r, nil
}
