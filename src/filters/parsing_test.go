package filters

import (
	"os"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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
	err = validate.RegisterValidation("preset_allowed", func(fl validator.FieldLevel) bool {
		paramType := ParamType(fl.Parent().FieldByName("Type").String())
		return paramType == StringListParam
	})
	require.NoError(t, err)
	err = validate.RegisterValidation("valid_only_if", func(fl validator.FieldLevel) bool {
		target := fl.Field().String()
		if fl.Parent().FieldByName("Name").String() == target {
			return false // Referencing itself
		}
		if filter, ok := (fl.Top().Interface()).(*Template); ok {
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

func TestParseTemplate(t *testing.T) {
	file, err := os.Open("testdata/templates/simple.yaml")
	require.NoError(t, err)
	defer file.Close()

	filter, err := parseTemplate("simple", file)
	require.NoError(t, err)

	expectedTemplate := Template{
		Name:  "simple",
		Title: "Template title",
		Params: []Parameter{
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
		Tags:     []string{"tag1", "tag2"},
		Template: "{{#each string_list}}\n{{ . }}\n{{/each}}\n",
		Tests: []testCase{{
			Params: map[string]interface{}{
				"boolean_param": true,
				"string_param":  "ignored",
				"string_list":   []interface{}{"one", "two", "three"},
			},
			Output: "one\ntwo\nthree\n",
		}},
		Description: "<h2>Test description title</h2>\n",
	}
	assert.EqualValues(t, &expectedTemplate, filter)
}

func TestParsePresets(t *testing.T) {
	tests := map[string]struct {
		input    *Template
		expected *Template
		err      vErrs
	}{
		"no_presets": {
			input:    &Template{Params: []Parameter{{Name: "param name"}}},
			expected: &Template{Params: []Parameter{{Name: "param name"}}},
			err:      nil,
		},
		"with_presets": {
			input: &Template{
				Name: "simple-template",
				Params: []Parameter{{
					Name: "param-name",
					Type: StringListParam,
					Presets: []Preset{{
						Name:        "internal-preset",
						Description: "preset description",
						Values:      []string{"1", "2"},
						Default:     true,
					}, {
						Name:        "sourced-preset",
						Description: "preset description",
						Source:      "preset source",
						License:     "preset license",
						Default:     false,
					}},
				}},
			},
			expected: &Template{
				Name: "simple-template",
				Params: []Parameter{{
					Name: "param-name",
					Type: StringListParam,
					Presets: []Preset{{
						Name:        "internal-preset",
						Description: "preset description",
						Values:      []string{"1", "2"},
						Default:     true,
					}, {
						Name:        "sourced-preset",
						Description: "preset description",
						Source:      "preset source",
						License:     "preset license",
						Values:      []string{"a", "b"},
						Default:     false,
					}},
				}},
				presets: []presetEntry{{
					EnableKey: "param-name---preset---internal-preset",
					Name:      "internal-preset",
					Value:     []string{"1", "2"},
					TargetKey: "param-name",
					Header:    "!! simple-template with internal-preset preset",
				}, {
					EnableKey: "param-name---preset---sourced-preset",
					Name:      "sourced-preset",
					Value:     []string{"a", "b"},
					TargetKey: "param-name",
					Header: `!! simple-template with sourced-preset preset
!! Source: preset source
!! License: preset license`,
				}},
			},
			err: nil,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			require.NoError(t, parsePresets(tc.input, testPresets))
			require.EqualValues(t, tc.expected, tc.input)
		})
	}
}

type vErrs map[string]string

func TestValidateTemplate(t *testing.T) {
	tests := map[string]struct {
		input *Template
		err   vErrs
	}{
		"simple_ok": {
			input: &Template{
				Name:        "name",
				Title:       "title",
				Template:    "template",
				Description: "desc",
			},
			err: nil,
		},
		"all_missing": {
			input: &Template{},
			err: vErrs{
				"Template.Name":        "required",
				"Template.Title":       "required",
				"Template.Template":    "required",
				"Template.Description": "required",
			},
		},
		"param_ok": {
			input: &Template{
				Name:        "name",
				Title:       "title",
				Template:    "template",
				Description: "desc",
				Params: []Parameter{
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
			input: &Template{
				Name:        "name",
				Title:       "title",
				Template:    "template",
				Description: "desc",
				Params: []Parameter{
					{
						Name:        "param",
						Description: "desc",
						Type:        "bad",
						Default:     true,
					},
				},
			},
			err: vErrs{
				"Template.Params[0].Type":    "oneof",
				"Template.Params[0].Default": "valid_default",
			},
		},
		"param_empty": {
			input: &Template{
				Name:        "name",
				Title:       "title",
				Template:    "template",
				Description: "desc",
				Params: []Parameter{
					{},
				},
			},
			err: vErrs{
				"Template.Params[0].Name":        "required",
				"Template.Params[0].Description": "required",
				"Template.Params[0].Type":        "required",
				"Template.Params[0].Default":     "valid_default",
			},
		},
		"onlyif_bad_type": {
			input: &Template{
				Name:        "name",
				Title:       "title",
				Template:    "template",
				Description: "desc",
				Params: []Parameter{
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
				"Template.Params[1].OnlyIf": "valid_only_if",
			},
		},
		"onlyif_unknown": {
			input: &Template{
				Name:        "name",
				Title:       "title",
				Template:    "template",
				Description: "desc",
				Params: []Parameter{
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
				"Template.Params[0].OnlyIf": "valid_only_if",
			},
		},
		"onlyif_self": {
			input: &Template{
				Name:        "name",
				Title:       "title",
				Template:    "template",
				Description: "desc",
				Params: []Parameter{
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
				"Template.Params[0].OnlyIf": "valid_only_if",
			},
		},
		"invalid_tags": {
			input: &Template{
				Name:        "name",
				Title:       "title",
				Template:    "template",
				Description: "desc",
				Tags:        []string{"abc", "%%"},
			},
			err: vErrs{
				"Template.Tags[1]": "alphaunicode",
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
