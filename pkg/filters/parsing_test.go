package filters

import (
	"os"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseFilter(t *testing.T) {
	file, err := os.Open("testdata/filter.yaml")
	require.NoError(t, err)
	defer file.Close()

	filter, err := ParseFilter("named", file)
	require.NoError(t, err)

	assert.EqualValues(t, &Filter{
		Name:  "named",
		Title: "Filter title",
		Params: map[string]FilterParam{
			"boolean_param": {
				Description: "A boolean parameter",
				Type:        BooleanParam,
				Default:     true,
			},
			"another_boolean": {
				Description: "A disabled boolean parameter",
				Type:        BooleanParam,
				Default:     false,
			},
			"string_param": {
				Description: "A string parameter",
				Type:        StringParam,
				Default:     "René Coty",
			},
			"string_list": {
				Description: "A list of strings",
				Type:        StringListParam,
				Default:     []interface{}{"abc", "123"},
			},
		},
		Template: `{{#each string_list}}
{{ . }}
{{/each}}
`,
		Description: []byte("<h2>Test description title</h2>\n"),
	}, filter)
}

func TestParseFilterAndTest(t *testing.T) {
	file, err := os.Open("testdata/filter.yaml")
	require.NoError(t, err)
	defer file.Close()

	filter, err := parseFilterAndTest("named", file)
	require.NoError(t, err)

	assert.EqualValues(t, &filterAndTests{
		Filter: Filter{
			Name:  "named",
			Title: "Filter title",
			Params: map[string]FilterParam{
				"boolean_param": {
					Description: "A boolean parameter",
					Type:        BooleanParam,
					Default:     true,
				},
				"another_boolean": {
					Description: "A disabled boolean parameter",
					Type:        BooleanParam,
					Default:     false,
				},
				"string_param": {
					Description: "A string parameter",
					Type:        StringParam,
					Default:     "René Coty",
				},
				"string_list": {
					Description: "A list of strings",
					Type:        StringListParam,
					Default:     []interface{}{"abc", "123"},
				},
			},
			Template:    "{{#each string_list}}\n{{ . }}\n{{/each}}\n",
			Description: []byte("<h2>Test description title</h2>\n"),
		},
		Tests: []testCase{{
			Params: map[string]interface{}{
				"boolean_param": true,
				"string_param":  "ignored",
				"string_list":   []interface{}{"one", "two", "three"},
			},
			Output: "one\ntwo\nthree\n",
		}},
	}, filter)
}

type vErrs map[string]string

func TestValidateFilter(t *testing.T) {
	tests := map[string]struct {
		input filter
		err   vErrs
	}{
		"simple_ok": {
			input: &Filter{
				Name:        "name",
				Title:       "title",
				Template:    "template",
				Description: []byte("desc"),
			},
			err: nil,
		},
		"all_missing": {
			input: &Filter{},
			err: vErrs{
				"Filter.Name":        "required",
				"Filter.Title":       "required",
				"Filter.Template":    "required",
				"Filter.Description": "required",
			},
		},
		"param_ok": {
			input: &Filter{
				Name:        "name",
				Title:       "title",
				Template:    "template",
				Description: []byte("desc"),
				Params: map[string]FilterParam{
					"param": {
						Description: "desc",
						Type:        "checkbox",
						Default:     true,
					},
				},
			},
		},
		"param_bad_type": {
			input: &Filter{
				Name:        "name",
				Title:       "title",
				Template:    "template",
				Description: []byte("desc"),
				Params: map[string]FilterParam{
					"param": {
						Description: "desc",
						Type:        "bad",
						Default:     true,
					},
				},
			},
			err: vErrs{
				"Filter.Params[param].Type": "oneof",
			},
		},
		"param_empty": {
			input: &Filter{
				Name:        "name",
				Title:       "title",
				Template:    "template",
				Description: []byte("desc"),
				Params: map[string]FilterParam{
					"param": {},
				},
			},
			err: vErrs{
				"Filter.Params[param].Description": "required",
				"Filter.Params[param].Type":        "required",
			},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			if tc.err == nil {
				assert.NoError(t, validate.Struct(tc.input))
			} else {
				err := validate.Struct(tc.input)
				require.Error(t, err)
				errs := make(vErrs)
				for _, e := range err.(validator.ValidationErrors) {
					errs[e.StructNamespace()] = e.Tag()
				}
				assert.EqualValues(t, tc.err, errs)
			}
		})
	}
}
