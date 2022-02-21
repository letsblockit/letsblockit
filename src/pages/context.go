package pages

import (
	"net/http"

	"github.com/google/uuid"
)

// RequestInfo is a subset of echo.Context
type RequestInfo interface {
	// Request returns `*http.Request`.
	Request() *http.Request
	// Scheme returns the HTTP protocol scheme, `http` or `https`.
	Scheme() string
}

type ContextData map[string]interface{}

type Context struct {
	NakedContent bool
	Page         *page
	Sidebar      *page
	NoBoost      bool
	HotReload    bool
	MainDomain   bool
	RequestInfo  RequestInfo

	CurrentSection  string
	NavigationLinks interface{}
	Title           string

	UserID         uuid.UUID
	UserLoggedIn   bool
	UserHasAccount bool
	CSRFToken      string

	Data ContextData
}

func (c *Context) Add(key string, value interface{}) {
	if c.Data == nil {
		c.Data = make(ContextData)
	}
	c.Data[key] = value
}
