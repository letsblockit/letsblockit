package server

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/xvello/letsblockit/src/db"
	"github.com/xvello/letsblockit/src/filters"
)

const listHeaderTemplate = `! Title: letsblock.it - My filters
! Expires: 1 hour
! Homepage: https://letsblock.it
! License: https://github.com/xvello/letsblockit/blob/main/LICENSE.txt
`

const filterHeaderTemplate = `
! %s
`

func (s *Server) renderList(c echo.Context) error {
	token, err := uuid.Parse(c.Param("token"))
	if err != nil {
		return err
	}

	var instances []db.GetInstancesForListRow
	if err := s.store.RunTx(c, func(ctx context.Context, q db.Querier) error {
		list, err := q.GetListForToken(ctx, token)
		if err == db.NotFound {
			return echo.ErrNotFound
		} else if err != nil {
			return err
		}

		if c.Request().Header.Get("Referer") == "" {
			err = q.MarkListDownloaded(ctx, list.ID)
			if err != nil {
				return err
			}
		}

		instances, err = q.GetInstancesForList(ctx, list.ID)
		return err
	}); err != nil {
		return err
	}

	_, err = fmt.Fprint(c.Response(), listHeaderTemplate)
	if err != nil {
		return err
	}

	var custom []db.GetInstancesForListRow
	printFilter := func(f *db.GetInstancesForListRow) error {
		_, e := fmt.Fprintf(c.Response(), filterHeaderTemplate, f.FilterName)
		if e != nil {
			return e
		}
		params := make(map[string]interface{})
		if e = f.Params.AssignTo(&params); e != nil {
			return e
		}
		return s.filters.Render(c.Response(), f.FilterName, params)
	}
	for _, f := range instances {
		if f.FilterName == filters.CustomRulesFilterName {
			custom = append(custom, f)
			continue
		}
		if err := printFilter(&f); err != nil {
			c.Logger().Warnf("skipping filter %s in list %s: %s", f.FilterName, token, err.Error())
		}
	}
	for _, f := range custom {
		if err := printFilter(&f); err != nil {
			c.Logger().Warnf("skipping filter %s in list %s: %s", f.FilterName, token, err.Error())
		}
	}
	return nil
}
