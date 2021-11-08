package server

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/labstack/echo/v4"
)

const (
	oryCookieNamePrefix  = "ory_session_"
	oryWhoamiPattern     = "https://%s.projects.oryapis.com/api/kratos/public/sessions/whoami"
	oryLogoutInfoPattern = "https://%s.projects.oryapis.com/api/kratos/public/self-service/logout/browser"
	userContextKey       = "_user"
)

// oryUser holds the parts of the kratos user we care about.
type oryUser struct {
	Active   bool
	Identity struct {
		Id        string
		Addresses []struct {
			Verified bool
		} `json:"verifiable_addresses"`
	}
}

type oryLogoutInfo struct {
	LogoutUrl string `json:"logout_url"`
}

func (u *oryUser) Id() string {
	if u == nil {
		return ""
	}
	return u.Identity.Id
}

func (u *oryUser) IsActive() bool {
	if u == nil {
		return false
	}
	return u.Active && u.Identity.Id != ""
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
func buildOryMiddleware(project string, logger echo.Logger) echo.MiddlewareFunc {
	client := buildRetryableClient(logger)
	endpoint := fmt.Sprintf(oryWhoamiPattern, project)

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
				return fmt.Errorf("failed to query kratos: %w", err)
			}
			defer res.Body.Close()

			var user oryUser
			if err := json.NewDecoder(res.Body).Decode(&user); err != nil {
				return fmt.Errorf("failed to decode session object: %w", err)
			}
			if user.IsActive() {
				c.Set(userContextKey, &user)
			}

			return next(c)
		}
	}
}

// getLogoutUrl retrieves the logout URL for the current session
// and patches it to go through the proxy.
func getLogoutUrl(project string, c echo.Context) (string, error) {
	client := buildRetryableClient(c.Logger())
	endpoint := fmt.Sprintf(oryLogoutInfoPattern, project)
	req, err := retryablehttp.NewRequest("GET", endpoint, nil)
	if err != nil {
		return "", fmt.Errorf("failed to instantiate request: %w", err)
	}
	req.Header.Set(echo.HeaderCookie, c.Request().Header.Get(echo.HeaderCookie))
	req.Header.Set("Accept", "application/json")
	res, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to query kratos: %w", err)
	}
	defer res.Body.Close()

	var info oryLogoutInfo
	if err := json.NewDecoder(res.Body).Decode(&info); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}
	logout, err := url.Parse(info.LogoutUrl)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("/.ory%s?%s", logout.Path, logout.RawQuery), nil
}

func buildRetryableClient(logger echo.Logger) *retryablehttp.Client {
	client := retryablehttp.NewClient()
	client.RetryMax = 3
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
