package pages

import (
	"net/http"

	"github.com/letsblockit/letsblockit/src/db"
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
	NakedContent     bool
	Page             *page
	Sidebar          *page
	NoBoost          bool
	HotReload        bool
	OfficialInstance bool
	GreyLogo         bool
	RequestInfo      RequestInfo

	CurrentSection  string
	NavigationLinks interface{}
	Title           string

	UserID         string
	UserLoggedIn   bool
	UserHasAccount bool
	HasNews        bool
	ColorMode      db.ColorMode
	Preferences    *db.UserPreference
	CSRFToken      string

	Data ContextData
}

func (c *Context) Add(key string, value interface{}) {
	if c.Data == nil {
		c.Data = make(ContextData)
	}
	c.Data[key] = value
}
