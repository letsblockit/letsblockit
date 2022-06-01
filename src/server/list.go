package server

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/letsblockit/letsblockit/src/db"
	"github.com/letsblockit/letsblockit/src/filters"
	"gopkg.in/yaml.v2"
)

const listExportTemplate = `# letsblock.it filter list export
#
# List token: %s
# Export date: %s
#
# You can edit this file and render it locally, check out instructions at:
# https://github.com/letsblockit/letsblockit/tree/main/cmd/render/README.md

`

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

	list, err := convertFilterList(storedInstances)
	if err != nil {
		return err
	}
	return list.Render(c.Response(), c.Logger(), s.filters)
}

func (s *Server) exportList(c echo.Context) error {
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
		} else if getUser(c).Id() != storedList.UserID {
			return echo.ErrForbidden
		}
		storedInstances, e = q.GetInstancesForList(ctx, storedList.ID)
		return e
	}); err != nil {
		return err
	}

	list, err := convertFilterList(storedInstances)
	if err != nil {
		return err
	}

	c.Response().Header().Set("Content-Type", "text/yaml")
	c.Response().Header().Set("Content-Disposition", "attachment; filename=\"exported-filter-list.yaml\"")
	c.Response().WriteHeader(200)
	_, err = fmt.Fprintf(c.Response(), listExportTemplate, token, s.now().Format("2006-01-02"))
	if err != nil {
		return nil
	}
	err = yaml.NewEncoder(c.Response()).Encode(&list)
	if err != nil {
		return nil
	}
	return nil
}

func convertFilterList(storedInstances []db.GetInstancesForListRow) (*filters.List, error) {
	list := &filters.List{Title: "My filters"}
	var customFilterInstances []*filters.Instance
	for _, storedInstance := range storedInstances {
		instance := &filters.Instance{
			Filter: storedInstance.FilterName,
			Params: make(map[string]interface{}),
		}
		err := storedInstance.Params.AssignTo(&instance.Params)
		if err != nil {
			return nil, err
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
	return list, nil
}
