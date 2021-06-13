package filters

import (
	"embed"

	"github.com/aymerick/raymond"
)

//go:embed data
var inputFiles embed.FS

var filenameSuffix = ".yaml"
var yamlSeparator = []byte("\n---")
var newLine = []byte("\n")

type filter interface {
	SetDescription([]byte)
	Parse() error
}

type Filter struct {
	Name        string                 `validate:"required"`
	Title       string                 `validate:"required"`
	Params      map[string]FilterParam `validate:"dive"`
	Template    string                 `validate:"required"`
	Description []byte                 `validate:"required"` // Rendered HTML bytes
	Parsed      *raymond.Template
}

type filterAndTests struct {
	Filter `yaml:"a,inline"`
	Tests  []testCase
}

type FilterParam struct {
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

func (f *Filter) SetDescription(desc []byte) {
	f.Description = desc
}

func (f *Filter) Parse() error {
	var err error
	f.Parsed, err = raymond.Parse(f.Template)
	return err
}

func (f *filterAndTests) SetDescription(desc []byte) {
	f.Filter.SetDescription(desc)
}

func (f *filterAndTests) Parse() error {
	return f.Filter.Parse()
}
