package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"

	"github.com/letsblockit/letsblockit/src/filters"
	"gopkg.in/yaml.v3"
)

type filterAndDescription struct {
	filter filters.Template
	desc   []byte
}

func read(file io.Reader) (*filterAndDescription, error) {
	var out filterAndDescription
	input, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	pos := bytes.Index(input, []byte("---\n"))
	if pos < 0 {
		return nil, errors.New("separator not found")
	}

	err = yaml.Unmarshal(input[:pos], &out.filter)
	if err != nil {
		return nil, fmt.Errorf("invalid metadata: %w", err)
	}
	out.desc = input[pos:]

	return &out, nil
}

func (f *filterAndDescription) write(file io.Writer) error {
	enc := yaml.NewEncoder(file)
	enc.SetIndent(2)
	err := enc.Encode(&f.filter)
	if err != nil {
		return err
	}
	_, err = file.Write(f.desc)
	return err
}
