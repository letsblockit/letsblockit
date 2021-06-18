package filters

import (
	"github.com/aymerick/raymond"
)

var filenameSuffix = ".yaml"
var yamlSeparator = []byte("\n---")
var newLine = []byte("\n")

type filter interface {
	SetDescription(string)
	Parse() error
}

type Filter struct {
	Name        string        `validate:"required"`
	Title       string        `validate:"required"`
	Params      []FilterParam `validate:"dive"`
	Template    string        `validate:"required"`
	Description string        `validate:"required"`
	Parsed      *raymond.Template
}

type filterAndTests struct {
	Filter `yaml:"a,inline"`
	Tests  []testCase
}

type FilterParam struct {
	Name        string      `validate:"required"`
	Description string      `validate:"required"`
	Type        ParamType   `validate:"required,oneof=checkbox string list"`
	Default     interface{} `validate:"valid_default"`
}

type ParamType string

const (
	BooleanParam    ParamType = "checkbox"
	StringParam     ParamType = "string"
	StringListParam ParamType = "list"
)

type testCase struct {
	Params map[string]interface{}
	Output string `validate:"required"`
}

func (f *Filter) SetDescription(desc string) {
	f.Description = desc
}

func (f *Filter) Parse() error {
	var err error
	f.Parsed, err = raymond.Parse(f.Template)
	return err
}

func (f *filterAndTests) SetDescription(desc string) {
	f.Filter.SetDescription(desc)
}

func (f *filterAndTests) Parse() error {
	return f.Filter.Parse()
}
