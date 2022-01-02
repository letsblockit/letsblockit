package filters

var filenameSuffix = ".yaml"
var yamlSeparator = []byte("\n---")
var newLine = []byte("\n")

type filter interface {
	setDescription(string)
}

type Filter struct {
	Name        string        `validate:"required"`
	Blurb       string        `validate:"required"`
	Title       string        `validate:"required"`
	Params      []FilterParam `validate:"dive"`
	Tags        []string      `validate:"dive,alphaunicode"`
	Template    string        `validate:"required"`
	Description string        `validate:"required"`
}

type filterAndTests struct {
	Filter `yaml:"a,inline"`
	Tests  []testCase
}

type FilterParam struct {
	Name        string      `validate:"required"`
	Description string      `validate:"required"`
	Type        ParamType   `validate:"required,oneof=checkbox string list multiline"`
	Default     interface{} `validate:"valid_default"`
}

type ParamType string

const (
	BooleanParam    ParamType = "checkbox"
	StringParam     ParamType = "string"
	StringListParam ParamType = "list"
	MultiLineParam  ParamType = "multiline"
)

type testCase struct {
	Params map[string]interface{}
	Output string `validate:"required"`
}

func (f *Filter) setDescription(desc string) {
	f.Description = desc
}

func (f *filterAndTests) setDescription(desc string) {
	f.Filter.setDescription(desc)
}

func (f *Filter) HasTag(tag string) bool {
	for _, t := range f.Tags {
		if t == tag {
			return true
		}
	}
	return false
}
