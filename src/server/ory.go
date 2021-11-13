package server

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/labstack/echo/v4"
)

const (
	oryCookieNamePrefix = "ory_session_"
	oryWhoamiPath       = "/api/kratos/public/sessions/whoami"
	oryLogoutInfoPath   = "/api/kratos/public/self-service/logout/browser"
	oryLogoutDestPath   = "/.ory/api/kratos/public/self-service/logout?token="
	userContextKey      = "_user"
)

var oryClientRetries = 3

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
	client := buildRetryableClient(s.echo.Logger)
	endpoint := s.options.OryUrl + oryLogoutInfoPath

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
	if info.Token == "" {
		return "", fmt.Errorf("empty logout token")
	}
	return oryLogoutDestPath + info.Token, nil
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
