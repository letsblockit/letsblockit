package server

import (
	"context"
	"net/http"
	"strings"

	"github.com/jackc/pgtype"
	"github.com/labstack/echo/v4"
	"github.com/letsblockit/letsblockit/src/db"
	"github.com/letsblockit/letsblockit/src/filters"
)

type filterAction uint8

const (
	actionRender = filterAction(iota)
	actionSave
	actionDelete
)

func (s *Server) listFilters(c echo.Context) error {
	tag := c.Param("tag")
	hc := s.buildPageContext(c, "Available uBlock filter templates")
	if tag != "" {
		hc.Title = "Available filter templates for " + tag
		hc.Add("tag_search", tag)
	}

	hc.Add("filter_tags", s.filters.GetTags())
	var activeNames map[string]struct{}
	if hc.UserLoggedIn {
		var updatedFilters map[string]bool
		var testingFilters map[string]bool
		activeNames = make(map[string]struct{})
		instances, _ := s.store.GetActiveFiltersForUser(c.Request().Context(), hc.UserID)
		for _, instance := range instances {
			activeNames[instance.FilterName] = struct{}{}
			if s.hasMissingParams(instance) {
				if updatedFilters == nil {
					updatedFilters = make(map[string]bool)
				}
				updatedFilters[instance.FilterName] = true
			}
			if instance.TestMode {
				if testingFilters == nil {
					testingFilters = make(map[string]bool)
				}
				testingFilters[instance.FilterName] = true
			}
		}
		if len(instances) > 0 {
			if info, err := s.store.GetListForUser(c.Request().Context(), hc.UserID); err == nil {
				hc.Add("list_token", info.Token.String())
				hc.Add("list_downloaded", info.Downloaded)
			}
		}
		if len(updatedFilters) > 0 {
			hc.Add("updated_filters", updatedFilters)
		}
		if len(testingFilters) > 0 {
			hc.Add("testing_filters", testingFilters)
		}
	}

	// Template and group filters, or quick return on homepage
	if len(activeNames) == 0 && len(tag) == 0 {
		hc.Add("available_filters", s.filters.GetAll())
	} else {
		var active, available []*filters.Template
		for _, f := range s.filters.GetAll() {
			if tag != "" {
				if !f.HasTag(tag) {
					continue
				}
			}
			if _, ok := activeNames[f.Name]; ok {
				active = append(active, f)
			} else {
				available = append(available, f)
			}
		}
		hc.Add("active_filters", active)
		hc.Add("available_filters", available)
	}
	return s.pages.Render(c, "list-filters", hc)
}

func (s *Server) viewFilter(c echo.Context) error {
	filter, err := s.filters.Get(c.Param("name"))
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound)
	}
	hc := s.buildPageContext(c, filter.Title)
	hc.Add("filter", filter)

	// Parse filters param and render output if non-empty
	instance, action, err := parseFilterParams(c, filter)
	if err != nil {
		return err
	}

	switch {
	case hc.UserLoggedIn && action == actionSave:
		// Save filter params if requested
		var out pgtype.JSONB
		if err = out.Set(instance.Params); err != nil {
			return err
		}
		if err = s.upsertFilterParams(c, hc.UserID, instance); err != nil {
			return err
		}
		hc.Add("saved_ok", true)
		hc.Add("has_instance", true)
	case hc.UserLoggedIn && action == actionDelete:
		// Handle deletion if requested, remove all instances matching a given name
		if err = s.store.DeleteInstanceForUserAndFilter(c.Request().Context(), db.DeleteInstanceForUserAndFilterParams{
			UserID:     hc.UserID,
			FilterName: filter.Name,
		}); err != nil {
			return err
		}
		return s.pages.RedirectToPage(c, "list-filters")
	case hc.UserLoggedIn:
		// If no params are passed, source from the user's filters
		stored, err := s.store.GetInstanceForUserAndFilter(c.Request().Context(), db.GetInstanceForUserAndFilterParams{
			UserID:     hc.UserID,
			FilterName: filter.Name,
		})
		switch err {
		case nil:
			hc.Add("has_instance", true)
			if instance.Params == nil {
				if err = stored.Params.AssignTo(&instance.Params); err != nil {
					return err
				}
			}
			instance.TestMode = stored.TestMode
		case db.NotFound: // ok
		default:
			return err
		}
	}

	if instance.Params == nil {
		// If no params found, inject the default values
		if len(filter.Params) > 0 {
			instance.Params = make(map[string]interface{})
			for _, param := range filter.Params {
				instance.Params[param.Name] = param.Default
				for _, preset := range param.Presets {
					instance.Params[param.BuildPresetParamName(preset.Name)] = preset.Default
				}
			}
		}
	} else {
		// Check whether new params have been added
		newParams := make(map[string]bool)
		for _, param := range filter.Params {
			if _, found := instance.Params[param.Name]; !found {
				newParams[param.Name] = true
			}
			for _, preset := range param.Presets {
				name := param.BuildPresetParamName(preset.Name)
				if _, found := instance.Params[name]; !found {
					newParams[name] = true
				}
			}
		}
		if len(newParams) > 0 {
			hc.Add("new_params", newParams)
		}
	}

	// Render the filter template
	var buf strings.Builder
	if err = s.filters.Render(&buf, instance); err != nil {
		return err
	}
	hc.Add("rendered", buf.String())
	hc.Add("params", instance.Params)
	hc.Add("test_mode", instance.TestMode)
	return s.pages.Render(c, "view-filter", hc)
}

func (s *Server) viewFilterRender(c echo.Context) error {
	filter, err := s.filters.Get(c.Param("name"))
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound)
	}

	// Parse filters param and render output if non empty
	instance, _, err := parseFilterParams(c, filter)
	if err != nil {
		return err
	}

	// If no params are passed, inject the default ones
	if instance.Params == nil {
		instance.Params = make(map[string]interface{})
		for _, p := range filter.Params {
			instance.Params[p.Name] = p.Default
		}
	}

	// Render the filter template
	var buf strings.Builder
	if err = s.filters.Render(&buf, instance); err != nil {
		return err
	}
	hc := s.buildPageContext(c, "")
	hc.NakedContent = true
	hc.Add("rendered", buf.String())

	// Detect user session from form params (this endpoint is unauthenticated)
	if formParams, err := c.FormParams(); err == nil {
		if _, ok := formParams["__logged_in"]; ok {
			hc.UserLoggedIn = true
		}
	}

	return s.pages.Render(c, "view-filter-render", hc)
}

func (s *Server) upsertFilterParams(c echo.Context, user string, instance *filters.Instance) error {
	out := pgtype.JSONB{
		Bytes:  nil,
		Status: pgtype.Null,
	}
	if len(instance.Params) > 0 {
		if err := out.Set(&instance.Params); err != nil {
			return err
		}
	}
	return s.store.RunTx(c, func(ctx context.Context, q db.Querier) error {
		count, err := q.CountInstanceForUserAndFilter(ctx, db.CountInstanceForUserAndFilterParams{
			UserID:     user,
			FilterName: instance.Template,
		})
		if err != nil {
			return err
		}
		if count == 0 {
			if listCount, _ := q.CountListsForUser(ctx, user); listCount == 0 {
				if _, err := q.CreateListForUser(ctx, user); err != nil {
					return err
				}
			}
			return q.CreateInstanceForUserAndFilter(ctx, db.CreateInstanceForUserAndFilterParams{
				UserID:     user,
				FilterName: instance.Template,
				Params:     out,
				TestMode:   instance.TestMode,
			})
		} else {
			return q.UpdateInstanceForUserAndFilter(ctx, db.UpdateInstanceForUserAndFilterParams{
				UserID:     user,
				FilterName: instance.Template,
				Params:     out,
				TestMode:   instance.TestMode,
			})
		}
	})
}

func parseFilterParams(c echo.Context, filter *filters.Template) (*filters.Instance, filterAction, error) {
	formParams, err := c.FormParams()
	if err != nil {
		return nil, actionRender, err
	}
	if len(formParams) == 0 {
		return &filters.Instance{
			Template: filter.Name,
		}, actionRender, nil
	}

	action := actionRender
	if _, ok := formParams["__save"]; ok {
		action = actionSave
	} else if _, ok := formParams["__disable"]; ok {
		action = actionDelete
	}

	instance := &filters.Instance{
		Template: filter.Name,
		Params:   make(map[string]interface{}),
		TestMode: formParams.Get("__test_mode") == "on",
	}

	for _, p := range filter.Params {
		switch p.Type {
		case filters.StringListParam:
			var values []string
			for _, v := range formParams[p.Name] {
				if v != "" {
					values = append(values, v)
				}
			}
			instance.Params[p.Name] = values
			for _, preset := range p.Presets {
				name := p.BuildPresetParamName(preset.Name)
				instance.Params[name] = formParams.Get(name) == "on"
			}
		case filters.StringParam:
			instance.Params[p.Name] = formParams.Get(p.Name)
		case filters.MultiLineParam:
			instance.Params[p.Name] = formParams.Get(p.Name)
		case filters.BooleanParam:
			instance.Params[p.Name] = formParams.Get(p.Name) == "on"
		default:
			return nil, action, echo.NewHTTPError(http.StatusInternalServerError, "unknown param type "+p.Type)
		}
	}
	return instance, action, err
}

func (s *Server) hasMissingParams(instance db.GetActiveFiltersForUserRow) bool {
	filter, err := s.filters.Get(instance.FilterName)
	if err != nil || len(filter.Params) == 0 {
		return false
	}

	params := make(map[string]interface{})
	if err := instance.Params.AssignTo(&params); err != nil {
		return true
	}
	for _, p := range filter.Params {
		if _, found := params[p.Name]; !found {
			return true
		}
		for _, preset := range p.Presets {
			if _, found := params[p.BuildPresetParamName(preset.Name)]; !found {
				return true
			}
		}
	}
	return false
}
