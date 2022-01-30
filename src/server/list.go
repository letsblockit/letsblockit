package server

import (
	"context"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/xvello/letsblockit/src/db"
	"github.com/xvello/letsblockit/src/filters"
)

func (s *Server) renderList(c echo.Context) error {
	token, err := uuid.Parse(c.Param("token"))
	if err != nil {
		return err
	}

	var storedInstances []db.GetInstancesForListRow
	if err := s.store.RunTx(c, func(ctx context.Context, q db.Querier) error {
		storedList, e := q.GetListForToken(ctx, token)
		if e == db.NotFound {
			return echo.ErrNotFound
		} else if e != nil {
			return e
		} else if s.isUserBanned(storedList.UserID) {
			return echo.ErrForbidden
		}

		if c.Request().Header.Get("Referer") == "" {
			e = q.MarkListDownloaded(ctx, storedList.ID)
			if e != nil {
				return e
			}
		}

		storedInstances, e = q.GetInstancesForList(ctx, storedList.ID)
		return e
	}); err != nil {
		return err
	}

	list := &filters.List{Title: "My filters"}
	var customFilterInstances []*filters.Instance
	for _, storedInstance := range storedInstances {
		instance := &filters.Instance{
			Filter: storedInstance.FilterName,
			Params: make(map[string]interface{}),
		}
		err = storedInstance.Params.AssignTo(&instance.Params)
		if err != nil {
			return err
		}
		if instance.Filter == filters.CustomRulesFilterName {
			customFilterInstances = append(customFilterInstances, instance)
		} else {
			list.Instances = append(list.Instances, instance)
		}
	}
	if len(customFilterInstances) > 0 {
		list.Instances = append(list.Instances, customFilterInstances...)
	}

	return list.Render(c.Response(), c.Logger(), s.filters)
}
