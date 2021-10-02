package server

import (
	"github.com/labstack/echo/v4"
	"github.com/xvello/letsblockit/src/models"
)

func (s *Server) renderList(c echo.Context) error {
	token := c.Param("token")
	list := models.FilterList{
		Token: token,
	}
	s.gorm.Where(&list).Preload("FilterInstances").First(&list)
	if list.ID == 0 {
		return echo.ErrNotFound
	}
	_, err := c.Response().Write([]byte("hello " + token))
	return err
}
