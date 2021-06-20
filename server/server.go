package server

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/xvello/weblock/filters"
)

var DryRunFinished = errors.New("dry run finished")

type Options struct {
	DryRun  bool   `arg:"--dry-run" help:"instantiate all components and exit"`
	Address string `default:"127.0.0.1:8765" help:"address to listen to"`
}

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
	options   *Options
	echo      *echo.Echo
	pages     *pages
	filters   *filters.Repository
	assetHash string
	assetETag string
}

func NewServer(options *Options) *Server {
	return &Server{
		options: options,
		echo:    echo.New(),
	}
}

func (s *Server) Start() error {
	concurrentRunOrPanic([]func([]error){
		func(_ []error) {
			s.assetHash = computeAssetsHash()
			s.assetETag = fmt.Sprintf("\"%s\"", s.assetHash)
		},
		func(errs []error) { s.pages, errs[0] = loadTemplates() },
		func(errs []error) { s.filters, errs[0] = filters.LoadFilters() },
		func(_ []error) { s.setupRouter() },
	})
	s.pages.registerHelpers(buildHelpers(s.echo, s.assetHash))
	if s.options.DryRun {
		return DryRunFinished
	}
	return s.echo.Start(s.options.Address)
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

	s.addStatic("/about", "about", "About the weBlock project")

	s.echo.GET("/filters", func(c echo.Context) error {
		hc := buildHandlebarsContext(c, "Available uBlock filter templates")
		hc["filters"] = s.filters.GetFilters()
		return s.pages.render(c, "list-filters", hc)
	}).Name = "list-filters"

	s.echo.GET("/filters/:name", s.viewFilter).Name = "view-filter"
	s.echo.POST("/filters/:name", s.viewFilter)
}

func (s *Server) addStatic(url, page, title string) {
	s.echo.GET(url, func(c echo.Context) error {
		return s.pages.render(c, page, buildHandlebarsContext(c, title))
	}).Name = page
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
	var output strings.Builder
	output.WriteString("Startup tasks finished |")
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
			output.WriteString(fmt.Sprintf(" %d: %s |", i, time.Since(start)))
		}(&wg)
	}
	wg.Wait()
	fmt.Println(output.String())
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
