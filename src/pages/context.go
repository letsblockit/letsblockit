package pages

import (
	"github.com/google/uuid"
)

type ContextData map[string]interface{}

type Context struct {
	NakedContent bool
	Scripts      []string
	Page         *page

	CurrentSection  string
	NavigationLinks interface{}
	Title           string

	UserID       uuid.UUID
	UserLoggedIn bool
	UserVerified bool

	Data ContextData
}

func (c *Context) Add(key string, value interface{}) {
	if c.Data == nil {
		c.Data = make(ContextData)
	}
	c.Data[key] = value
}
