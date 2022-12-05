package component

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/open-component-model/ocm/pkg/common"
	"github.com/open-component-model/ocm/pkg/common/accessobj"
	"github.com/open-component-model/ocm/pkg/contexts/ocm/compdesc"
	"github.com/open-component-model/ocm/pkg/contexts/ocm/repositories/comparch"
	"github.com/open-component-model/ocm/pkg/contexts/ocm/signing"
	signctx "github.com/open-component-model/ocm/pkg/contexts/ocm/signing"
)

type VerifyOpts struct {
	Name       string
	Signature  string
	Repository string
	PublicKey  string
}

// TODO: support signing remote components using oci:// prefix
// or registry flag

// Verify verifies a component using the supplied options
func (c *Context) Verify(opts *VerifyOpts) error {
	key, err := os.ReadFile(opts.PublicKey)
	if err != nil {
		return err
	}

	cacheDir, err := os.UserCacheDir()
	if err != nil {
		return err
	}

	loc := path.Join(cacheDir, "ocm", strings.ReplaceAll(opts.Name, ":", "/"))

	handlerOpts := signctx.NewOptions(
		signctx.PublicKey(opts.Signature, key),
		signing.VerifySignature(opts.Signature),
	)

	handlerOpts.NormalizationAlgo = compdesc.JsonNormalisationV2
	if err := handlerOpts.Complete(nil); err != nil {
		return err
	}

	cv, err := comparch.Open(c.builder.OCMContext(), accessobj.ACC_READONLY, loc, os.ModePerm)
	if err != nil {
		return err
	}

	digest, err := signctx.Apply(common.NewPrinter(os.Stdout), nil, cv, handlerOpts, true)
	if err != nil {
		return err
	}

	for _, s := range cv.GetDescriptor().Signatures {
		if s.Name == opts.Signature {
			if digest.Value != s.Digest.Value {
				return fmt.Errorf("Signature %s did not match expected value", opts.Signature)
			}
			return nil
		}
	}

	return fmt.Errorf("Signature %s not found", opts.Signature)
}
