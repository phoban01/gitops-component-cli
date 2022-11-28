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
	File   string `optional:"" default:"Application.cue" short:"f"`
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

	opts := component.ResolveOpts{
		Filename: r.File,
	}

	res, err := ctx.Resolve(&opts)
	if err != nil {
		return err
	}

	//TODO: check if this is actually a list?
	out, err := res.LookupPath(cue.MakePath(cue.Str(r.Expr))).List()
	if err != nil {
		return err
	}

	if r.Format == "yaml" {
		data, err := yaml.EncodeStream(out)
		if err != nil {
			return err
		}
		fmt.Fprint(os.Stdout, string(data))
		return nil
	}

	fmt.Fprint(os.Stdout, out)

	return nil
}
