package server

import (
	"fmt"
	"net/http"
	"strconv"
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
	echo  *echo.Echo
	pages *templates
	repo  *filters.Repository
}

func NewServer() *Server {
	return &Server{echo: echo.New()}
}

func (s *Server) Start() error {
	concurrentRunOrPanic([]func([]error){
		func(errs []error) { s.pages, errs[0] = loadTemplates(buildHelpers(s.echo)) },
		func(errs []error) { s.repo, errs[0] = filters.LoadFilters() },
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

	s.echo.GET("/assets/*", echo.WrapHandler(http.FileServer(http.FS(openAssets()))))

	s.echo.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "pong")
	})

	s.echo.GET("/about", func(c echo.Context) error {
		return s.pages.render(c, "about", buildHandlebarsContext(c, "About weBlock"))
	}).Name = "list-filters"

	s.echo.GET("/filters", func(c echo.Context) error {
		hc := buildHandlebarsContext(c, "Available uBlock filter templates")
		hc["filters"] = s.repo.GetFilters()
		return s.pages.render(c, "list-filters", hc)
	}).Name = "list-filters"

	s.echo.GET("/filters/:name", s.viewFilter).Name = "view-filter"
	s.echo.POST("/filters/:name", s.viewFilter)
}

func (s *Server) viewFilter(c echo.Context) error {
	filter, err := s.repo.GetFilter(c.Param("name"))
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
			params[p.Name] = formParams[p.Name]
		case filters.StringParam:
			params[p.Name] = formParams.Get(p.Name)
		case filters.BooleanParam:
			params[p.Name], _ = strconv.ParseBool(formParams.Get(p.Name))
		default:
			return nil, echo.NewHTTPError(http.StatusInternalServerError, "unknown param type "+p.Type)
		}
	}
	return params, err
}