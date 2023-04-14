package component

import (
	"os"
	"path"
	"strings"

	"github.com/open-component-model/ocm/pkg/common"
	"github.com/open-component-model/ocm/pkg/common/accessobj"
	"github.com/open-component-model/ocm/pkg/contexts/ocm/compdesc"
	"github.com/open-component-model/ocm/pkg/contexts/ocm/repositories/comparch"
	signctx "github.com/open-component-model/ocm/pkg/contexts/ocm/signing"
	"github.com/open-component-model/ocm/pkg/signing"
	"github.com/open-component-model/ocm/pkg/signing/handlers/rsa"
)

type SignOpts struct {
	Name       string
	Signature  string // TODO: the identity should be taken from the signature itself?
	Repository string
	PrivateKey string
}

// TODO: support signing remote components using oci:// prefix
// or registry flag

// Sign signs a component using the supplied options
func (c *Context) Sign(opts *SignOpts) error {
	key, err := os.ReadFile(opts.PrivateKey)
	if err != nil {
		return err
	}

	cacheDir, err := os.UserCacheDir()
	if err != nil {
		return err
	}

	loc := path.Join(cacheDir, "ocm", strings.ReplaceAll(opts.Name, ":", "/"))

	handlerOpts := signctx.NewOptions(
		signctx.PrivateKey(opts.Signature, key),
		signctx.Sign(signing.DefaultHandlerRegistry().GetSigner(rsa.Algorithm), opts.Signature),
		signctx.Update(),
		signctx.VerifyDigests(),
	)
	handlerOpts.NormalizationAlgo = compdesc.JsonNormalisationV2
	if err := handlerOpts.Complete(nil); err != nil {
		return err
	}

	cv, err := comparch.Open(c.builder.OCMContext(), accessobj.ACC_WRITABLE, loc, os.ModePerm)
	if err != nil {
		return err
	}

	if _, err := signctx.Apply(common.NewPrinter(os.Stdout), nil, cv, handlerOpts, true); err != nil {
		return err
	}

	return nil
}
