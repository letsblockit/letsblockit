package server

import (
	"fmt"
	"net/http"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/labstack/echo/v4"
	"github.com/xvello/letsblockit/src/filters"
	"github.com/xvello/letsblockit/src/store"
)

func (s *Server) viewFilter(c echo.Context) error {
	filter, err := s.filters.GetFilter(c.Param("name"))
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound)
	}
	hc := s.buildHandlebarsContext(c, fmt.Sprintf("How to %s with uBlock or Adblock", lowerFirst(filter.Title)))
	hc["filter"] = filter
	u := getUser(c)

	// Parse filters param and render output if non empty
	params, save, disable, err := parseFilterParams(c, filter)
	if err != nil {
		return err
	}

	// Save filter params if requested
	if save && u.IsVerified() {
		if err = s.store.UpsertFilterInstance(u.Id(), filter.Name, params); err != nil {
			return err
		}
		hc["saved_ok"] = true
		hc["has_instance"] = true
	}

	// Handle deletion if requested, remove all instances matching a given name
	if disable && u.IsVerified() {
		if err = s.store.DropFilterInstance(u.Id(), filter.Name); err != nil {
			return err
		}
		return s.redirect(c, "list-filters")
	}

	// If no params are passed, source from the user's filters
	if !save && params == nil && u.IsVerified() {
		f, err := s.store.GetFilterInstance(u.Id(), filter.Name)
		switch err {
		case nil:
			params = f.Params
			hc["has_instance"] = true
		case store.ErrRecordNotFound: // ok
		default:
			return err
		}
	}

	// If no config found, inject the default ones
	if params == nil {
		params = make(map[string]interface{})
		for _, p := range filter.Params {
			params[p.Name] = p.Default
		}
	}

	// Render the filter template
	var buf strings.Builder
	if err = s.filters.Render(c.Request().Context(), &buf, filter.Name, params); err != nil {
		return err
	}
	hc["rendered"] = buf.String()
	hc["params"] = params

	return s.pages.render(c, "view-filter", hc)
}

func (s *Server) viewFilterRender(c echo.Context) error {
	filter, err := s.filters.GetFilter(c.Param("name"))
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound)
	}

	// Parse filters param and render output if non empty
	params, _, _, err := parseFilterParams(c, filter)
	if err != nil {
		return err
	}

	// If no params are passed, inject the default ones
	if params == nil {
		params = make(map[string]interface{})
		for _, p := range filter.Params {
			params[p.Name] = p.Default
		}
	}

	// Render the filter template
	var buf strings.Builder
	if err = s.filters.Render(c.Request().Context(), &buf, filter.Name, params); err != nil {
		return err
	}
	hc := s.buildHandlebarsContext(c, "")
	hc["_naked"] = true
	hc["rendered"] = buf.String()

	return s.pages.render(c, "view-filter-render", hc)
}

func parseFilterParams(c echo.Context, filter *filters.Filter) (map[string]interface{}, bool, bool, error) {
	formParams, err := c.FormParams()
	if err != nil {
		return nil, false, false, err
	}
	if len(formParams) == 0 {
		return nil, false, false, nil
	}

	_, save := formParams["__save"]
	_, disable := formParams["__disable"]
	params := make(map[string]interface{})
	for _, p := range filter.Params {
		switch p.Type {
		case filters.StringListParam:
			var values []string
			for _, v := range formParams[p.Name] {
				if v != "" {
					values = append(values, v)
				}
			}
			params[p.Name] = values
		case filters.StringParam:
			params[p.Name] = formParams.Get(p.Name)
		case filters.BooleanParam:
			params[p.Name] = formParams.Get(p.Name) == "on"
		default:
			return nil, false, false, echo.NewHTTPError(http.StatusInternalServerError, "unknown param type "+p.Type)
		}
	}
	return params, save, disable, err
}

func lowerFirst(s string) string {
	if s == "" {
		return ""
	}
	r, n := utf8.DecodeRuneInString(s)
	return string(unicode.ToLower(r)) + s[n:]
}
