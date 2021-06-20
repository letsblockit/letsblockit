package server

import (
	"embed"
	"io"
	"net/http"
	"strings"

	"github.com/aymerick/raymond"
	"github.com/labstack/echo/v4"
	"github.com/russross/blackfriday/v2"
	"github.com/xvello/weblock/utils"
)

//go:embed templates/*
var templateFiles embed.FS

type page struct {
	Partial  string
	Contents string
}

// templates holds parsed pages ready for rendering
type templates struct {
	main  *raymond.Template
	pages map[string]*page
}

// loadTemplates parses all web templates found in the templates folder
func loadTemplates(helpers map[string]interface{}) (*templates, error) {
	tpl := templates{
		pages: make(map[string]*page),
	}
	// Parse toplevel layout template
	contents, err := templateFiles.ReadFile("templates/_layout.handlebars")
	if err != nil {
		return nil, err
	}
	tpl.main, err = raymond.Parse(string(contents))
	if err != nil {
		return nil, err
	}
	if helpers != nil {
		tpl.main.RegisterHelpers(helpers)
	}

	// Parse handlebars templates
	err = utils.Walk(templateFiles, ".handlebars", func(name string, file io.Reader) error {
		if strings.HasPrefix(name, "_") {
			return nil
		}
		contents, e := io.ReadAll(file)
		if e != nil {
			return e
		}
		partial, e := raymond.Parse(string(contents))
		if e != nil {
			return e
		}
		tpl.main.RegisterPartialTemplate(name, partial)
		tpl.pages[name] = &page{Partial: name}
		return e
	})
	if err != nil {
		return nil, err
	}

	// Parse markdown pages
	err = utils.Walk(templateFiles, ".md", func(name string, file io.Reader) error {
		if strings.HasPrefix(name, "_") {
			return nil
		}
		rawContents, e := io.ReadAll(file)
		if e != nil {
			return e
		}
		tpl.pages[name] = &page{Contents: string(blackfriday.Run(rawContents))}
		return e
	})

	return &tpl, err
}

func (t *templates) render(c echo.Context, name string, ctx map[string]interface{}) error {
	var found bool
	ctx["_page"], found = t.pages[name]
	if !found {
		return echo.NewHTTPError(http.StatusNotFound, "template %s not found", name)
	}
	contents, err := t.main.Exec(ctx)
	if err != nil {
		return err
	}
	return c.HTML(http.StatusOK, contents)
}
