package filters

import (
	"fmt"
	"io"
	"io/fs"
	"sort"
	"strings"

	"github.com/imantung/mario"
	"github.com/letsblockit/letsblockit/data"
)

const CustomRulesFilterName = "custom-rules"

// Repository holds parsed Filters ready for use
type Repository struct {
	main       *mario.Template
	filterMap  map[string]*Filter
	filterList []*Filter
	tagList    []string
}

// LoadFilters parses all filter definitions in the given filesystem
func LoadFilters(input fs.FS) (*Repository, error) {
	main, err := mario.New().Parse("{{>(_filter)}}")
	main.WithHelperFunc("string_split", func(args string) []string {
		return strings.Split(args, " ")
	})
	if err != nil {
		return nil, fmt.Errorf("failed to parse toplevel template: %w", err)
	}
	repo := &Repository{
		main:      main,
		filterMap: make(map[string]*Filter),
	}
	allTags := make(map[string]struct{})

	err = data.Walk(input, filenameSuffix, func(name string, file io.Reader) error {
		f, e := parseFilter(name, file)
		if e != nil {
			return e
		}
		partial, e := mario.New().Parse(f.Template)
		if e != nil {
			return fmt.Errorf("failed to parse filter template: %w", e)
		}
		_ = main.WithPartial(name, partial)
		repo.filterMap[name] = f
		repo.filterList = append(repo.filterList, f)
		for _, tag := range f.Tags {
			allTags[tag] = struct{}{}
		}
		return nil
	})
	sortFilters(repo.filterList)
	repo.tagList = flattenTagMap(allTags)

	return repo, err
}

func (r *Repository) GetFilter(name string) (*Filter, error) {
	filter, found := r.filterMap[name]
	if !found {
		return nil, fmt.Errorf("unknown filter '%s'", name)
	}
	return filter, nil
}

func (r *Repository) GetFilters() []*Filter {
	return r.filterList
}

func (r *Repository) GetTags() []string {
	return r.tagList
}

func (r *Repository) Render(w io.Writer, instance *Instance) error {
	filter, found := r.filterMap[instance.Filter]
	if !found {
		return fmt.Errorf("template '%s' not found", instance.Filter)
	}
	params := shallowCopy(instance.Params)
	params["_filter"] = instance.Filter

	if instance.TestMode {
		w = NewTestModeTransformer(w)
	}
	if err := r.main.Execute(w, params); err != nil {
		return err
	}

	for _, preset := range filter.presets {
		if params[preset.EnableKey] == true {
			if _, err := fmt.Fprintln(w, "!!", instance.Filter, "with", preset.Name, "preset"); err != nil {
				return err
			}
			params := shallowCopy(params)
			params[preset.TargetKey] = preset.Value
			if err := r.main.Execute(w, params); err != nil {
				return err
			}
		}
	}

	return nil
}

func flattenTagMap(tags map[string]struct{}) []string {
	out := make([]string, 0, len(tags))
	for tag := range tags {
		out = append(out, tag)
	}
	sort.Strings(out)
	return out
}

// sortFilters moves custom-rules at the end of the list, keeping the other filters in alphabetical order
func sortFilters(filters []*Filter) {
	sort.Slice(filters, func(i, j int) bool {
		if filters[i].Name == CustomRulesFilterName {
			return false
		}
		if filters[j].Name == CustomRulesFilterName {
			return true
		}
		return strings.Compare(filters[i].Name, filters[j].Name) < 0
	})
}

func shallowCopy(input map[string]interface{}) map[string]interface{} {
	output := make(map[string]interface{}, len(input))
	for k, v := range input {
		output[k] = v
	}
	return output
}
