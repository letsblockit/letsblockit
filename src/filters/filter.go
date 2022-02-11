package filters

var (
	PresetNameSeparator = "---preset---"
	filenameSuffix      = ".yaml"
	yamlSeparator       = []byte("\n---")
	newLine             = []byte("\n")
)

type filter interface {
	finishParsing(string)
}

type Preset struct {
	Name        string   `validate:"required"`
	Description string   `validate:"required"`
	Source      string   `validate:"omitempty,url"`
	Values      []string `validate:"required"`
	Default     bool
}

type Filter struct {
	Name        string        `validate:"required"`
	Title       string        `validate:"required"`
	Params      []FilterParam `validate:"dive"`
	Tags        []string      `validate:"dive,alphaunicode"`
	Template    string        `validate:"required"`
	Description string        `validate:"required"`
	presets     []presetEntry // Generated on parse from params and presets
}

type presetEntry struct {
	EnableKey string
	Name      string
	TargetKey string
	Value     interface{}
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
	OnlyIf      string      `validate:"omitempty,valid_only_if"`
	Presets     []Preset    `validate:"dive"`
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

func (f *filterAndTests) finishParsing(desc string) {
	f.Filter.finishParsing(desc)
}

func (f *Filter) finishParsing(desc string) {
	f.Description = desc
	for _, param := range f.Params {
		if param.Type != StringListParam {
			continue
		}
		for _, preset := range param.Presets {
			f.presets = append(f.presets, presetEntry{
				EnableKey: param.BuildPresetParamName(preset.Name),
				Name:      preset.Name,
				TargetKey: param.Name,
				Value:     preset.Values,
			})
		}
	}
}

func (f *Filter) HasTag(tag string) bool {
	for _, t := range f.Tags {
		if t == tag {
			return true
		}
	}
	return false
}

func (p *FilterParam) BuildPresetParamName(preset string) string {
	return p.Name + PresetNameSeparator + preset
}
