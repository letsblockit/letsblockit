package filters

import (
	"os"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var expectedFilter = Filter{
	Name:  "simple",
	Title: "Filter title",
	Params: []FilterParam{
		{
			Name:        "boolean_param",
			Description: "A boolean parameter",
			Type:        BooleanParam,
			Default:     true,
		},
		{
			Name:        "another_boolean",
			Description: "A disabled boolean parameter",
			Type:        BooleanParam,
			Default:     false,
		},
		{
			Name:        "string_param",
			Description: "A string parameter",
			Type:        StringParam,
			Default:     "René Coty",
		},
		{
			Name:        "string_list",
			Description: "A list of strings",
			Type:        StringListParam,
			Default:     []interface{}{"abc", "123"},
		},
	},
	Tags:        []string{"tag1", "tag2"},
	Template:    "{{#each string_list}}\n{{ . }}\n{{/each}}\n",
	Description: "<h2>Test description title</h2>\n",
}

func buildValidator(t *testing.T) *validator.Validate {
	var validate = validator.New()
	err := validate.RegisterValidation("valid_default", func(fl validator.FieldLevel) bool {
		paramType := ParamType(fl.Parent().FieldByName("Type").String())
		switch fl.Field().Interface().(type) {
		case bool:
			return paramType == BooleanParam
		case string:
			return paramType == StringParam || paramType == MultiLineParam
		case []string, []interface{}:
			return paramType == StringListParam
		default:
			return false
		}
	})
	require.NoError(t, err)
	err = validate.RegisterValidation("valid_only_if", func(fl validator.FieldLevel) bool {
		target := fl.Field().String()
		if fl.Parent().FieldByName("Name").String() == target {
			return false // Referencing itself
		}
		if filter, ok := (fl.Top().Interface()).(*Filter); ok {
			for _, p := range filter.Params {
				if p.Name == target && p.Type == BooleanParam {
					return true
				}
			}
		}
		return false
	})
	require.NoError(t, err)
	return validate
}

func TestParseFilter(t *testing.T) {
	file, err := os.Open("testdata/simple.yaml")
	require.NoError(t, err)
	defer file.Close()

	filter, err := parseFilter("simple", file)
	require.NoError(t, err)

	assert.EqualValues(t, &expectedFilter, filter)
}

func TestParseFilterAndTest(t *testing.T) {
	file, err := os.Open("testdata/simple.yaml")
	require.NoError(t, err)
	defer file.Close()

	filter, err := parseFilterAndTest("simple", file)
	require.NoError(t, err)

	assert.EqualValues(t, &FilterAndTests{
		Filter: expectedFilter,
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
				Description: "desc",
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
				Description: "desc",
				Params: []FilterParam{
					{
						Name:        "param1",
						Description: "desc",
						Type:        BooleanParam,
						Default:     true,
					},
					{
						Name:        "param2",
						Description: "desc",
						Type:        StringParam,
						Default:     "example",
					},
					{
						Name:        "param3",
						Description: "desc",
						OnlyIf:      "param1",
						Type:        StringListParam,
						Default:     []string{"abc", "123"},
					},
				},
			},
		},
		"param_bad_type": {
			input: &Filter{
				Name:        "name",
				Title:       "title",
				Template:    "template",
				Description: "desc",
				Params: []FilterParam{
					{
						Name:        "param",
						Description: "desc",
						Type:        "bad",
						Default:     true,
					},
				},
			},
			err: vErrs{
				"Filter.Params[0].Type":    "oneof",
				"Filter.Params[0].Default": "valid_default",
			},
		},
		"param_empty": {
			input: &Filter{
				Name:        "name",
				Title:       "title",
				Template:    "template",
				Description: "desc",
				Params: []FilterParam{
					{},
				},
			},
			err: vErrs{
				"Filter.Params[0].Name":        "required",
				"Filter.Params[0].Description": "required",
				"Filter.Params[0].Type":        "required",
				"Filter.Params[0].Default":     "valid_default",
			},
		},
		"onlyif_bad_type": {
			input: &Filter{
				Name:        "name",
				Title:       "title",
				Template:    "template",
				Description: "desc",
				Params: []FilterParam{
					{
						Name:        "param1",
						Description: "desc",
						Type:        StringParam,
						Default:     "example",
					},
					{
						Name:        "param2",
						Description: "desc",
						OnlyIf:      "param1",
						Type:        StringListParam,
						Default:     []string{"abc", "123"},
					},
				},
			},
			err: vErrs{
				"Filter.Params[1].OnlyIf": "valid_only_if",
			},
		},
		"onlyif_unknown": {
			input: &Filter{
				Name:        "name",
				Title:       "title",
				Template:    "template",
				Description: "desc",
				Params: []FilterParam{
					{
						Name:        "param3",
						Description: "desc",
						OnlyIf:      "param1",
						Type:        StringListParam,
						Default:     []string{"abc", "123"},
					},
				},
			},
			err: vErrs{
				"Filter.Params[0].OnlyIf": "valid_only_if",
			},
		},
		"onlyif_self": {
			input: &Filter{
				Name:        "name",
				Title:       "title",
				Template:    "template",
				Description: "desc",
				Params: []FilterParam{
					{
						Name:        "param3",
						Description: "desc",
						OnlyIf:      "param3",
						Type:        StringListParam,
						Default:     []string{"abc", "123"},
					},
				},
			},
			err: vErrs{
				"Filter.Params[0].OnlyIf": "valid_only_if",
			},
		},
		"invalid_tags": {
			input: &Filter{
				Name:        "name",
				Title:       "title",
				Template:    "template",
				Description: "desc",
				Tags:        []string{"abc", "%%"},
			},
			err: vErrs{
				"Filter.Tags[1]": "alphaunicode",
			},
		},
	}

	validate := buildValidator(t)
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
