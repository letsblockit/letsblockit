package pages

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/imantung/mario"
	"github.com/labstack/echo/v4"
	"github.com/letsblockit/letsblockit/data"
	"github.com/russross/blackfriday/v2"
)

type page struct {
	Partial  string
	Contents string
}

type ContextBuilder func(c echo.Context, title string) *Context

// Pages holds parsed pages ready for rendering
type Pages struct {
	main    *mario.Template
	naked   *mario.Template
	pages   map[string]*page
	builder ContextBuilder
}

// LoadPages parses all web pages found in the pages folder
func LoadPages() (*Pages, error) {
	pp := Pages{
		pages: make(map[string]*page),
	}
	// Parse toplevel layout template
	contents, err := data.Pages.ReadFile("pages/_layout.hbs")
	if err != nil {
		return nil, err
	}
	pp.main, err = mario.New().Parse(string(contents))
	if err != nil {
		return nil, fmt.Errorf("failed to parse toplevel template: %w", err)
	}

	// Parse toplevel naked template
	contents, err = data.Pages.ReadFile("pages/_naked.hbs")
	if err != nil {
		return nil, err
	}
	pp.naked, err = mario.New().Parse(string(contents))
	if err != nil {
		return nil, fmt.Errorf("failed to parse naked template: %w", err)
	}

	// Parse handlebars pages
	err = data.Walk(data.Pages, ".hbs", func(name string, file io.Reader) error {
		if strings.HasPrefix(name, "_") {
			return nil
		}
		contents, e := io.ReadAll(file)
		if e != nil {
			return e
		}
		partial, e := mario.New().Parse(string(contents))
		if e != nil {
			return e
		}
		_ = pp.main.WithPartial(name, partial)
		_ = pp.naked.WithPartial(name, partial)

		pp.pages[name] = &page{Partial: name}
		return e
	})
	if err != nil {
		return nil, err
	}

	// Parse markdown pages
	err = data.Walk(data.Pages, ".md", func(name string, file io.Reader) error {
		if strings.HasPrefix(name, "_") {
			return nil
		}
		rawContents, e := io.ReadAll(file)
		if e != nil {
			return e
		}
		pp.pages[name] = &page{Contents: string(blackfriday.Run(rawContents))}
		return e
	})
	if err != nil {
		return nil, err
	}

	// Load raw html pages
	err = data.Walk(data.Pages, ".html", func(name string, file io.Reader) error {
		if strings.HasPrefix(name, "_") {
			return nil
		}
		rawContents, e := io.ReadAll(file)
		if e != nil {
			return e
		}
		pp.pages[name] = &page{Contents: string(rawContents)}
		return e
	})

	return &pp, err
}

func (p *Pages) RegisterHelpers(helpers map[string]interface{}) {
	for n, h := range helpers {
		_ = p.main.WithHelperFunc(n, h)
	}
}

func (p *Pages) RegisterContextBuilder(b ContextBuilder) {
	p.builder = b
}

func (p *Pages) BuildPageContext(c echo.Context, title string) *Context {
	if p.builder == nil {
		return nil
	}
	return p.builder(c, title)
}

func (p *Pages) Render(c echo.Context, name string, data *Context) error {
	var found bool
	data.Page, found = p.pages[name]
	if !found {
		return echo.NewHTTPError(http.StatusNotFound, "template not found: "+name)
	}
	tpl := p.main
	if data.NakedContent {
		tpl = p.naked
	}
	buf := new(bytes.Buffer)
	if err := tpl.Execute(buf, data); err != nil {
		return err
	}
	return c.HTMLBlob(http.StatusOK, buf.Bytes())
}

func (p *Pages) RenderWithSidebar(c echo.Context, name, sidebar string, data *Context) error {
	var found bool
	data.Sidebar, found = p.pages[sidebar]
	if !found {
		return echo.NewHTTPError(http.StatusNotFound, "sidebar template not found: "+name)
	}
	return p.Render(c, name, data)
}

// RedirectToPage the user to another page, either via htmx client-side Redirect (form submissions)
// or http 302 Redirect (direct access, js disabled)
func (p *Pages) RedirectToPage(c echo.Context, name string, params ...interface{}) error {
	return p.Redirect(c, http.StatusFound, c.Echo().Reverse(name, params...))
}

func (p *Pages) Redirect(c echo.Context, code int, target string) error {
	if c.Request().Header.Get("HX-Request") == "true" {
		c.Response().Header().Set("HX-Redirect", target)
		return c.NoContent(200)
	}
	return c.Redirect(code, target)
}
