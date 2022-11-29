package main

import (
	"fmt"
	"os"

	"cuelang.org/go/cue"
	"cuelang.org/go/encoding/yaml"
	"github.com/phoban01/gitops-components/pkg/component"
	"github.com/phoban01/gitops-components/pkg/cuelibs"
)

type Render struct {
	File   string `optional:"" default:"application.cue" short:"f"`
	Expr   string `optional:"" default:"out" short:"e"`
	Format string `optional:"" short:"o"`
}

func (r *Render) Help() string {
	return "This is the full help text for render"
}

func (r *Render) Validate() error {
	return nil
}

func (r *Render) Run() error {
	ctx := component.New()
	if err := ctx.WithFS(cuelibs.Files); err != nil {
		return err
	}

	opts := component.RenderOpts{
		Filename: r.File,
	}

	res, err := ctx.Render(&opts)
	if err != nil {
		return err
	}

	out := res.LookupPath(cue.ParsePath(r.Expr))

	if r.Format == "yaml" {
		//TODO: check if this is actually a list?
		l, err := out.List()
		if err != nil {
			return err
		}
		data, err := yaml.EncodeStream(l)
		if err != nil {
			return err
		}
		fmt.Fprint(os.Stdout, string(data))
		return nil
	}

	fmt.Fprint(os.Stdout, out)

	return nil
}
