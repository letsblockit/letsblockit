package filters

import (
	"embed"
	"fmt"
	"io"
	"io/fs"

	"github.com/xvello/weblock/utils"
)

//go:embed data
var definitionFiles embed.FS

// Repository holds parsed Filters ready for use
type Repository struct {
	fMap  map[string]*Filter
	fList []*Filter
}

// LoadFilters parses all filter definitions found in the data folder
func LoadFilters() (*Repository, error) {
	return load(definitionFiles)
}

func load(input fs.FS) (*Repository, error) {
	repo := &Repository{
		fMap: make(map[string]*Filter),
	}

	err := utils.Walk(input, filenameSuffix, func(name string, file io.Reader) error {
		f, e := parseFilter(name, file)
		if e != nil {
			return e
		}
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
