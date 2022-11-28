package main

import (
	"context"
	"os"
	"path"
	"strings"

	"github.com/open-component-model/ocm/pkg/common/accessobj"
	"github.com/open-component-model/ocm/pkg/contexts/credentials/repositories/dockerconfig"
	"github.com/open-component-model/ocm/pkg/contexts/ocm"
	"github.com/open-component-model/ocm/pkg/contexts/ocm/repositories/comparch"
	ocmreg "github.com/open-component-model/ocm/pkg/contexts/ocm/repositories/ocireg"
	"github.com/open-component-model/ocm/pkg/contexts/ocm/transfer"
	"github.com/open-component-model/ocm/pkg/contexts/ocm/transfer/transferhandler"
	"github.com/open-component-model/ocm/pkg/contexts/ocm/transfer/transferhandler/standard"
	"github.com/phoban01/gitops-components/pkg/component"
	"github.com/pkg/errors"
)

type Push struct {
	Component  string `arg:"" name:"component" help:"component" type:"component"`
	Repository string `arg:"" name:"repository" help:"repository" type:"repository"`
	Copy       bool   `help:"Copy remote resources to the target repository" short:"c"`
}

func (p *Push) Help() string {
	return "This is the full help text for push"
}

func (p *Push) Validate() error {
	return nil
}

func (p *Push) Run() error {
	opts := &component.PushOpts{
		Name:       p.Component,
		Repository: p.Repository,
		Copy:       p.Copy,
	}
	return p.run(opts)
}

func (p *Push) run(opts *component.PushOpts) error {
	octx := ocm.ForContext(context.Background())
	spec := dockerconfig.NewRepositorySpec("~/.docker/config.json", true)
	if _, err := octx.CredentialsContext().RepositoryForSpec(spec); err != nil {
		return errors.Wrapf(err, "cannot access default docker config")
	}

	repo, err := octx.RepositoryForSpec(ocmreg.NewRepositorySpec(opts.Repository, nil))
	if err != nil {
		return err
	}
	defer repo.Close()

	handlerOpts := []transferhandler.TransferOption{
		standard.Recursive(opts.Copy),
		standard.ResourcesByValue(opts.Copy),
		standard.Overwrite(true),
		standard.Resolver(repo),
	}

	handler, err := standard.New(handlerOpts...)
	if err != nil {
		return err
	}

	cacheDir, err := os.UserCacheDir()
	if err != nil {
		return err
	}

	loc := path.Join(cacheDir, "ocm", strings.ReplaceAll(opts.Name, ":", "/"))

	cv, err := comparch.Open(octx, accessobj.ACC_READONLY, loc, os.ModePerm)
	if err != nil {
		return err
	}
	defer cv.Close()

	if err := transfer.TransferVersion(nil, nil, cv, repo, handler); err != nil {
		return err
	}

	return nil
}
