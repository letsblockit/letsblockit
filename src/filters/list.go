package filters

import (
	"fmt"
	"io"

	"github.com/go-playground/validator/v10"
)

const (
	listHeaderTemplate = `! Title: letsblock.it - %s
! Expires: 1 day
! Homepage: https://letsblock.it
! License: https://github.com/xvello/letsblockit/blob/main/LICENSE.txt
`
	filterHeaderTemplate = `
! %s
`
)

type Instance struct {
	Filter string                 `yaml:"filter" validate:"required"`
	Params map[string]interface{} `yaml:"params,omitempty"`
}

type List struct {
	Title     string      `yaml:"title" validate:"required"`
	Instances []*Instance `yaml:"instances" validate:"dive,required"`
}

type repository interface {
	GetFilter(name string) (*Filter, error)
	Render(w io.Writer, name string, data map[string]interface{}) error
}

func (i *Instance) Render(out io.Writer, repo repository) error {
	_, e := fmt.Fprintf(out, filterHeaderTemplate, i.Filter)
	if e != nil {
		return e
	}
	return repo.Render(out, i.Filter, i.Params)
}

func (l *List) Render(out io.Writer, logger logger, repo repository) error {
	_, err := fmt.Fprintf(out, listHeaderTemplate, l.Title)
	if err != nil {
		return err
	}

	for _, i := range l.Instances {
		if err := i.Render(out, repo); err != nil {
			logger.Warnf("skipping filter %s: %s", i.Filter, err)
		}
	}
	return nil
}

func (l *List) Validate() error {
	return validator.New().Struct(l)
}
