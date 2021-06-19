package filters

import (
	"github.com/aymerick/raymond"
)

var filenameSuffix = ".yaml"
var yamlSeparator = []byte("\n---")
var newLine = []byte("\n")

type filter interface {
	setDescription(string)
	parse() error

}

type Filter struct {
	Name        string        `validate:"required"`
	Title       string        `validate:"required"`
	Params      []FilterParam `validate:"dive"`
	Tags        []string      `validate:"dive,alphaunicode"`
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

func (f *Filter) Render(data interface{}) (string, error) {
	return f.Parsed.Exec(data)
}

func (f *Filter) setDescription(desc string) {
	f.Description = desc
}

func (f *Filter) parse() error {
	var err error
	f.Parsed, err = raymond.Parse(f.Template)
	return err
}

func (f *filterAndTests) setDescription(desc string) {
	f.Filter.setDescription(desc)
}

func (f *filterAndTests) parse() error {
	return f.Filter.parse()
}
