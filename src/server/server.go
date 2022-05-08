package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/DataDog/datadog-go/v5/statsd"
	"github.com/coreos/go-systemd/activation"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/xvello/letsblockit/src/db"
	"github.com/xvello/letsblockit/src/filters"
	"github.com/xvello/letsblockit/src/news"
	"github.com/xvello/letsblockit/src/pages"
	"github.com/xvello/letsblockit/src/users"
)

var ErrDryRunFinished = errors.New("dry run finished")

const (
	loggerFormat = `{"http":{"host":"${host}","remote_ip":"${remote_ip}",` +
		`"method":"${method}","uri":"${uri}","path":"${path}",` +
		`"user_agent":"${user_agent}","referer":"${referer}","is_htmx":"${header:HX-Request}",` +
		`"status":"${status}","error":"${error}","latency":${latency},` +
		`"bytes_in":${bytes_in},"bytes_out":${bytes_out}}}` + "\n"
	mainDomain = "letsblock.it"
	csrfLookup = "_csrf"
)

type Options struct {
	Address       string `default:"127.0.0.1:8765" help:"address to listen to"`
	SystemdSocket bool   `help:"use a systemd socket instead of opening a port"`
	Debug         bool   `help:"log with debug level"`
	DryRun        bool   `help:"instantiate all components and exit"`
	KratosURL     string `default:"http://localhost:4000/.ory" help:"url of the kratos API, defaults to using local proxy"`
	Reload        bool   `help:"reload frontend when the backend restarts"`
	Statsd        string `placeholder:"localhost:8125" help:"address to send statsd metrics to, disabled by default"`
	DatabaseName  string `default:"letsblockit" help:"psql database name to use"`
	DatabaseHost  string `default:"/var/run/postgresql" help:"psql host to connect to"`
	silent        bool
}

var navigationLinks = []struct {
	Name   string
	Target string
}{{
	Name:   "Filter list",
	Target: "filters",
}, {
	Name:   "Help",
	Target: "help",
}, {
	Name:   "About",
	Target: "help/about",
}, {
	Name:   "Contributing",
	Target: "help/contributing",
}}

type Server struct {
	assets      *wrappedAssets
	banned      map[uuid.UUID]struct{}
	echo        *echo.Echo
	filters     FilterRepository
	now         func() time.Time
	options     *Options
	pages       PageRenderer
	preferences UserPreferenceManager
	releases    ReleaseClient
	statsd      statsd.ClientInterface
	store       db.Store
}

func NewServer(options *Options) *Server {
	return &Server{
		options: options,
		echo:    echo.New(),
		now:     time.Now,
	}
}

func (s *Server) Start() error {
	concurrentRunOrPanic([]func([]error){
		func(_ []error) { s.assets = loadAssets() },
		func(errs []error) { s.pages, errs[0] = pages.LoadPages() },
		func(errs []error) { s.filters, errs[0] = filters.LoadFilters() },
		func(errs []error) {
			s.store, errs[0] = db.Connect(s.options.DatabaseHost, s.options.DatabaseName)
			if errs[0] == nil {
				errs[0] = db.Migrate(s.options.DatabaseHost, s.options.DatabaseName)
			}
			if errs[0] == nil {
				errs[0] = s.loadBannedUsers()
			}
			if errs[0] == nil {
				s.preferences, errs[0] = users.NewPreferenceManager(s.store)
			}
		},
	})

	s.releases = news.NewReleaseClient(news.GithubReleasesEndpoint)
	if s.options.Statsd != "" {
		dsd, err := statsd.New(s.options.Statsd)
		if err != nil {
			return err
		}
		s.statsd = dsd
		s.echo.Use(buildDogstatsMiddleware(dsd))
		go collectStats(s.echo.Logger, s.store, dsd)
	} else {
		s.statsd = &statsd.NoOpClient{}
	}

	s.pages.RegisterHelpers(buildHelpers(s.echo, s.assets.hash))
	s.setupRouter()
	if s.options.DryRun {
		return ErrDryRunFinished
	}

	if s.options.SystemdSocket {
		listeners, err := activation.Listeners()
		if err != nil {
			return err
		}
		if len(listeners) != 1 {
			return errors.New("unexpected number of socket activation fds")
		} else {
			fmt.Println("reusing systemd socket...")
		}
		s.echo.Listener = listeners[0]
	}
	return s.echo.Start(s.options.Address)
}

func (s *Server) setupRouter() {
	s.echo.Use(middleware.Recover())
	if !s.options.silent {
		if s.options.Debug {
			s.echo.Logger.SetLevel(log.DEBUG)
		} else {
			s.echo.Logger.SetLevel(log.INFO)
		}
		s.echo.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
			Format: loggerFormat,
		}))
	}

	s.echo.HideBanner = true
	s.echo.IPExtractor = echo.ExtractIPFromXFFHeader(
		echo.TrustLoopback(true),
		echo.TrustLinkLocal(false),
		echo.TrustPrivateNet(false)) // upstream proxy sets the X-Forwarded-For header

	s.echo.Pre(middleware.RemoveTrailingSlash())
	s.echo.Pre(middleware.Rewrite(map[string]string{
		"/favicon.ico": "/assets/images/favicon.ico",
		"/robots.txt":  "/assets/robots.txt",
		"/about":       "/help/about",
	}))

	anon := s.echo.Group("")
	anon.GET("/assets/*", s.assets.serve)
	anon.GET("/list/:token", s.renderList).Name = "render-filterlist"
	anon.POST("/filters/:name/render", s.viewFilterRender).Name = "view-filter-render"
	anon.GET("/should-reload", shouldReload)
	anon.GET("/news.atom", s.newsAtomHandler).Name = "news-atom"

	withAuth := s.echo.Group("")
	if s.options.KratosURL != "" {
		withAuth.Use(s.buildOryMiddleware())
	}
	withAuth.Use(middleware.CSRFWithConfig(middleware.CSRFConfig{
		TokenLookup:    "form:" + csrfLookup,
		ContextKey:     csrfLookup,
		CookieName:     csrfLookup,
		CookiePath:     "/",
		CookieSameSite: http.SameSiteStrictMode,
		CookieHTTPOnly: true,
	}))

	withAuth.GET("/", s.landingPageHandler).Name = "landing"
	withAuth.GET("/help", s.helpPages).Name = "help-main"
	withAuth.GET("/help/:page", s.helpPages).Name = "help"
	withAuth.GET("/news", s.newsHandler).Name = "news"

	withAuth.GET("/filters", s.listFilters).Name = "list-filters"
	withAuth.GET("/filters/tag/:tag", s.listFilters).Name = "filters-for-tag"

	withAuth.GET("/filters/:name", s.viewFilter).Name = "view-filter"
	withAuth.POST("/filters/:name", s.viewFilter)

	withAuth.GET("/list/:token/export", s.exportList).Name = "export-filterlist"
	withAuth.GET("/user/account", s.userAccount).Name = "user-account"
	withAuth.GET("/user/forms/:type", s.renderKratosForm)
	withAuth.POST("/user/start/:type", s.startKratosFlow).Name = "start-flow"
	withAuth.POST("/user/rotate-token", s.rotateListToken).Name = "rotate-list-token"
}

func shouldReload(c echo.Context) error {
	if !strings.HasPrefix(c.Request().Host, "localhost") {
		return echo.NewHTTPError(http.StatusNotFound)
	}
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
}

func (s *Server) absoluteReverse(c echo.Context, name string, params ...interface{}) string {
	u := &url.URL{
		Scheme: c.Scheme(),
		Host:   c.Request().Host,
		Path:   s.echo.Reverse(name, params...),
	}
	return u.String()
}

// redirectToPage the user to another page, either via htmx client-side redirect (form submissions)
// or http 302 redirect (direct access, js disabled)
func (s *Server) redirectToPage(c echo.Context, name string, params ...interface{}) error {
	return s.redirect(c, http.StatusFound, s.echo.Reverse(name, params...))
}

func (s *Server) redirect(c echo.Context, code int, target string) error {
	if c.Request().Header.Get("HX-Request") == "true" {
		c.Response().Header().Set("HX-Redirect", target)
		return c.NoContent(200)
	}
	return c.Redirect(code, target)
}

func (s *Server) buildPageContext(c echo.Context, title string) *pages.Context {
	var section string
	switch c.Request().URL.Path {
	case "/help/about":
		section = "help/about"
	case "/help/contributing":
		section = "help/contributing"
	default:
		for _, s := range strings.Split(c.Path(), "/") {
			if s != "" {
				section = s
				break
			}
		}
	}

	context := &pages.Context{
		CurrentSection:  section,
		NavigationLinks: navigationLinks,
		Title:           title,
		MainDomain:      c.Request().Host == mainDomain,
		HotReload:       s.options.Reload,
		RequestInfo:     c,
	}
	if _, err := c.Cookie(hasAccountCookieName); err == nil {
		context.UserHasAccount = true
	}
	if t, ok := c.Get(csrfLookup).(string); ok {
		context.CSRFToken = t
	}
	if u := getUser(c); u != nil {
		context.UserID = u.Id()
		context.UserLoggedIn = true
		context.Preferences, _ = s.preferences.Get(c, context.UserID)
		if context.Preferences != nil {
			latest, _ := s.releases.GetLatestAt()
			context.HasNews = latest.After(context.Preferences.NewsCursor)
		}
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
			loggedTag := fmt.Sprintf("logged:%t", c.Get(hasKratosContextKey) != nil)
			duration := time.Since(start)
			_ = dsd.Distribution("letsblockit.request_duration", float64(duration.Nanoseconds()), []string{loggedTag}, 1)
			_ = dsd.Incr("letsblockit.request_count", []string{loggedTag, fmt.Sprintf("status:%d", c.Response().Status)}, 1)
			return nil
		}
	}
}

func collectStats(log echo.Logger, store db.Store, dsd *statsd.Client) {
	collect := func() {
		stats, err := store.GetStats(context.Background())
		if err != nil {
			log.Error("cannot collect db stats: " + err.Error())
			return
		}
		_ = dsd.Gauge("letsblockit.total_list_count", float64(stats.ListsTotal), nil, 1)
		_ = dsd.Gauge("letsblockit.active_list_count", float64(stats.ListsActive), nil, 1)
		_ = dsd.Gauge("letsblockit.fresh_list_count", float64(stats.ListsFresh), nil, 1)

		instances, err := store.GetInstanceStats(context.Background())
		if err != nil {
			log.Error("cannot collect db stats: " + err.Error())
			return
		}
		for _, i := range instances {
			tags := []string{"filter_name:" + i.FilterName}
			_ = dsd.Gauge("letsblockit.instance_count", float64(i.Count), tags, 1)
		}
	}

	_ = dsd.Incr("letsblockit.startup", nil, 1)
	collect()
	for range time.Tick(5 * time.Minute) {
		collect()
	}
}
