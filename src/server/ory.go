package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/labstack/echo/v4"
	"github.com/xvello/letsblockit/src/pages"
)

const (
	oryCookieNamePrefix = "ory_session_"
	oryWhoamiPath       = "/api/kratos/public/sessions/whoami"
	oryLogoutInfoPath   = "/api/kratos/public/self-service/logout/browser"
	oryLogoutDestPath   = "/.ory/api/kratos/public/self-service/logout?token="
	userContextKey      = "_user"
	oryGetFlowPattern   = "http://localhost:4000/.ory/api/kratos/public/self-service/%s/flows?id=%s"
)

var oryClientRetries = 3

type formSettings struct {
	Title string
	Intro string
}

var supportedForms = map[string]formSettings{
	"error": {
		Title: "Account management error",
		Intro: `There was an error. If it persists, please <a href="https://github.com/xvello/letsblockit/issues">open an issue</a>.`,
	},
	"login": {
		Title: "Log into your account",
		Intro: `Enter your e-mail and password to login, or <a href="/.ory/ui/registration">click here to create a new account</a>.`,
	},
	"registration": {
		Title: "Create a new account",
		Intro: `Already have an account? <a href="/.ory/ui/login">Sign in instead</a>.`,
	},
	"settings": {
		Title: "Account settings",
		Intro: `You can change your e-mail or password here. If you change your e-mail, you will receive a new validation email.`,
	},
	"verification": {
		Title: "Verify your account",
	},
}

// oryUser holds the parts of the kratos user we care about.
type oryUser struct {
	Active   bool
	Identity struct {
		Id        uuid.UUID
		Addresses []struct {
			Verified bool
		} `json:"verifiable_addresses"`
	}
}

type oryLogoutInfo struct {
	Token string `json:"logout_token"`
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

func (u *oryUser) IsVerified() bool {
	if u == nil {
		return false
	}
	for _, addr := range u.Identity.Addresses {
		if addr.Verified {
			return true
		}
	}
	return false
}

// leveledLogger implements retryablehttp.LeveledLogger around an echo.Logger
type leveledLogger struct {
	log echo.Logger
}

func (l *leveledLogger) Error(msg string, keysAndValues ...interface{}) {
	l.log.Errorf(msg, keysAndValues)
}
func (l *leveledLogger) Info(msg string, keysAndValues ...interface{}) {
	l.log.Infof(msg, keysAndValues)
}
func (l *leveledLogger) Debug(msg string, keysAndValues ...interface{}) {
	l.log.Debugf(msg, keysAndValues)
}
func (l *leveledLogger) Warn(msg string, keysAndValues ...interface{}) {
	l.log.Warnf(msg, keysAndValues)
}

// buildOryMiddleware tries to resolve an Ory Cloud session from the cookies.
// If it succeeds, a "user" value is added to the context for use by handlers.
func (s *Server) buildOryMiddleware() echo.MiddlewareFunc {
	client := buildRetryableClient(s.echo.Logger)
	endpoint := s.options.OryUrl + oryWhoamiPath

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cookies := c.Request().Header.Get(echo.HeaderCookie)
			if !strings.Contains(cookies, oryCookieNamePrefix) {
				return next(c)
			}

			req, err := retryablehttp.NewRequest("GET", endpoint, nil)
			if err != nil {
				return fmt.Errorf("failed to instantiate request: %w", err)
			}
			req.Header.Set(echo.HeaderCookie, cookies)
			req.Header.Set("Accept", "application/json")
			res, err := client.Do(req)
			if err != nil {
				s.echo.Logger.Error("failed to query kratos: %w", err)
				return next(c)
			}
			defer res.Body.Close()

			var user oryUser
			if err = json.NewDecoder(res.Body).Decode(&user); err != nil {
				s.echo.Logger.Error("failed to decode session object: %w", err)
				return next(c)
			}
			if user.IsActive() {
				c.Set(userContextKey, &user)
			}

			return next(c)
		}
	}
}

// getLogoutUrl retrieves the logout token for the current session
// and and builds the redirect URL.
func (s *Server) getLogoutUrl(c echo.Context) (string, error) {
	var info oryLogoutInfo
	if err := s.queryKratos(c, s.options.OryUrl+oryLogoutInfoPath, &info); err != nil {
		return "", err
	}
	if info.Token == "" {
		return "", fmt.Errorf("empty logout token")
	}
	return oryLogoutDestPath + info.Token, nil
}

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
		endpoint := fmt.Sprintf(oryGetFlowPattern, formType, flowID)
		if err := s.queryKratos(c, endpoint, &body); err != nil {
			return nil, err
		}
		ui, ok := body["ui"]
		if !ok {
			return nil, errors.New("no UI field in payload")
		}

		hc := s.buildPageContext(c, formSettings.Title)
		hc.Add("type", formType)
		hc.Add("ui", ui)
		hc.Add("settings", formSettings)
		c.Logger().Warn(ui)
		return hc, nil
	}()
	if err != nil {
		c.Logger().Warnf("Kratos form failure, falling-back to managed UI: %s", err.Error())
		return c.Redirect(http.StatusFound, fmt.Sprintf("/.ory/ui/%s?flow=%s", formType, flowID))
	}
	return s.pages.Render(c, "kratos-form", hc)
}

func (s *Server) queryKratos(c echo.Context, endpoint string, body interface{}) error {
	client := buildRetryableClient(s.echo.Logger)

	req, err := retryablehttp.NewRequest("GET", endpoint, nil)
	if err != nil {
		return fmt.Errorf("failed to instantiate request: %w", err)
	}
	req.Header.Set(echo.HeaderCookie, c.Request().Header.Get(echo.HeaderCookie))
	req.Header.Set("Accept", "application/json")
	res, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to query kratos: %w", err)
	}
	defer res.Body.Close()

	if err := json.NewDecoder(res.Body).Decode(body); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}
	return nil
}

func buildRetryableClient(logger echo.Logger) *retryablehttp.Client {
	client := retryablehttp.NewClient()
	client.RetryMax = oryClientRetries
	client.HTTPClient.Timeout = 2 * time.Second
	client.Logger = &leveledLogger{logger}
	return client
}

func getUser(c echo.Context) *oryUser {
	if u, ok := c.Get(userContextKey).(*oryUser); ok {
		return u
	}
	return nil
}
