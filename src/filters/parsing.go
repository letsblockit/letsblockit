package filters

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/fs"

	"github.com/russross/blackfriday/v2"
	"gopkg.in/yaml.v3"
)

const presetFilePattern string = "filters/presets/%s/%s.txt"

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
			f.presets = append(f.presets, presetEntry{
				EnableKey: param.BuildPresetParamName(preset.Name),
				Name:      preset.Name,
				TargetKey: param.Name,
				Value:     preset.Values,
			})
		}
	}
	return nil
}
