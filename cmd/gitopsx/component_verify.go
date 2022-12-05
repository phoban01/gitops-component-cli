package main

import (
	"github.com/phoban01/gitops-components/pkg/component"
)

type Verify struct {
	Component string `arg:"" name:"component" help:"component" type:"component"`
	Signature string `name:"signature" default:"default" help:"the name of the signature"`
	PublicKey string `name:"public-key" help:"publickey key" type:"publickey"`
}

func (p *Verify) Help() string {
	return ""
}

func (p *Verify) Validate() error {
	return nil
}

func (p *Verify) Run() error {
	opts := &component.VerifyOpts{
		Name:      p.Component,
		PublicKey: p.PublicKey,
		Signature: p.Signature,
	}
	return p.run(opts)
}

func (p *Verify) run(opts *component.VerifyOpts) error {
	ctx := component.New()
	return ctx.Verify(opts)
}
