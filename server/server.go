package server

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/honeycombio/beeline-go"
	"github.com/honeycombio/beeline-go/wrappers/hnyecho"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/xvello/weblock/filters"
)

var DryRunFinished = errors.New("dry run finished")

type Options struct {
	DryRun         bool   `arg:"--dry-run" help:"instantiate all components and exit"`
	Address        string `default:"127.0.0.1:8765" help:"address to listen to"`
	HoneycombKey   string `arg:"--honeycomb-key" help:"write key for honeycomb.io trace submission"`
	HoneycombDebug bool   `arg:"--honeycomb-debug" help:"write traces to stdout instead"`
}

var navigationLinks = []struct {
	Name   string
	Target string
}{{
	Name:   "Filter list",
	Target: "filters",
}, {
	Name:   "About",
	Target: "about",
}}

type Server struct {
	options *Options
	echo    *echo.Echo
	pages   *pages
	filters *filters.Repository
	assets  *wrappedAssets
}

func NewServer(options *Options) *Server {
	return &Server{
		options: options,
		echo:    echo.New(),
	}
}

func (s *Server) Start() error {
	if s.options.HoneycombKey != "" || s.options.HoneycombDebug {
		beeline.Init(beeline.Config{
			WriteKey: s.options.HoneycombKey,
			Dataset: "weblock",
			STDOUT:   s.options.HoneycombDebug,
		})
		s.echo.Use(hnyecho.New().Middleware())
		defer beeline.Close()
	}
	concurrentRunOrPanic([]func([]error){
		func(_ []error) { s.assets = loadAssets() },
		func(errs []error) { s.pages, errs[0] = loadPages() },
		func(errs []error) { s.filters, errs[0] = filters.LoadFilters() },
	})

	s.pages.registerHelpers(buildHelpers(s.echo, s.assets.hash))
	s.setupRouter()
	if s.options.DryRun {
		return DryRunFinished
	}
	return s.echo.Start(s.options.Address)
}

func (s *Server) setupRouter() {
	s.echo.Use(middleware.Logger())
	s.echo.Pre(middleware.RemoveTrailingSlash())

	s.echo.GET("/assets/*", s.assets.serve)

	s.echo.GET("/", func(c echo.Context) error {
		return c.Redirect(http.StatusFound, "/filters")
	}).Name = "index"

	s.addStatic("/about", "about", "About the weBlock project")

	s.echo.GET("/filters", func(c echo.Context) error {
		hc := buildHandlebarsContext(c, "Available uBlock filter templates")
		hc["filters"] = s.filters.GetFilters()
		return s.pages.render(c, "list-filters", hc)
	}).Name = "list-filters"

	s.echo.GET("/filters/tag/:tag", func(c echo.Context) error {
		tag := c.Param("tag")
		hc := buildHandlebarsContext(c, "Filter templates for "+tag)
		var matching []*filters.Filter
		for _, f := range s.filters.GetFilters() {
			for _, t := range f.Tags {
				if t == tag {
					matching = append(matching, f)
					break
				}
			}
		}
		hc["filters"] = matching
		// TODO: link to go back to all tags
		return s.pages.render(c, "list-filters", hc)
	}).Name = "filters-for-tag"

	s.echo.GET("/filters/:name", s.viewFilter).Name = "view-filter"
	s.echo.POST("/filters/:name", s.viewFilter)
}

func (s *Server) addStatic(url, page, title string) {
	s.echo.GET(url, func(c echo.Context) error {
		return s.pages.render(c, page, buildHandlebarsContext(c, title))
	}).Name = page
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
