package server

import (
	"embed"
	"io"

	"github.com/aymerick/raymond"
	"github.com/xvello/weblock/utils"
)

//go:embed templates
var templateFiles embed.FS

// templates holds parsed web templates
type templates map[string]*raymond.Template

// loadTemplates parses all web templates found in the templates folder
func loadTemplates() (templates, error) {
	repo := make(templates)
	err := utils.Walk(templateFiles, ".handlebars", func(name string, file io.Reader) error {
		contents, e := io.ReadAll(file)
		if e != nil {
			return e
		}
		repo[name], e = raymond.Parse(string(contents))
		return e
	})

	return repo, err
}
