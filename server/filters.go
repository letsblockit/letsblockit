package server

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/xvello/weblock/filters"
)

func (s *Server) viewFilter(c echo.Context) error {
	filter, err := s.filters.GetFilter(c.Param("name"))
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound)
	}
	hc := buildHandlebarsContext(c, fmt.Sprintf("How to %s with uBlock or Adblock", filter.Title))
	hc["filter"] = filter

	// Parse filters param and render output if non empty
	params, err := parseFilterParams(c, filter)
	if err != nil {
		return err
	}
	if params != nil || len(filter.Params) == 0 {
		var buf strings.Builder
		if err = filter.Render(&buf, params); err != nil {
			return err
		}
		hc["rendered"] = buf.String()
		hc["params"] = params
	} else {
		defaultParams := make(map[string]interface{})
		for _, p := range filter.Params {
			defaultParams[p.Name] = p.Default
		}
		hc["params"] = defaultParams
	}
	return s.pages.render(c, "view-filter", hc)
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
