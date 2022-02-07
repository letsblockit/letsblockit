package main

import (
	"fmt"
	"io"
	"os"

	"github.com/alecthomas/kong"
	"github.com/xvello/letsblockit/src/filters"
	"gopkg.in/yaml.v2"
)

// Alias outputs to allow capturing them in tests
var (
	stdout io.Writer = os.Stdout
	stderr io.Writer = os.Stderr
)

type renderCmd struct {
	Strict bool   `help:"validate the input data before rendering the output"`
	Input  string `default:"-" help:"input file to use, defaults to stdin" arg:"" type:"existingfile"`
}

type logger struct{}

func (l *logger) Warnf(format string, args ...interface{}) {
	if _, err := fmt.Fprintf(stderr, "WARNING: "+format+"\n", args...); err != nil {
		panic(err)
	}
}

func (c *renderCmd) Run() error {
	var err error
	var input io.Reader
	if c.Input == "-" {
		input = os.Stdin
	} else {
		input, err = os.Open(c.Input)
		if err != nil {
			return fmt.Errorf("cannot open input file: %w", err)
		}
	}

	repo, err := filters.LoadFilters()
	if err != nil {
		return fmt.Errorf("cannot load filter templates: %w", err)
	}
	var list filters.List
	err = yaml.NewDecoder(input).Decode(&list)
	if err != nil {
		return fmt.Errorf("cannot decode input file: %w", err)
	}

	if c.Strict {
		err = list.Validate()
		if err != nil {
			return fmt.Errorf("invalid input data: %w", err)
		}
	}

	return list.Render(stdout, &logger{}, repo)
}

func main() {
	cmd := &renderCmd{}
	k := kong.Parse(cmd)
	k.FatalIfErrorf(cmd.Run())
}
