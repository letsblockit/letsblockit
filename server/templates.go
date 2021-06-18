package server

import (
	"embed"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/aymerick/raymond"
	"github.com/gin-gonic/gin"
	"github.com/xvello/weblock/utils"
)

//go:embed templates/*
var templateFiles embed.FS

// templates holds parsed web templates
type templates struct {
	templates map[string]*raymond.Template
}

// loadTemplates parses all web templates found in the templates folder
func loadTemplates() (*templates, error) {
	// Parse toplevel layout template
	contents, err := templateFiles.ReadFile("templates/_layout.handlebars")
	if err != nil {
		return nil, err
	}
	layout, err := raymond.Parse(string(contents))
	if err != nil {
		return nil, err
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

func (t *templates) render(c *gin.Context, name string, ctx interface{}) {
	if tpl, found := t.templates[name]; found {
		contents, err := tpl.Exec(ctx)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
		} else {
			c.Header("Content-Type", "text/html; charset=utf-8")
			c.String(http.StatusOK, contents)
		}
	} else {
		c.AbortWithError(http.StatusNotFound, fmt.Errorf("template %s not found", name))
	}
}
