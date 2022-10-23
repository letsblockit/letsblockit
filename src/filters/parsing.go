package filters

import (
	"bytes"
	"errors"
	"fmt"
	"io"

	"github.com/russross/blackfriday/v2"
	"gopkg.in/yaml.v2"
)

func parseTemplate(name string, reader io.Reader) (*Template, error) {
	tpl := &Template{
		Name: name,
	}
	return tpl, parse(reader, tpl)
}

func parseTemplateAndTests(name string, reader io.Reader) (*TemplateAndTests, error) {
	tpl := &TemplateAndTests{
		Template: Template{
			Name: name,
		},
	}
	return tpl, parse(reader, tpl)
}

func parse(reader io.Reader, template template) error {
	// Read the whole input file and parse the YAML block
	input, err := io.ReadAll(reader)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(input, template)
	if err != nil {
		return fmt.Errorf("invalid metadata: %w", err)
	}

	// Find the separator and parse the markdown after it
	pos := bytes.Index(input, yamlSeparator)
	if pos < 0 {
		return errors.New("separator not found")
	}
	pos += len(yamlSeparator)
	pos += bytes.Index(input[pos:], newLine)
	template.finishParsing(string(blackfriday.Run(input[pos:])))
	return nil
}
