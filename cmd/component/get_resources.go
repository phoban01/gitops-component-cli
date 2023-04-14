package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"

	"cuelang.org/go/cue"
	"cuelang.org/go/encoding/yaml"
	"github.com/phoban01/gitops-components/pkg/component"
)

type Resources struct {
	Component string `arg:"" name:"component" help:"component" type:"component"`
	Format    string `name:"format" short:"o" default:"table"`
}

func (p *Resources) Help() string {
	return ""
}

func (p *Resources) Validate() error {
	return nil
}

func (p *Resources) Run() error {
	opts := &component.DescribeOpts{
		Name:   p.Component,
		Expr:   "spec.resources",
		Format: p.Format,
	}
	return p.run(opts)
}

func (p *Resources) run(opts *component.DescribeOpts) error {
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
	case "yaml":
		result, err = yaml.Encode(*value)
	default:
		buf := &bytes.Buffer{}
		w := tabwriter.NewWriter(buf, 10, 1, 5, ' ', 0)
		fmt.Fprintln(w, "NAME\tTYPE\tVERSION")
		items, err := value.List()
		if err != nil {
			return err
		}
		for items.Next() {
			resources := items.Value()
			name := resources.LookupPath(cue.ParsePath("name"))
			rtype := resources.LookupPath(cue.ParsePath("type"))
			version := resources.LookupPath(cue.ParsePath("version"))
			fmt.Fprintf(w, "%s\t%s\t%s\n", name, rtype, version)
		}
		w.Flush()
		result = buf.Bytes()
	}

	if err != nil {
		return err
	}

	_, err = fmt.Fprint(os.Stdout, string(result))

	return err
}
