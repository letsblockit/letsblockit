package filters

import (
	"fmt"
	"io"

	"github.com/go-playground/validator/v10"
)

const (
	listHeaderTemplate = `! Title: letsblock.it - %s
! Expires: 12 hours
! Homepage: https://letsblock.it
! License: https://github.com/letsblockit/letsblockit/blob/main/LICENSE.txt
`
	filterHeaderTemplate = `
! %s
`
)

type Instance struct {
	Filter   string                 `yaml:"filter" validate:"required"`
	Params   map[string]interface{} `yaml:"params,omitempty"`
	TestMode bool                   `yaml:"test_mode,omitempty"`
}

type List struct {
	Title     string      `yaml:"title" validate:"required"`
	Instances []*Instance `yaml:"instances" validate:"dive,required"`
	TestMode  bool        `yaml:"test_mode,omitempty"`
}

type repository interface {
	GetFilter(name string) (*Filter, error)
	Render(w io.Writer, instance *Instance) error
}

func (i *Instance) Render(out io.Writer, repo repository) error {
	_, e := fmt.Fprintf(out, filterHeaderTemplate, i.Filter)
	if e != nil {
		return e
	}
	return repo.Render(out, i)
}

func (l *List) Render(out io.Writer, logger logger, repo repository) error {
	_, err := fmt.Fprintf(out, listHeaderTemplate, l.Title)
	if err != nil {
		return err
	}

	for _, i := range l.Instances {
		if l.TestMode {
			i.TestMode = true
		}
		if err := i.Render(out, repo); err != nil {
			logger.Warnf("skipping filter %s: %s", i.Filter, err)
		}
	}
	return nil
}

func (l *List) Validate() error {
	return validator.New().Struct(l)
}
