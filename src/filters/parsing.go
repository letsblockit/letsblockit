package filters

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"strings"

	"github.com/russross/blackfriday/v2"
	"golang.org/x/exp/slices"
	"gopkg.in/yaml.v3"
)

const (
	newline                  string = "\n"
	presetFilePattern        string = "%s/%s.txt"
	presetHeaderPattern      string = "!! %s with %s preset"
	presetAttributionPattern string = `
!! Source: %s
!! License: %s`
)

func parseTemplate(name string, reader io.Reader) (*Template, error) {
	tpl := &Template{Name: name}

	// Read the whole input file and find the separator
	input, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	pos := bytes.Index(input, yamlSeparator)
	if pos < 0 {
		return nil, errors.New("separator not found")
	}

	// Parse the yaml
	err = yaml.Unmarshal(input[:pos+1], tpl)
	if err != nil {
		return nil, fmt.Errorf("invalid metadata: %w", err)
	}

	// Parse the markdown description
	pos += len(yamlSeparator)
	pos += bytes.Index(input[pos:], newLine)
	tpl.Description = string(blackfriday.Run(input[pos:]))

	// Check for presence of raw rules if template not given.
	// Must have either, but not both.
	hasTemplate := len(tpl.Template) > 0
	tpl.rawRules = len(tpl.Params) > 0
	for i, p := range tpl.Params {
		if p.Type == BooleanParam && len(p.Rules) > 0 {
			if hasTemplate {
				return nil, fmt.Errorf("%s has a template AND raw rules on %s, not allowed", tpl.Name, p.Name)
			}
			if !strings.HasSuffix(p.Rules, newline) {
				tpl.Params[i].Rules = p.Rules + newline
			}
		} else {
			if !hasTemplate {
				return nil, fmt.Errorf("%s has no template but param %s has no raw rules", tpl.Name, p.Name)
			}
			tpl.rawRules = false
			break
		}
	}

	// Make sure contributors and sponsors are sorted
	slices.Sort(tpl.Contributors)
	slices.Sort(tpl.Sponsors)

	return tpl, nil
}

func parsePresets(f *Template, presets fs.FS) error {
	for i, param := range f.Params {
		if param.Type != StringListParam {
			continue
		}

		// Load preset values from file if needed
		for j, preset := range param.Presets {
			if len(preset.Values) > 0 {
				continue
			}
			filename := fmt.Sprintf(presetFilePattern, f.Name, preset.Name)
			file, err := presets.Open(filename)
			if err != nil {
				return fmt.Errorf("preset has no value and no preset file found at %s: %w", filename, err)
			}
			lines := bufio.NewScanner(file)
			lines.Split(bufio.ScanLines)
			var values []string
			for lines.Scan() {
				values = append(values, lines.Text())
			}
			if len(values) == 0 {
				return fmt.Errorf("preset file %s is empty", filename)
			}
			f.Params[i].Presets[j].Values = values
		}
	}

	// Prepare presets slice used in rendering logic
	for _, param := range f.Params {
		if param.Type != StringListParam {
			continue
		}
		for _, preset := range param.Presets {
			header := fmt.Sprintf(presetHeaderPattern, f.Name, preset.Name)
			if len(preset.Source) > 0 {
				header += fmt.Sprintf(presetAttributionPattern, preset.Source, preset.License)
			}
			f.presets = append(f.presets, presetEntry{
				EnableKey: param.BuildPresetParamName(preset.Name),
				Name:      preset.Name,
				TargetKey: param.Name,
				Value:     preset.Values,
				Header:    header,
			})
		}
	}
	return nil
}
