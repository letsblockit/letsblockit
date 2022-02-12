package filters

import (
	"bytes"
	"errors"
	"fmt"
	"io"

	"github.com/russross/blackfriday/v2"
	"gopkg.in/yaml.v2"
)

func parseFilter(name string, reader io.Reader) (*Filter, error) {
	filter := &Filter{
		Name: name,
	}
	return filter, parse(reader, filter)
}

func parseFilterAndTest(name string, reader io.Reader) (*FilterAndTests, error) {
	filter := &FilterAndTests{
		Filter: Filter{
			Name: name,
		},
	}
	return filter, parse(reader, filter)
}

func parse(reader io.Reader, filter filter) error {
	// Read the whole input file and parse the YAML block
	input, err := io.ReadAll(reader)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(input, filter)
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
	filter.finishParsing(string(blackfriday.Run(input[pos:])))
	return nil
}
