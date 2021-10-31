package server

import (
	"fmt"
	"sort"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/xvello/letsblockit/src/filters"
	"github.com/xvello/letsblockit/src/store"
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
	list, err := s.store.GetListForToken(token)
	switch err {
	case nil:
		// ok
	case store.ErrRecordNotFound:
		return echo.ErrNotFound
	default:
		return err
	}

	_, err = fmt.Fprintf(c.Response(), listHeaderTemplate, list.Name)
	if err != nil {
		return err
	}

	sortFilters(list.FilterInstances)
	for _, f := range list.FilterInstances {
		_, err := fmt.Fprintf(c.Response(), filterHeaderTemplate, f.FilterName)
		if err != nil {
			return err
		}
		err = s.filters.Render(c.Response(), f.FilterName, f.Params)
		if err != nil {
			c.Logger().Warnf("skipping filter %s in list %s: %s", f.FilterName, token, err.Error())
		}
	}
	return nil
}

// sortFilters sorts filters instances per filter name, moving custom-rules at the end
func sortFilters(instances []*store.FilterInstance) {
	sort.Slice(instances, func(i, j int) bool {
		if instances[i].FilterName == filters.CustomRulesFilterName {
			return false
		}
		if instances[j].FilterName == filters.CustomRulesFilterName {
			return true
		}
		return strings.Compare(instances[i].FilterName, instances[j].FilterName) < 0
	})
}
