package server

import (
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/xvello/weblock/filters"
)

type navigationLink struct {
	Name   string
	Target string
}

var navigationLinks = []navigationLink{{
	Name:   "Filter list",
	Target: "filters",
}, {
	Name:   "About",
	Target: "about",
}}

type Server struct {
	echo      *echo.Echo
	pages     *templates
	filters   *filters.Repository
	assetETag string
}

func NewServer() *Server {
	return &Server{echo: echo.New()}
}

func (s *Server) Start() error {
	concurrentRunOrPanic([]func([]error){
		func(errs []error) {
			assetHash := computeAssetsHash()
			s.assetETag = fmt.Sprintf("\"%s\"", assetHash)
			s.pages, errs[0] = loadTemplates(buildHelpers(s.echo, assetHash))
		},
		func(errs []error) { s.filters, errs[0] = filters.LoadFilters() },
		func(_ []error) { s.setupRouter() },
	})
	return s.echo.Start(":8080")
}

func (s *Server) setupRouter() {
	s.echo.Use(middleware.Logger())
	s.echo.Pre(middleware.RemoveTrailingSlash())

	s.echo.GET("/", func(c echo.Context) error {
		return c.Redirect(http.StatusFound, "/filters")
	}).Name = "index"

	assetsServer := http.FileServer(http.FS(openAssets()))
	s.echo.GET("/assets/*", func(c echo.Context) error {
		if c.Request().Header.Get("If-None-Match") == s.assetETag {
			return c.NoContent(http.StatusNotModified)
		}
		c.Response().Before(func() {
			c.Response().Header().Set("Vary", "Accept-Encoding")
			c.Response().Header().Set("Cache-Control", "public, max-age=86400")
			c.Response().Header().Set("ETag", s.assetETag)
		})
		assetsServer.ServeHTTP(c.Response(), c.Request())
		return nil
	})

	s.echo.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "pong")
	})

	s.echo.GET("/about", func(c echo.Context) error {
		return s.pages.render(c, "about", buildHandlebarsContext(c, "About weBlock"))
	}).Name = "list-filters"

	s.echo.GET("/filters", func(c echo.Context) error {
		hc := buildHandlebarsContext(c, "Available uBlock filter templates")
		hc["filters"] = s.filters.GetFilters()
		return s.pages.render(c, "list-filters", hc)
	}).Name = "list-filters"

	s.echo.GET("/filters/:name", s.viewFilter).Name = "view-filter"
	s.echo.POST("/filters/:name", s.viewFilter)
}

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
	if params != nil {
		hc["rendered"], err = filter.Render(params)
		if err != nil {
			return err
		}
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

func buildHandlebarsContext(c echo.Context, title string) map[string]interface{} {
	var section string
	for _, s := range strings.Split(c.Path(), "/") {
		if s != "" {
			section = s
			break
		}
	}
	return map[string]interface{}{
		"navLinks":   navigationLinks,
		"navCurrent": section,
		"title":      title,
	}
}

func concurrentRunOrPanic(tasks []func([]error)) {
	var wg sync.WaitGroup
	for i, f := range tasks {
		wg.Add(1)
		f := f
		i := i
		go func(wg *sync.WaitGroup) {
			start := time.Now()
			defer wg.Done()
			errs := []error{nil}
			f(errs)
			if errs[0] != nil {
				panic(errs[0])
			}
			fmt.Printf("Task %d took %s\n", i, time.Since(start))
		}(&wg)
	}
	wg.Wait()
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
