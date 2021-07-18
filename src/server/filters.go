package server

import (
	"fmt"
	"net/http"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/labstack/echo/v4"
	"github.com/xvello/letsblockit/src/filters"
)

func (s *Server) viewFilter(c echo.Context) error {
	filter, err := s.filters.GetFilter(c.Param("name"))
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound)
	}
	hc := buildHandlebarsContext(c, fmt.Sprintf("How to %s with uBlock or Adblock", lowerFirst(filter.Title)))
	hc["filter"] = filter

	// Parse filters param and render output if non empty
	params, err := parseFilterParams(c, filter)
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
	params, err := parseFilterParams(c, filter)
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
	hc := map[string]interface{}{
		"_naked":   true,
		"rendered": buf.String(),
	}
	return s.pages.render(c, "view-filter-render", hc)
}

func parseFilterParams(c echo.Context, filter *filters.Filter) (map[string]interface{}, error) {
	formParams, err := c.FormParams()
	if err != nil {
		return nil, err
	}
	if len(formParams) == 0 {
		return nil, nil
	}
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
			return nil, echo.NewHTTPError(http.StatusInternalServerError, "unknown param type "+p.Type)
		}
	}
	return params, err
}

func lowerFirst(s string) string {
	if s == "" {
		return ""
	}
	r, n := utf8.DecodeRuneInString(s)
	return string(unicode.ToLower(r)) + s[n:]
}
