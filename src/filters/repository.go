package filters

import (
	"context"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"strings"

	"github.com/imantung/mario"
	"github.com/labstack/echo/v4"
	"github.com/xvello/letsblockit/data"
)

// Repository holds parsed Filters ready for use
type Repository struct {
	main  *mario.Template
	fMap  map[string]*Filter
	fList []*Filter
}

// LoadFilters parses all filter definitions found in the data folder
func LoadFilters() (*Repository, error) {
	return load(data.Filters)
}

func load(input fs.FS) (*Repository, error) {
	main, err := mario.New().Parse("{{>(_filter)}}")
	main.WithHelperFunc("string_split", func(args string) []string {
		return strings.Split(args, " ")
	})
	repo := &Repository{
		main: main,
		fMap: make(map[string]*Filter),
	}
	if err != nil {
		return nil, fmt.Errorf("failed to parse toplevel template: %w", err)
	}

	err = data.Walk(input, filenameSuffix, func(name string, file io.Reader) error {
		f, e := parseFilter(name, file)
		if e != nil {
			return e
		}
		partial, e := mario.New().Parse(f.Template)
		if e != nil {
			return fmt.Errorf("failed to parse filter template: %w", err)
		}
		_ = main.WithPartial(name, partial)
		repo.fMap[name] = f
		repo.fList = append(repo.fList, f) // list is naturally sorted because Walkdir iterates on lexical order
		return nil
	})

	return repo, err
}

func (r *Repository) GetFilter(name string) (*Filter, error) {
	filter, found := r.fMap[name]
	if !found {
		return nil, fmt.Errorf("unknown filter %s", name)
	}
	return filter, nil
}

func (r *Repository) GetFilters() []*Filter {
	return r.fList
}

func (r *Repository) Render(ctx context.Context, w io.Writer, name string, data map[string]interface{}) error {
	_, found := r.fMap[name]
	if !found {
		return echo.NewHTTPError(http.StatusNotFound, "template %s not found", name)
	}
	data["_filter"] = name
	return r.main.Execute(w, data)
}
