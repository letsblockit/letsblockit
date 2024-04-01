package server

import (
	"errors"
	"fmt"
	"io/fs"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/DataDog/datadog-go/v5/statsd"
	"github.com/coreos/go-systemd/activation"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/letsblockit/letsblockit/data"
	"github.com/letsblockit/letsblockit/src/db"
	"github.com/letsblockit/letsblockit/src/filters"
	"github.com/letsblockit/letsblockit/src/news"
	"github.com/letsblockit/letsblockit/src/pages"
	"github.com/letsblockit/letsblockit/src/users"
	"github.com/letsblockit/letsblockit/src/users/auth"
	"github.com/vearutop/statigz"
	"gopkg.in/natefinch/lumberjack.v2"
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
	healthPath = "/_health"
)

type Options struct {
	Address             string `group:"Networking" default:"127.0.0.1:8765" help:"address to listen to"`
	UseSystemdSocket    bool   `group:"Networking" help:"use a systemd socket instead of opening a port"`
	GzipResponses       bool   `group:"Networking" help:"compress most responses with gzip"`
	DatabaseUrl         string `group:"Database" default:"postgresql:///letsblockit" help:"psql database to connect to"`
	DatabasePoolOptions string `group:"Database" default:"" help:"pgxpool additional options"`
	AuthMethod          string `group:"Authentication" required:"" enum:"kratos,proxy" help:"authentication method to use"`
	AuthKratosUrl       string `group:"Authentication" default:"http://localhost:4000/.ory" help:"url of the kratos API, defaults to using local ory proxy"`
	AuthProxyHeaderName string `group:"Authentication" placeholder:"X-Auth-Request-User" help:"name for the cookie set by the reverse proxy"`
	LogLevel            string `group:"Development" default:"info" enum:"debug,info,warn,error,off" help:"http log level"`
	CacheDir            string `group:"Development" placeholder:"/tmp" help:"folder to cache external resources in during local development"`
	HotReload           bool   `group:"Development" help:"reload frontend when the backend restarts"`
	StatsdTarget        string `group:"Monitoring" placeholder:"localhost:8125" help:"address to send statsd metrics to, disabled by default"`
	LogsFolder          string `group:"Monitoring" help:"output access logs to files instead of stdout"`
	PlausibleScript     string `group:"Monitoring" help:"URL to the Plausible script for web analytics"`
	ListDownloadDomain  string `group:"Miscellaneous" help:"domain to use for list downloads, leave empty to use the main domain"`
	OfficialInstance    bool   `group:"Miscellaneous" help:"turn on behaviours specific to the official letsblock.it instances"`
	Sunset              bool   `group:"Miscellaneous" help:"Project sunset mode: require an existing account to use"`
	DryRun              bool   `hidden:""`
}

var navigationLinks = []struct {
	Name   string
	Target string
}{{
	Name:   "Template list",
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
	assets      http.Handler
	auth        auth.Backend
	bans        *users.BanManager
	echo        *echo.Echo
	filters     *filters.Repository
	now         func() time.Time
	options     *Options
	pages       PageRenderer
	preferences *users.PreferenceManager
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
	// Open the socket as soon as possible to avoid dropping requests on restarts
	if s.options.UseSystemdSocket {
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
	} else {
		listener, err := net.Listen("tcp", s.options.Address)
		if err != nil {
			return err
		}
		s.echo.Listener = listener
	}

	if s.options.StatsdTarget != "" {
		dsd, err := statsd.New(s.options.StatsdTarget, statsd.WithoutTelemetry())
		if err != nil {
			return err
		}
		s.statsd = dsd
		s.echo.Use(buildDogstatsMiddleware(dsd))
	} else {
		s.statsd = &statsd.NoOpClient{}
	}

	var helpers map[string]interface{}
	concurrentRunOrPanic([]func([]error){
		func(errs []error) { s.assets = statigz.FileServer(data.Assets) },
		func(errs []error) { s.pages, errs[0] = pages.LoadPages() },
		func(errs []error) { s.filters, errs[0] = filters.Load(data.Templates, data.Presets) },
		func(errs []error) {
			s.store, errs[0] = db.Connect(s.options.DatabaseUrl, s.options.DatabasePoolOptions, s.statsd)
			if errs[0] == nil {
				errs[0] = db.Migrate(s.options.DatabaseUrl)
			}
			if errs[0] == nil {
				s.bans, errs[0] = users.LoadUserBans(s.store)
			}
			if errs[0] == nil {
				s.preferences, errs[0] = users.NewPreferenceManager(s.store)
			}
		},
		func(errs []error) { helpers, errs[0] = buildHelpers(s.echo) },
		func(errs []error) {
			tpl := data.Templates.(fs.ReadDirFS)
			s.releases, errs[0] = news.DownloadReleases(news.GithubReleasesEndpoint, s.options.CacheDir, s.options.OfficialInstance, tpl)
		},
	})

	if s.options.LogsFolder != "" {
		if err := os.MkdirAll(s.options.LogsFolder, 0750); err != nil {
			return err
		}
		target := fmt.Sprintf("%s/lbi-%d.log", s.options.LogsFolder, os.Getpid())
		fmt.Println("Writing access logs to", target)
		s.echo.Logger.SetOutput(&lumberjack.Logger{
			Filename:   target,
			MaxSize:    100, // megabytes
			MaxBackups: 3,
		})
	}

	switch s.options.AuthMethod {
	case "kratos":
		if s.options.AuthKratosUrl == "" {
			return fmt.Errorf("missing required parameter auth-kratos-url")
		}
		s.auth = auth.NewOryBackend(s.options.AuthKratosUrl, s.pages, s.statsd, s.options.Sunset)
	case "proxy":
		if s.options.AuthProxyHeaderName == "" {
			return fmt.Errorf("missing required parameter auth-proxy-header-name")
		}
		s.auth = auth.NewProxy(s.options.AuthProxyHeaderName)
	default:
		return fmt.Errorf("unsupported auth method %s", s.options.AuthMethod)
	}

	s.pages.RegisterHelpers(helpers)
	s.pages.RegisterContextBuilder(s.buildPageContext)
	s.setupRouter()
	if s.options.DryRun {
		return ErrDryRunFinished
	}

	if s.options.StatsdTarget != "" {
		go collectBusinessStats(s.echo.Logger, s.store, s.statsd)
		go collectMemStats(s.statsd)
	}

	return s.echo.Start(s.options.Address)
}

func (s *Server) setupRouter() {
	switch s.options.LogLevel {
	case "debug":
		s.echo.Logger.SetLevel(log.DEBUG)
	case "info":
		s.echo.Logger.SetLevel(log.INFO)
	case "warn":
		s.echo.Logger.SetLevel(log.WARN)
	case "error":
		s.echo.Logger.SetLevel(log.ERROR)
	case "off":
		s.echo.Logger.SetLevel(log.OFF)
	}
	s.echo.Use(
		middleware.Recover(),
		middleware.LoggerWithConfig(middleware.LoggerConfig{
			Format: loggerFormat,
			Skipper: func(c echo.Context) bool {
				return c.Request().URL.Path == healthPath
			},
		}),
	)

	s.echo.HideBanner = true
	s.echo.IPExtractor = echo.ExtractIPFromXFFHeader()

	s.echo.Pre(middleware.RemoveTrailingSlash())
	s.echo.Pre(middleware.Rewrite(map[string]string{
		"/favicon.ico": "/assets/images/favicon.ico",
		"/robots.txt":  "/assets/robots.txt",
		"/about":       "/help/about",
	}))

	// Raw routes
	s.echo.GET(healthPath, func(c echo.Context) error { return c.String(200, "OK") })
	s.echo.GET("/assets/*", echo.WrapHandler(s.assets), assetCacheControlMiddleware)
	s.echo.HEAD("/assets/*", echo.WrapHandler(s.assets))
	s.echo.GET("/filters/youtube-streams-chat", func(c echo.Context) error {
		return s.pages.RedirectToPage(c, "view-filter", "youtube-cleanup")
	})
	if s.options.HotReload {
		s.echo.GET("/should-reload", shouldReload)
	}

	var middlewares []echo.MiddlewareFunc
	if s.options.GzipResponses {
		middlewares = append(middlewares, middleware.GzipWithConfig(middleware.GzipConfig{Level: 6}))
	}
	zippedRoutes := s.echo.Group("", middlewares...)
	zippedRoutes.POST("/filters/:name/render", s.viewFilterRender).Name = "view-filter-render"
	zippedRoutes.GET("/list/:token", s.renderList).Name = "render-filterlist"
	zippedRoutes.GET("/news.atom", s.newsAtomHandler).Name = "news-atom"

	authedRoutes := zippedRoutes.Group("",
		s.auth.BuildMiddleware(),
		func(next echo.HandlerFunc) echo.HandlerFunc {
			return func(c echo.Context) error {
				if s.bans.IsBanned(auth.GetUserId(c)) {
					return echo.ErrForbidden
				}
				return next(c)
			}
		},
		middleware.CSRFWithConfig(middleware.CSRFConfig{
			TokenLookup:    "form:" + csrfLookup,
			ContextKey:     csrfLookup,
			CookieName:     csrfLookup,
			CookiePath:     "/",
			CookieSameSite: http.SameSiteStrictMode,
			CookieHTTPOnly: true,
			ErrorHandler: func(err error, c echo.Context) error {
				if auth.GetUserId(c) == "" {
					return nil // Ignore missing CSRF if the session is anonymous
				}
				return err
			},
		}),
	)
	s.auth.RegisterRoutes(authedRoutes)

	authedRoutes.GET("/", s.landingPageHandler).Name = "landing"
	authedRoutes.GET("/help", s.helpPages).Name = "help-main"
	authedRoutes.GET("/help/:page", s.helpPages).Name = "help"
	authedRoutes.GET("/news", s.newsHandler).Name = "news"
	authedRoutes.GET("/sunset", s.sunsetHandler).Name = "sunset"

	authedRoutes.GET("/filters", s.listFilters).Name = "list-filters"
	authedRoutes.GET("/filters/tag/:tag", s.listFilters).Name = "filters-for-tag"

	authedRoutes.GET("/filters/:name", s.viewFilter).Name = "view-filter"
	authedRoutes.POST("/filters/:name", s.viewFilter)

	authedRoutes.GET("/export/:token", s.exportList).Name = "export-filterlist"
	authedRoutes.GET("/user/account", s.userAccount).Name = "user-account"
	authedRoutes.POST("/user/rotate-token", s.rotateListToken).Name = "rotate-list-token"
	authedRoutes.POST("/user/preferences", s.updatePreferences).Name = "update-preferences"
}

func assetCacheControlMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if c.Request().URL.Query().Has("h") {
			c.Response().Header().Set("Cache-Control", "max-age=604800") // 7 days
		}
		return next(c)
	}
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
	c.Response().Header().Set("Access-Control-Allow-Origin", "*")
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
		CurrentSection:   section,
		NavigationLinks:  navigationLinks,
		Title:            title,
		OfficialInstance: s.options.OfficialInstance,
		PlausibleScript:  s.options.PlausibleScript,
		GreyLogo:         s.options.OfficialInstance && c.Request().Host != mainDomain,
		HotReload:        s.options.HotReload,
		RequestInfo:      c,
		UserHasAccount:   auth.HasAccount(c),
	}
	if t, ok := c.Get(csrfLookup).(string); ok {
		context.CSRFToken = t
	}
	if u := auth.GetUserId(c); u != "" {
		context.UserID = u
		context.UserLoggedIn = true
		context.Preferences, _ = s.preferences.Get(c, context.UserID)
		if context.Preferences != nil {
			latest := s.releases.GetLatestAt()
			context.HasNews = latest.After(context.Preferences.NewsCursor.Time)
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

func getEtag(c echo.Context) string {
	return c.Request().Header.Get("If-None-Match")
}
