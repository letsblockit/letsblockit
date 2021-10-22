package server

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/xvello/letsblockit/src/models"
)

const listHeaderTemplate = `! Title: letsblock.it - %s
! Expires: 1 hour
! Homepage: https://letsblock.it
! License: https://github.com/xvello/letsblockit/blob/main/LICENSE.txt
`

const filterHeaderTemplate = `
! %s
`

func (s *Server) renderList(c echo.Context) error {
	token := c.Param("token")
	list := models.FilterList{
		Token: token,
	}
	s.gorm.Where(&list).Preload("FilterInstances").First(&list)
	if list.ID == 0 {
		return echo.ErrNotFound
	}

	_, err := fmt.Fprintf(c.Response(), listHeaderTemplate, list.Name)
	if err != nil {
		return err
	}

	for _, f := range list.FilterInstances {
		_, err := fmt.Fprintf(c.Response(), filterHeaderTemplate, f.FilterName)
		if err != nil {
			return err
		}
		err = s.filters.Render(c.Request().Context(), c.Response(), f.FilterName, f.Params)
		if err != nil {
			c.Logger().Warnf("skipping filter %s in list %s: %s", f.FilterName, token, err.Error())
		}
	}
	return nil
}
