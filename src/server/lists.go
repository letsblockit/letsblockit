package server

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/xvello/letsblockit/src/models"
)

const listHeaderTemplate = `! Title: letsblock.it - %s
! Expires: 1 hour

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
			return err
		}
	}
	return nil
}

func (s *Server) editList(c echo.Context) error {
	user := getUser(c)
	if user == nil {
		// Cannot find user session, redirect to login
		return s.redirect(c, "user-login")
	}
	if !user.IsVerified() {
		return s.redirect(c, "user-account")
	}
	var filters models.FilterInstance
	s.gorm.Where("user_id = ?", user.Id()).Order("filter_name").Find(&filters)

	hc := s.buildHandlebarsContext(c, "My filters")
	hc["filters"] = &filters
	return s.pages.render(c, "user-account", hc)
}
