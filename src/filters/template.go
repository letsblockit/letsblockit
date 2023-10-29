package filters

import (
	"fmt"
	"io"
)

var (
	presetNameSeparator = "---preset---"
	filenameSuffix      = ".yaml"
	yamlSeparator       = []byte("\n---")
	newLine             = []byte("\n")
)

type Preset struct {
	Name        string   `validate:"required"`
	Description string   `validate:"required"`
	Source      string   `validate:"omitempty,url" yaml:",omitempty"`
	License     string   `validate:"required_with=Source" yaml:",omitempty"`
	Values      []string `validate:"required_without=Source"`
	Default     bool     `yaml:",omitempty"`
}

type Template struct {
	Name         string      `validate:"required" yaml:"-"`
	Title        string      `validate:"required"`
	Params       []Parameter `validate:"dive" yaml:",omitempty"`
	Tags         []string    `validate:"dive,alphaunicode" yaml:",omitempty"`
	Template     string      `validate:"required_without=rawRules,excluded_with=rawRules"`
	Tests        []testCase
	Description  string `validate:"required" yaml:"-"`
	Contributors []string
	Sponsors     []string
	presets      []presetEntry `yaml:"-"` // Generated on parse from params and presets
	rawRules     bool
}

type presetEntry struct {
	EnableKey string
	Name      string
	TargetKey string
	Header    string
	Value     interface{}
}

type Parameter struct {
	Name        string      `validate:"required"`
	Description string      `validate:"required"`
	Link        string      `validate:"omitempty,url" yaml:",omitempty"`
	Type        ParamType   `validate:"required,oneof=checkbox string list multiline"`
	OnlyIf      string      `validate:"omitempty,valid_only_if" yaml:",omitempty"`
	Default     interface{} `validate:"valid_default"`
	Rules       string      `validate:"omitempty,raw_rules_allowed" yaml:",omitempty"`
	Presets     []Preset    `validate:"omitempty,preset_allowed,dive" yaml:",omitempty"`
}

type ParamType string

const (
	BooleanParam    ParamType = "checkbox"
	StringParam     ParamType = "string"
	StringListParam ParamType = "list"
	MultiLineParam  ParamType = "multiline"
)

type testCase struct {
	Params map[string]interface{} `yaml:",omitempty"`
	Output string                 `validate:"required"`
}

func (f *Template) HasTag(tag string) bool {
	for _, t := range f.Tags {
		if t == tag {
			return true
		}
	}
	return false
}

func (f *Template) renderRawRules(w io.Writer, params map[string]interface{}) error {
	for _, param := range f.Params {
		if param.Type != BooleanParam {
			return fmt.Errorf("unsupported param type %s for %s", param.Type, param.Name)
		}
		if len(param.Rules) == 0 {
			return fmt.Errorf("no rules for param %s", param.Name)
		}
		if value, found := params[param.Name]; found && value == true {
			if _, err := fmt.Fprint(w, param.Rules); err != nil {
				return err
			}
		}
	}
	return nil
}

func (p *Parameter) BuildPresetParamName(preset string) string {
	return p.Name + presetNameSeparator + preset
}
