package filters

import (
	"embed"
	"fmt"
	"io"
	"io/fs"

	"github.com/aymerick/raymond"
	"github.com/xvello/weblock/utils"
)

//go:embed data
var definitionFiles embed.FS

// Repository holds parsed Filters ready for use
type Repository struct {
	filters map[string]*Filter
}

// LoadFilters parses all filter definitions found in the data folder
func LoadFilters() (*Repository, error) {
	return load(definitionFiles)
}

func load(input fs.FS) (*Repository, error) {
	repo := &Repository{
		filters: make(map[string]*Filter),
	}

	err := utils.Walk(input, filenameSuffix, func(name string, file io.Reader) error {
		var e error
		repo.filters[name], e = parseFilter(name, file)
		if e != nil {
			return fmt.Errorf("cannot parse %s: %w", name, e)
		}
		return nil
	})

	return repo, err
}

// RenderFilter takes arguments and returns the result of filter
// templating for inclusion in an adblock filter list.
func (r *Repository) RenderFilter(name string, data interface{}) (string, error) {
	filter, err := r.getFilter(name)
	if err != nil {
		return "", err
	}
	return filter.Parsed.Exec(data)
}

// RenderPage renders the description page for a given filter.
// The passed template object will be given the full Filter object as input.
func (r *Repository) RenderPage(name string, template *raymond.Template) (string, error) {
	filter, err := r.getFilter(name)
	if err != nil {
		return "", err
	}
	return template.Exec(filter)
}

func (r *Repository) getFilter(name string) (*Filter, error) {
	filter, found := r.filters[name]
	if !found {
		return nil, fmt.Errorf("unknown filter %s", name)
	}
	return filter, nil
}
