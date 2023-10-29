package filters

import (
	"fmt"
	"io"
	"io/fs"
	"sort"
	"strings"

	"github.com/imantung/mario"
	"github.com/letsblockit/letsblockit/data"
	"github.com/samber/lo"
)

const (
	CustomRulesFilterName = "custom-rules"
	CustomTagName         = "custom"
)

// Repository holds parsed Templates ready for use
type Repository struct {
	main         *mario.Template
	templateMap  map[string]*Template
	templateList []*Template
	tagList      []string
}

// Load parses template definitions from the given filesystem
func Load(templates, presets fs.FS) (*Repository, error) {
	main, err := mario.New().Parse("{{>(_template)}}")
	main.WithHelperFunc("string_split", func(args string) []string {
		return strings.Split(args, " ")
	})
	if err != nil {
		return nil, fmt.Errorf("failed to parse toplevel template: %w", err)
	}
	repo := &Repository{
		main:        main,
		templateMap: make(map[string]*Template),
	}
	allTags := make(map[string]struct{})

	err = data.Walk(templates, filenameSuffix, func(name string, file io.Reader) error {
		tpl, e := parseTemplate(name, file)
		if e != nil {
			return e
		}
		if e = parsePresets(tpl, presets); err != nil {
			return e
		}
		partial, e := mario.New().Parse(tpl.Template)
		if e != nil {
			return fmt.Errorf("failed to parse template template: %w", e)
		}
		_ = main.WithPartial(name, partial)
		repo.templateMap[name] = tpl
		repo.templateList = append(repo.templateList, tpl)
		for _, tag := range tpl.Tags {
			allTags[tag] = struct{}{}
		}
		return nil
	})
	sortTemplates(repo.templateList)
	repo.tagList = flattenTagMap(allTags)

	return repo, err
}

func (r *Repository) Get(name string) (*Template, error) {
	tpl, found := r.templateMap[name]
	if !found {
		return nil, fmt.Errorf("unknown template '%s'", name)
	}
	return tpl, nil
}

func (r *Repository) GetAll() []*Template {
	return r.templateList
}

func (r *Repository) GetTags() []string {
	return r.tagList
}

func (r *Repository) Has(name string) bool {
	_, found := r.templateMap[name]
	return found
}

func (r *Repository) Render(w io.Writer, instance *Instance) error {
	tpl, found := r.templateMap[instance.Template]
	if !found {
		return fmt.Errorf("template '%s' not found", instance.Template)
	}
	params := shallowCopy(instance.Params)
	params["_template"] = instance.Template

	if instance.TestMode {
		w = NewTestModeTransformer(w)
	}
	if tpl.rawRules {
		return tpl.renderRawRules(w, params)
	}
	if err := r.main.Execute(w, params); err != nil {
		return err
	}

	for _, preset := range tpl.presets {
		if params[preset.EnableKey] == true {
			if _, err := fmt.Fprintln(w, preset.Header); err != nil {
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

// sortTemplates moves custom templates at the end of the list, keeping the other templates in alphabetical order
func sortTemplates(filters []*Template) {
	sort.Slice(filters, func(i, j int) bool {
		iCustom := lo.Contains(filters[i].Tags, CustomTagName)
		jCustom := lo.Contains(filters[j].Tags, CustomTagName)
		if iCustom == jCustom {
			return strings.Compare(filters[i].Title, filters[j].Title) < 0
		}
		return jCustom
	})
}

func shallowCopy(input map[string]interface{}) map[string]interface{} {
	output := make(map[string]interface{}, len(input))
	for k, v := range input {
		output[k] = v
	}
	return output
}
