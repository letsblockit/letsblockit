package server

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/xvello/letsblockit/src/filters"
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

func (s *Server) filterList(c echo.Context) error {
	// Check user is authenticated
	user := getUser(c)
	if user == nil {
		// Cannot find user session, redirect to login
		return s.redirect(c, "user-login")
	}
	if !user.IsVerified() {
		return s.redirect(c, "user-account")
	}

	// Handle deletion if requested, remove all instances matching a given name
	if name := c.Request().PostFormValue("delete"); name != "" {
		target := &models.FilterInstance{
			UserID:     user.Id(),
			FilterName: name,
		}
		s.gorm.Where(target).Delete(target)
	}

	// List filters and render list
	var active []*filters.Filter
	for name := range s.getActiveFilterNames(user) {
		if f, err := s.filters.GetFilter(name); err == nil {
			active = append(active, f)
		}
	}

	hc := s.buildHandlebarsContext(c, "My filters")
	hc["filters"] = &active
	return s.pages.render(c, "view-list", hc)
}
