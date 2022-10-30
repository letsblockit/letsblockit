package filters

import (
	"bytes"
	"errors"
	"fmt"
	"io"

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
