package component

import (
	"os"
	"path"
	"strings"

	"cuelang.org/go/cue"
	"cuelang.org/go/encoding/yaml"
	"github.com/open-component-model/ocm/pkg/common/accessobj"
	"github.com/open-component-model/ocm/pkg/contexts/ocm/compdesc"
	"github.com/open-component-model/ocm/pkg/contexts/ocm/repositories/comparch"
)

type DescribeOpts struct {
	Name   string
	Format string
	Expr   string
}

// TODO: support signing remote components using oci:// prefix
// or registry flag

// Describe verifies a component using the supplied options
func (c *Context) Describe(opts *DescribeOpts) (*cue.Value, error) {
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		return nil, err
	}

	loc := path.Join(cacheDir, "ocm", strings.ReplaceAll(opts.Name, ":", "/"))

	cv, err := comparch.Open(c.builder.OCMContext(), accessobj.ACC_READONLY, loc, os.ModePerm)
	if err != nil {
		return nil, err
	}

	cd := cv.GetDescriptor()

	data, err := compdesc.Encode(cd)
	if err != nil {
		return nil, err
	}

	file, err := yaml.Extract("cd", data)
	if err != nil {
		return nil, err
	}

	value := c.context.BuildFile(file)

	if opts.Expr != "" {
		value = value.LookupPath(cue.ParsePath(opts.Expr))
	}

	return &value, nil
}
