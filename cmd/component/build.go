package main

import (
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/phoban01/gitops-components/pkg/component"
	"github.com/phoban01/gitops-components/pkg/cuelibs"
)

type Build struct {
	Name          string `arg:"" name:"name" help:"Name and optionally a tag" type:"name"`
	ComponentFile string `optional:"" default:"componentfile.cue" short:"f"`
}

func (b *Build) Help() string {
	return "This is the full help text for build"
}

// TODO: use proper validation
func (b *Build) Validate() error {
	if _, err := name.NewTag(b.Name); err != nil {
		return err
	}
	return nil
}

func (b *Build) Run() error {
	tag, err := name.NewTag(b.Name)
	if err != nil {
		return err
	}
	opts := &component.BuildOpts{
		Name:     tag.Repository.Name(),
		Version:  tag.TagStr(),
		Filename: b.ComponentFile,
	}
	return b.run(opts)
}

// TODO: add support for oci images
func (b *Build) run(opts *component.BuildOpts) error {
	ctx := component.New()
	if err := ctx.WithPackage("ocm.software", cuelibs.Files); err != nil {
		return err
	}
	_, err := ctx.Build(opts)
	if err != nil {
		return err
	}
	return nil
}
