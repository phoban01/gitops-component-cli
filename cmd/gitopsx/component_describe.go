package main

import (
	"encoding/json"
	"fmt"
	"os"

	"cuelang.org/go/encoding/yaml"
	"github.com/phoban01/gitops-components/pkg/component"
)

type Describe struct {
	Component string `arg:"" name:"component" help:"component" type:"component"`
	Format    string `name:"format" short:"o" default:"yaml"`
	Expr      string `name:"expression" short:"e"`
}

func (p *Describe) Help() string {
	return ""
}

func (p *Describe) Validate() error {
	return nil
}

func (p *Describe) Run() error {
	opts := &component.DescribeOpts{
		Name:   p.Component,
		Format: p.Format,
		Expr:   p.Expr,
	}
	return p.run(opts)
}

func (p *Describe) run(opts *component.DescribeOpts) error {
	var (
		result []byte
		err    error
	)

	ctx := component.New()

	value, err := ctx.Describe(opts)
	if err != nil {
		return err
	}

	switch opts.Format {
	case "json":
		result, err = json.Marshal(value)
	default:
		result, err = yaml.Encode(*value)
	}

	if err != nil {
		return err
	}

	_, err = fmt.Fprint(os.Stdout, string(result))

	return err
}
