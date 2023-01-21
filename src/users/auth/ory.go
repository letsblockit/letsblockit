package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/DataDog/datadog-go/v5/statsd"
	"github.com/google/uuid"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/labstack/echo/v4"
	"github.com/letsblockit/letsblockit/src/pages"
	"zgo.at/zcache/v2"
)

const (
	oryCookieNamePrefix = "ory_session_"
	oryGetFlowPattern   = "/self-service/%s/flows?id=%s"
	oryStartFlowPattern = "/self-service/%s/browser"
	oryReturnToPattern  = "?return_to=%s"
	oryLogoutInfoPath   = "/self-service/logout/browser"
	oryWhoamiPath       = "/sessions/whoami"
	returnToKey         = "return_to"
)

// renderer is fulfilled by pages.Pages
type renderer interface {
	BuildPageContext(c echo.Context, title string) *pages.Context
	Redirect(c echo.Context, code int, target string) error
	Render(c echo.Context, name string, data *pages.Context) error
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
			Intro: `Enter your e-mail below to receive a recovery code from <code>no-reply@oryapis.com</code>. Input that code to login and set a new password.`,
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

// OryBackend is used for the Kratos auth method. It is used by the official instance with Ory Cloud.
// It should work with a self-hosted Kratos, but this has not been tested.
type OryBackend struct {
	client   *retryablehttp.Client
	rootUrl  string
	renderer renderer
	statsd   statsd.ClientInterface
}

func NewOryBackend(rootUrl string, renderer renderer, statsd statsd.ClientInterface) *OryBackend {
	client := retryablehttp.NewClient()
	client.RetryMax = 2
	client.RetryWaitMin = 100 * time.Millisecond
	client.RetryWaitMax = time.Second
	client.HTTPClient.Timeout = 5 * time.Second
	return &OryBackend{
		client:   client,
		rootUrl:  rootUrl,
		renderer: renderer,
		statsd:   statsd,
	}
}

func (o *OryBackend) RegisterRoutes(group EchoRouter) {
	group.GET("/user/forms/:type", o.renderKratosForm)
	group.POST("/user/action/:type", o.startKratosFlow).Name = userActionRouteName
	group.GET("/user/action/loginOrRegistration", o.startLoginFlow).Name = "simple-login-link"
}

// BuildMiddleware tries to resolve an Ory Cloud session from the cookies.
// If it succeeds, a "user" value is added to the context for use by handlers.
func (o *OryBackend) BuildMiddleware() echo.MiddlewareFunc {
	authCache := zcache.New[string, string](15*time.Minute, 10*time.Minute)
	endpoint := o.rootUrl + oryWhoamiPath

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cookies := c.Request().Header.Get(echo.HeaderCookie)
			if !strings.Contains(cookies, oryCookieNamePrefix) {
				return next(c)
			}

			if u, ok := authCache.Get(cookies); ok {
				setUserId(c, u)
				return next(c)
			}

			var user oryUser
			if err := o.queryKratos(c, "whoami", endpoint, &user); err != nil {
				c.Logger().Error("auth error: %w", err)
			} else if user.IsActive() {
				id := user.Id()
				authCache.Set(cookies, id)
				setUserId(c, id)
			}

			if _, err := c.Cookie(hasAccountCookieName); err == http.ErrNoCookie {
				c.SetCookie(&http.Cookie{
					Name:     hasAccountCookieName,
					Value:    hasAccountCookieValue,
					Path:     "/",
					Expires:  time.Now().AddDate(10, 0, 0),
					HttpOnly: true,
					SameSite: http.SameSiteStrictMode,
				})
			}

			return next(c)
		}
	}
}

// getLogoutUrl retrieves the logout url for the current session by calling the proxy
func (o *OryBackend) getLogoutUrl(c echo.Context) (string, error) {
	var info oryLogoutInfo
	if err := o.queryKratos(c, "logout", o.rootUrl+oryLogoutInfoPath, &info); err != nil {
		return "", err
	}
	if info.URL == "" {
		return "", fmt.Errorf("empty logout url")
	}
	return info.URL, nil
}

// renderKratosForm retrieves the flow definition from the proxy and renders the form.
// If any error occurs, the user is redirected to the Cloud Managed UI as a fallback.
func (o *OryBackend) renderKratosForm(c echo.Context) error {
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
		endpoint := o.rootUrl + fmt.Sprintf(oryGetFlowPattern, formType, flowID)
		if err := o.queryKratos(c, "flow", endpoint, &body); err != nil {
			return nil, err
		}
		ui, ok := body["ui"]
		if !ok {
			return nil, errors.New("no UI field in payload")
		}

		hc := o.renderer.BuildPageContext(c, formSettings.Title)
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
		return o.renderer.Redirect(c, http.StatusFound, fmt.Sprintf("%s/ui/%s?flow=%s", o.rootUrl, formType, flowID))
	}
	return o.renderer.Render(c, "kratos-form", hc)
}

// startKratosFlow redirects the requested user to the Kratos flow. It is used via POST instead
// of direct GET links to avoid search engines and preloading logics starting Kratos flows.
func (o *OryBackend) startKratosFlow(c echo.Context) error {
	target := c.Param("type")
	var allowReturnTo bool

	switch target {
	case "logout":
		target, err := o.getLogoutUrl(c)
		if err != nil {
			return nil
		}
		return o.renderer.Redirect(c, http.StatusSeeOther, target)
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

	redirect := o.rootUrl + fmt.Sprintf(oryStartFlowPattern, target)
	if allowReturnTo {
		if returnTo, inDomain := computeReturnTo(c); inDomain {
			redirect += fmt.Sprintf(oryReturnToPattern, returnTo)
		}
	}
	return o.renderer.Redirect(c, http.StatusSeeOther, redirect)
}

func (o *OryBackend) startLoginFlow(c echo.Context) error {
	c.SetParamNames("type")
	c.SetParamValues("loginOrRegistration")
	return o.startKratosFlow(c)
}

func (o *OryBackend) queryKratos(c echo.Context, typeTag, endpoint string, body interface{}) error {
	start := time.Now()
	c.Set(hasAuthContextKey, true)

	req, err := retryablehttp.NewRequest("GET", endpoint, nil)
	if err != nil {
		return fmt.Errorf("failed to instantiate request: %w", err)
	}
	req.Header.Set(echo.HeaderCookie, c.Request().Header.Get(echo.HeaderCookie))
	req.Header.Set("Accept", "application/json")
	res, err := o.client.Do(req)
	_ = o.statsd.Distribution("letsblockit.ory_request_duration", float64(time.Since(start).Nanoseconds()),
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
