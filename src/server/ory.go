package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/xvello/letsblockit/src/pages"
)

const (
	userContextKey      = "_user"
	oryCookieNamePrefix = "ory_session_"
	oryGetFlowPattern   = "/api/kratos/public/self-service/%s/flows?id=%s"
	oryStartFlowPattern = "/api/kratos/public/self-service/%s/browser"
	oryLogoutInfoPath   = "/api/kratos/public/self-service/logout/browser"
	oryWhoamiPath       = "/api/kratos/public/sessions/whoami"
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
			Intro: `There was an error. If it persists, please <a href="https://github.com/xvello/letsblockit/issues">open an issue</a>.`,
		},
		"login": {
			Title: "Log into your account",
			Tabs:  loginTabs,
		},
		"recovery": {
			Title: "Recover your account",
			Tabs:  loginTabs,
			Intro: `<div class="alert alert-warning">
<h5>We are currently investigating a bug leading to the recovery links pointing to the wrong domain</h5>
In the meantime, please replace
<code>https://goofy-tereshkova-bmg1nc7mqv.projects.oryapis.com/api/</code> with
<code>https://letsblock.it/.ory/api/</code> in the link you will receive.
</div>
Enter your e-mail below to receive a recovery link by e-mail.`,
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

func (u *oryUser) Id() uuid.UUID {
	if u == nil {
		return uuid.Nil
	}
	return u.Identity.Id
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
	endpoint := s.options.KratosURL + oryWhoamiPath

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cookies := c.Request().Header.Get(echo.HeaderCookie)
			if !strings.Contains(cookies, oryCookieNamePrefix) {
				return next(c)
			}
			var user oryUser
			if err := s.queryKratos(c, "whoami", endpoint, &user); err != nil {
				s.echo.Logger.Error("auth error: %w", err)
			} else if user.IsActive() {
				c.Set(userContextKey, &user)
			}
			return next(c)
		}
	}
}

// getLogoutUrl retrieves the logout url for the current session by calling the proxy
func (s *Server) getLogoutUrl(c echo.Context) (string, error) {
	var info oryLogoutInfo
	if err := s.queryKratos(c, "logout", s.options.KratosURL+oryLogoutInfoPath, &info); err != nil {
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
		endpoint := s.options.KratosURL + fmt.Sprintf(oryGetFlowPattern, formType, flowID)
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

	switch target {
	case "logout":
		target, err := s.getLogoutUrl(c)
		if err != nil {
			return nil
		}
		return s.redirect(c, http.StatusSeeOther, target)
	case "loginOrRegistration":
		if _, err := c.Cookie(hasAccountCookieName); err == nil {
			target = "login"
		} else {
			target = "registration"
		}
	}
	return s.redirect(c, http.StatusSeeOther, s.options.KratosURL+fmt.Sprintf(oryStartFlowPattern, target))
}

func (s *Server) queryKratos(c echo.Context, typeTag, endpoint string, body interface{}) error {
	start := time.Now()

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
