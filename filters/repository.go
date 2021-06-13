package filters

import (
	"fmt"
	"io/fs"
	"strings"
)

type Repository struct {
	filters map[string]*Filter
}

func LoadFilters() (*Repository, error) {
	repo := &Repository{
		filters: make(map[string]*Filter),
	}

	err := fs.WalkDir(inputFiles, "data", func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() || !strings.HasSuffix(d.Name(), filenameSuffix) {
			return nil
		}
		name := strings.TrimSuffix(d.Name(), filenameSuffix)
		file, err := inputFiles.Open(path)
		if err != nil {
			return fmt.Errorf("cannot open %s: %w", path, err)
		}
		repo.filters[name], err = parseFilter(name, file)
		_ = file.Close()
		if err != nil {
			return fmt.Errorf("cannot parse %s: %w", path, err)
		}
		return nil
	})

	return repo, err
}

func (r *Repository) Render(name string, data interface{}) (string, error) {
	filter, found := r.filters[name]
	if !found {
		return "", fmt.Errorf("unknown filter %s", name)
	}
	return filter.Parsed.Exec(data)
}
