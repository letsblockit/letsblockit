package pages

type Context struct {
	NakedContent bool
	Scripts      []string
	Page         *page

	CurrentSection  string
	NavigationLinks interface{}
	Title           string

	UserID       string
	UserVerified bool

	Data map[string]interface{}
}

func (c *Context) Add(key string, value interface{}) {
	if c.Data == nil {
		c.Data = make(map[string]interface{})
	}
	c.Data[key] = value
}
