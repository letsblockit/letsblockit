package server

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/xvello/weblock/filters"
)

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

	s.echo.GET("/filters", func(c echo.Context) error {
		return s.pages.render(c, "list-filters", map[string]interface{} {
			"title": "Available uBlock filter templates",
			"filters": s.repo.GetFilters(),
		})
	}).Name = "list-filters"

	s.echo.GET("/filters/:name", func(c echo.Context) error {
		if filter, err := s.repo.GetFilter(c.Param("name")); err == nil {
			return s.pages.render(c, "view-filter", map[string]interface{} {
				"filter": filter,
			})
		} else {
			return echo.NewHTTPError(http.StatusNotFound)
		}
	}).Name = "view-filter"
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
