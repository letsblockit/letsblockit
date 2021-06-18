package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetupRouter() (*gin.Engine, error) {
	rawAssets, err := openAssets()
	if err != nil {
		return nil, err
	}
	r := gin.Default()
	r.StaticFS("/assets", http.FS(rawAssets))
	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})
	return r, nil
}
