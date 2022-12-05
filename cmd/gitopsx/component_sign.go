package main

import (
	"github.com/phoban01/gitops-components/pkg/component"
)

type Sign struct {
	Component  string `arg:"" name:"component" help:"component" type:"component"`
	Signature  string `name:"signature" default:"default" help:"the name of the signature"`
	PrivateKey string `name:"private-key" help:"private key" type:"privatekey"`
}

func (p *Sign) Help() string {
	return ""
}

func (p *Sign) Validate() error {
	return nil
}

func (p *Sign) Run() error {
	opts := &component.SignOpts{
		Name:       p.Component,
		PrivateKey: p.PrivateKey,
		Signature:  p.Signature,
	}
	return p.run(opts)
}

func (p *Sign) run(opts *component.SignOpts) error {
	ctx := component.New()
	return ctx.Sign(opts)
}
