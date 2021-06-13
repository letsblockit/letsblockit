package filters

import (
	"fmt"
	"io/fs"
	"strings"

	"github.com/aymerick/raymond"
)

// Repository holds parsed Filters ready for use
type Repository struct {
	filters map[string]*Filter
}

// LoadFilters parses all filter definitions found in the data folder
func LoadFilters() (*Repository, error) {
	repo := &Repository{
		filters: make(map[string]*Filter),
	}

	err := fs.WalkDir(inputFiles, "data", func(path string, d fs.DirEntry, _ error) error {
		if d.IsDir() || !strings.HasSuffix(d.Name(), filenameSuffix) {
			return nil
		}
		name := strings.TrimSuffix(d.Name(), filenameSuffix)
		file, e := inputFiles.Open(path)
		if e != nil {
			return fmt.Errorf("cannot open %s: %w", path, e)
		}
		repo.filters[name], e = parseFilter(name, file)
		_ = file.Close()
		if e != nil {
			return fmt.Errorf("cannot parse %s: %w", path, e)
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
