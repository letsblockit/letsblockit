package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/letsblockit/letsblockit/src/pages"
	"zgo.at/zcache"
)

const (
	userContextKey      = "_user"
	hasKratosContextKey = "_has_kratos"
	oryCookieNamePrefix = "ory_session_"
	oryGetFlowPattern   = "/self-service/%s/flows?id=%s"
	oryStartFlowPattern = "/self-service/%s/browser"
	oryReturnToPattern  = "?return_to=%s"
	oryLogoutInfoPath   = "/self-service/logout/browser"
	oryWhoamiPath       = "/sessions/whoami"
	returnToKey         = "return_to"
)

var proxyClient = http.Client{
	Timeout: 30 * time.Second,
}

type formTab struct {
	Title string
	Type  string
}

var (
	loginTabs = []formTab{{
		Title: "Create account",
		Type:  "registration",
	}, {
		Title: "Login",
		Type:  "login",
	}, {
		Title: "Recover password",
		Type:  "recovery",
	}}

	supportedForms = map[string]struct {
		Title string
		Tabs  []formTab
		Intro string
	}{
		"error": {
			Title: "Account management error",
			Intro: `There was an error. If it persists, please <a href="https://github.com/letsblockit/letsblockit/issues">open an issue</a>.`,
		},
		"login": {
			Title: "Log into your account",
			Tabs:  loginTabs,
		},
		"recovery": {
			Title: "Recover your account",
			Tabs:  loginTabs,
			Intro: `Enter your e-mail below, we will send you a recovery link from <code>no-reply@ory.sh</code> to login and set a new password.`,
		},
		"registration": {
			Title: "Create a new account",
			Tabs:  loginTabs,
		},
		"settings": {
			Title: "Account settings",
			Intro: `You can change your e-mail or password here. If you change your e-mail, you will receive a new validation e-mail.`,
		},
		"verification": {
			Title: "Verify your account",
		},
	}
)

// oryUser holds the parts of the kratos user we care about.
type oryUser struct {
	Active   bool
	Identity struct {
		Id uuid.UUID
	}
}

type oryLogoutInfo struct {
	URL string `json:"logout_url"`
}

func (u *oryUser) Id() string {
	if u == nil {
		return ""
	}
	return u.Identity.Id.String()
}

func (u *oryUser) IsActive() bool {
	if u == nil {
		return false
	}
	return u.Active && u.Identity.Id != uuid.Nil
}

// buildOryMiddleware tries to resolve an Ory Cloud session from the cookies.
// If it succeeds, a "user" value is added to the context for use by handlers.
func (s *Server) buildOryMiddleware() echo.MiddlewareFunc {
	authCache := zcache.New(15*time.Minute, 10*time.Minute)
	endpoint := s.options.AuthKratosUrl + oryWhoamiPath

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cookies := c.Request().Header.Get(echo.HeaderCookie)
			if !strings.Contains(cookies, oryCookieNamePrefix) {
				return next(c)
			}

			if u, ok := authCache.Get(cookies); ok {
				c.Set(userContextKey, u)
				return next(c)
			}

			var user oryUser
			if err := s.queryKratos(c, "whoami", endpoint, &user); err != nil {
				s.echo.Logger.Error("auth error: %w", err)
			} else if s.bans.IsBanned(user.Id()) {
				return echo.ErrForbidden
			} else if user.IsActive() {
				authCache.SetDefault(cookies, &user)
				c.Set(userContextKey, &user)
			}

			return next(c)
		}
	}
}

// getLogoutUrl retrieves the logout url for the current session by calling the proxy
func (s *Server) getLogoutUrl(c echo.Context) (string, error) {
	var info oryLogoutInfo
	if err := s.queryKratos(c, "logout", s.options.AuthKratosUrl+oryLogoutInfoPath, &info); err != nil {
		return "", err
	}
	if info.URL == "" {
		return "", fmt.Errorf("empty logout url")
	}
	return info.URL, nil
}

// renderKratosForm retrieves the flow definition from the proxy and renders the form.
// If any error occurs, the user is redirected to the Cloud Managed UI as a fallback.
func (s *Server) renderKratosForm(c echo.Context) error {
	formType := c.Param("type")
	flowID := c.QueryParams().Get("flow")
	if formType == "" || flowID == "" {
		return fmt.Errorf("missing args, got type '%s', flow '%s'", formType, flowID)
	}
	hc, err := func() (*pages.Context, error) {
		formSettings, ok := supportedForms[formType]
		if !ok {
			return nil, fmt.Errorf("unsupported form type %s", formType)
		}

		body := make(map[string]interface{})
		endpoint := s.options.AuthKratosUrl + fmt.Sprintf(oryGetFlowPattern, formType, flowID)
		if err := s.queryKratos(c, "flow", endpoint, &body); err != nil {
			return nil, err
		}
		ui, ok := body["ui"]
		if !ok {
			return nil, errors.New("no UI field in payload")
		}

		hc := s.buildPageContext(c, formSettings.Title)
		hc.NoBoost = true
		hc.Add("type", formType)
		hc.Add("ui", ui)
		hc.Add("settings", formSettings)
		if rto, ok := body[returnToKey]; ok {
			hc.Add(returnToKey, rto)
		}
		return hc, nil
	}()
	if err != nil {
		c.Logger().Warnf("falling-back to managed UI: %s", err.Error())
		return c.Redirect(http.StatusFound, fmt.Sprintf("/.ory/ui/%s?flow=%s", formType, flowID))
	}
	return s.pages.Render(c, "kratos-form", hc)
}

// startKratosFlow redirects the requested user to the Kratos flow. It is used via POST instead
// of direct GET links to avoid search engines and preloading logics starting Kratos flows.
func (s *Server) startKratosFlow(c echo.Context) error {
	target := c.Param("type")
	var allowReturnTo bool

	switch target {
	case "logout":
		target, err := s.getLogoutUrl(c)
		if err != nil {
			return nil
		}
		return s.redirect(c, http.StatusSeeOther, target)
	case "loginOrRegistration":
		allowReturnTo = true
		if _, err := c.Cookie(hasAccountCookieName); err == nil {
			target = "login"
		} else {
			target = "registration"
		}
	case "login", "registration":
		allowReturnTo = true
	}

	redirect := s.options.AuthKratosUrl + fmt.Sprintf(oryStartFlowPattern, target)
	if allowReturnTo {
		if returnTo, inDomain := computeReturnTo(c); inDomain {
			redirect += fmt.Sprintf(oryReturnToPattern, returnTo)
		}
	}
	return s.redirect(c, http.StatusSeeOther, redirect)
}

func (s *Server) queryKratos(c echo.Context, typeTag, endpoint string, body interface{}) error {
	start := time.Now()
	c.Set(hasKratosContextKey, true)

	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return fmt.Errorf("failed to instantiate request: %w", err)
	}
	req.Header.Set(echo.HeaderCookie, c.Request().Header.Get(echo.HeaderCookie))
	req.Header.Set("Accept", "application/json")
	res, err := proxyClient.Do(req)
	_ = s.statsd.Distribution("letsblockit.ory_request_duration", float64(time.Since(start).Nanoseconds()),
		[]string{"type:" + typeTag, fmt.Sprintf("ok:%t", err == nil && res.StatusCode == http.StatusOK)}, 1)

	if err != nil {
		return fmt.Errorf("failed to query kratos: %w", err)
	}
	defer res.Body.Close()

	if err := json.NewDecoder(res.Body).Decode(body); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}
	return nil
}

func getUser(c echo.Context) *oryUser {
	if u, ok := c.Get(userContextKey).(*oryUser); ok {
		return u
	}
	return nil
}

// computeReturnTo looks into the form input for a return_to value and falls back to the referer.
// It also checks that the return to is in the current domain to avoid phishing.
func computeReturnTo(c echo.Context) (returnTo string, inDomain bool) {
	if values, err := c.FormParams(); err == nil {
		returnTo = values.Get(returnToKey)
	}
	if returnTo == "" {
		returnTo = c.Request().Referer()
	}

	parsed, err := url.Parse(returnTo)
	inDomain = err == nil && parsed.Host == c.Request().Host
	return
}
