package server

import (
	"embed"
	"io"
	"net/http"
	"strings"

	"github.com/aymerick/raymond"
	"github.com/labstack/echo/v4"
	"github.com/xvello/weblock/utils"
)

//go:embed templates/*
var templateFiles embed.FS

// templates holds parsed web templates
type templates struct {
	templates map[string]*raymond.Template
}

// loadTemplates parses all web templates found in the templates folder
func loadTemplates(helpers map[string]interface{}) (*templates, error) {
	// Parse toplevel layout template
	contents, err := templateFiles.ReadFile("templates/_layout.handlebars")
	if err != nil {
		return nil, err
	}
	layout, err := raymond.Parse(string(contents))
	if err != nil {
		return nil, err
	}
	if helpers != nil {
		layout.RegisterHelpers(helpers)
	}

	// Parse pages
	repo := templates{make(map[string]*raymond.Template)}
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
		page := layout.Clone()
		page.RegisterPartialTemplate("content", partial)
		repo.templates[name] = page
		return e
	})

	return &repo, err
}

func (t *templates) render(c echo.Context, name string, data interface{}) error {
	tpl, found := t.templates[name]
	if !found {
		return echo.NewHTTPError(http.StatusNotFound, "template %s not found", name)
	}
	contents, err := tpl.Exec(data)
	if err != nil {
		return err
	}
	return c.HTML(http.StatusOK, contents)
}
