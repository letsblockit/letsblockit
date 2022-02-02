package main

import (
	"fmt"
	"io"
	"os"

	"github.com/alexflint/go-arg"
	"github.com/xvello/letsblockit/src/filters"
	"gopkg.in/yaml.v2"
)

// Alias outputs to allow capturing them in tests
var (
	stdout io.Writer = os.Stdout
	stderr io.Writer = os.Stderr
)

type renderCmd struct {
	Input  string `default:"-" help:"input file to use, defaults to stdin" arg:"positional"`
	Strict bool
}

type logger struct{}

func (l *logger) Warnf(format string, args ...interface{}) {
	_, err := fmt.Fprintf(stderr, "WARNING: "+format+"\n", args...)
	if err != nil {
		panic(err)
	}
}

func (c *renderCmd) Run() error {
	var err error
	var input io.ReadCloser
	if c.Input == "-" {
		input = os.Stdin
	} else {
		input, err = os.Open(c.Input)
		if err != nil {
			return err
		}
	}

	repo, err := filters.LoadFilters()
	if err != nil {
		return err
	}
	var list filters.List
	err = yaml.NewDecoder(input).Decode(&list)
	if err != nil {
		return err
	}

	if c.Strict {
		err = list.Validate()
		if err != nil {
			return err
		}
	}

	return list.Render(stdout, &logger{}, repo)
}

func main() {
	cmd := &renderCmd{}
	arg.MustParse(cmd)
	if err := cmd.Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
