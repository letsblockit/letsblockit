package server

import (
	"fmt"
	"net/http"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/labstack/echo/v4"
	"github.com/xvello/letsblockit/src/filters"
	"github.com/xvello/letsblockit/src/models"
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
	params, save, err := parseFilterParams(c, filter)
	if err != nil {
		return err
	}

	// Save filter params if requested
	if save && u.IsVerified() {
		f := &models.FilterInstance{
			UserID:     u.Id(),
			FilterName: filter.Name,
		}
		s.gorm.Where(f).First(f)
		f.Params = params
		if f.FilterListID == 0 {
			list := s.getOrCreateFilterList(u)
			f.FilterListID = list.ID
			s.gorm.Create(&f)
		} else {
			s.gorm.Save(&f)
		}
		hc["saved_ok"] = true
	}

	// If no params are passed, source from the user's filters
	if !save && params == nil && u.IsVerified() {
		f := &models.FilterInstance{
			UserID:     u.Id(),
			FilterName: filter.Name,
		}
		s.gorm.Where(f).First(f)
		if f.ID > 0 {
			params = f.Params
			hc["has_instance"] = true
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
	params, _, err := parseFilterParams(c, filter)
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

func parseFilterParams(c echo.Context, filter *filters.Filter) (map[string]interface{}, bool, error) {
	formParams, err := c.FormParams()
	if err != nil {
		return nil, false, err
	}
	if len(formParams) == 0 {
		return nil, false, nil
	}

	_, save := formParams["__save"]
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
			return nil, false, echo.NewHTTPError(http.StatusInternalServerError, "unknown param type "+p.Type)
		}
	}
	return params, save, err
}

func lowerFirst(s string) string {
	if s == "" {
		return ""
	}
	r, n := utf8.DecodeRuneInString(s)
	return string(unicode.ToLower(r)) + s[n:]
}
