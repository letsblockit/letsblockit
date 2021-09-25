package server

import (
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/DataDog/datadog-go/statsd"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/xvello/letsblockit/src/filters"
)

var ErrDryRunFinished = errors.New("dry run finished")

type Options struct {
	DryRun     bool   `arg:"--dry-run" help:"instantiate all components and exit"`
	Address    string `default:"127.0.0.1:8765" help:"address to listen to"`
	Statsd     string `help:"address to send statsd metrics to"`
	OryProject string `help:"oxy cloud project to check credentials against"`
	Reload     bool   `help:"reload frontend when the backend restarts"`
	Debug      bool   `help:"log with debug level"`
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
	concurrentRunOrPanic([]func([]error){
		func(_ []error) { s.assets = loadAssets() },
		func(errs []error) { s.pages, errs[0] = loadPages() },
		func(errs []error) { s.filters, errs[0] = filters.LoadFilters() },
	})

	if s.options.Statsd != "" {
		dsd, err := statsd.New(s.options.Statsd)
		if err != nil {
			return err
		}
		s.echo.Use(buildDogstatsMiddleware(dsd))
	}

	if s.options.OryProject != "" {
		s.echo.Use(buildOryMiddleware(s.options.OryProject, s.echo.Logger))
	}

	if s.options.Debug {
		s.echo.Logger.SetLevel(log.DEBUG)
	}

	s.pages.registerHelpers(buildHelpers(s.echo, s.assets.hash))
	s.setupRouter()
	if s.options.DryRun {
		return ErrDryRunFinished
	}
	return s.echo.Start(s.options.Address)
}

func (s *Server) setupRouter() {
	s.echo.Use(middleware.Logger())
	s.echo.Pre(middleware.RemoveTrailingSlash())
	s.echo.Pre(middleware.Rewrite(map[string]string{
		"/favicon.ico": "/assets/images/favicon.ico",
		"/":            "/filters",
	}))

	if s.options.Reload {
		s.echo.GET("/should-reload", func(c echo.Context) error {
			// Set the headers related to event streaming.
			c.Response().Header().Set("Content-Type", "text/event-stream")
			c.Response().Header().Set("Cache-Control", "no-cache")
			c.Response().Header().Set("Connection", "keep-alive")
			c.Response().Header().Set("Transfer-Encoding", "chunked")
			if _, err := fmt.Fprintln(c.Response(), "retry:1000"); err != nil {
				return nil
			}
			c.Response().Flush()

			// Block indefinitely to keep the SSE open
			<-c.Request().Context().Done()
			return nil
		})
	}

	s.echo.GET("/assets/*", s.assets.serve)

	s.addStatic("/about", "about", "About: Let’s block it!")

	s.echo.GET("/filters", func(c echo.Context) error {
		hc := s.buildHandlebarsContext(c, "Available uBlock filter templates")
		hc["filters"] = s.filters.GetFilters()
		return s.pages.render(c, "list-filters", hc)
	}).Name = "list-filters"

	s.echo.GET("/filters/tag/:tag", func(c echo.Context) error {
		tag := c.Param("tag")
		hc := s.buildHandlebarsContext(c, "Filter templates for "+tag)
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
	s.echo.POST("/filters/:name/render", s.viewFilterRender).Name = "view-filter-render"

	s.echo.GET("/user/login", s.userLogin).Name = "user-login"
	s.echo.GET("/user/account", s.userAccount).Name = "user-account"
}

func (s *Server) addStatic(url, page, title string) {
	s.echo.GET(url, func(c echo.Context) error {
		return s.pages.render(c, page, s.buildHandlebarsContext(c, title))
	}).Name = page
}

func (s *Server) buildHandlebarsContext(c echo.Context, title string) map[string]interface{} {
	var section string
	for _, s := range strings.Split(c.Path(), "/") {
		if s != "" {
			section = s
			break
		}
	}
	context := map[string]interface{}{
		"navLinks":   navigationLinks,
		"navCurrent": section,
		"title":      title,
		"logged":     c.Get(userContextKey) != nil,
	}
	if s.options.Reload {
		context["jsImports"] = []string{"reload.js"}
	}
	return context
}

func concurrentRunOrPanic(tasks []func([]error)) {
	var wg sync.WaitGroup
	timings := make([]time.Duration, len(tasks))
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
			timings[i] = time.Since(start)
		}(&wg)
	}
	wg.Wait()

	var output strings.Builder
	output.WriteString("Startup tasks finished |")
	for i, d := range timings {
		output.WriteString(fmt.Sprintf(" %d: %s |", i, d))
	}
	fmt.Println(output.String())
}

func buildDogstatsMiddleware(dsd statsd.ClientInterface) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()
			if err := next(c); err != nil {
				c.Error(err)
			}
			duration := time.Since(start)
			_ = dsd.Distribution("letsblockit.request_duration", float64(duration.Nanoseconds()), nil, 1)
			_ = dsd.Incr("letsblockit.request_count", []string{fmt.Sprintf("status:%d", c.Response().Status)}, 1)
			return nil
		}
	}
}
