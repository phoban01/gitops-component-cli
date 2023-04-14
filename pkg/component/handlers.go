package component

import (
	"github.com/gabriel-vasile/mimetype"
	"github.com/mandelsoft/vfs/pkg/osfs"
	"github.com/open-component-model/ocm/cmds/ocm/commands/ocmcmds/common/inputs"
	"github.com/open-component-model/ocm/cmds/ocm/commands/ocmcmds/common/inputs/types/file"
	"github.com/open-component-model/ocm/cmds/ocm/commands/ocmcmds/common/inputs/types/ociimage"
	"github.com/open-component-model/ocm/pkg/common/accessio"
	"github.com/open-component-model/ocm/pkg/common/accessobj"
	"github.com/open-component-model/ocm/pkg/contexts/clictx"
	"github.com/open-component-model/ocm/pkg/contexts/datacontext/attrs/tmpcache"
	"github.com/open-component-model/ocm/pkg/contexts/ocm/accessmethods/ociartifact"
	"github.com/open-component-model/ocm/pkg/contexts/ocm/compdesc"
	metav1 "github.com/open-component-model/ocm/pkg/contexts/ocm/compdesc/meta/v1"
	"github.com/open-component-model/ocm/pkg/mime"
)

type addFileOpts struct {
	name     string
	path     string
	labels   map[string]string
	fileType string
}

// TODO: enable compression
func (c *Context) fileHandler(opts *addFileOpts) error {
	tmpcache.Set(clictx.DefaultContext(), &tmpcache.Attribute{Path: "/tmp"})

	mtype, err := mimetype.DetectFile(opts.path)
	if err != nil {
		return err
	}

	ftype := file.TYPE
	if opts.fileType != "" {
		ftype = opts.fileType
	}

	fs := osfs.New()
	acc := accessio.BlobAccessForFile(mtype.String(), opts.path, fs)

	r := &compdesc.ResourceMeta{
		ElementMeta: compdesc.ElementMeta{
			Name:          opts.name,
			ExtraIdentity: opts.labels,
		},
		Relation: metav1.LocalRelation,
		Type:     ftype,
	}

	if err := c.archive.SetResourceBlob(r, acc, "", nil); err != nil {
		return err
	}

	if err := c.archive.Update(); err != nil {
		return err
	}

	return nil
}

type addImageOpts struct {
	name    string
	labels  map[string]string
	image   string
	version string
}

// TODO:phoban01 should enable setting th etype
func (c *Context) imageHandler(opts *addImageOpts) error {
	r := &compdesc.ResourceMeta{
		ElementMeta: compdesc.ElementMeta{
			Name:          opts.name,
			ExtraIdentity: opts.labels,
			Version:       opts.version,
		},
		Relation: metav1.ExternalRelation,
		Type:     ociimage.TYPE,
	}

	// spec := ociimage.New(opts.image)
	spec := ociartifact.New(opts.image)

	if err := c.archive.SetResource(r, spec); err != nil {
		return err
	}

	if err := c.archive.Update(); err != nil {
		return err
	}

	return nil
}

type addCuelangOpts struct {
	name    string
	labels  map[string]string
	version string
	data    []byte
}

func (c *Context) cuelangHandler(opts *addCuelangOpts) error {
	r := &compdesc.ResourceMeta{
		ElementMeta: compdesc.ElementMeta{
			Name:          opts.name,
			ExtraIdentity: opts.labels,
			Version:       opts.version,
		},
		Relation: metav1.LocalRelation,
		Type:     "cuelang",
	}

	ictx := inputs.NewContext(clictx.DefaultContext(), nil, nil)
	tmpcache.Set(clictx.DefaultContext(), &tmpcache.Attribute{Path: "/tmp"})

	src := accessio.DataAccessForBytes(opts.data)
	blob := accessobj.CachedBlobAccessForDataAccess(ictx, mime.MIME_TEXT, src)

	acc, err := c.archive.AddBlob(blob, "cuelang", "", nil)
	if err != nil {
		return err
	}
	blob.Close()

	if err := c.archive.SetResource(r, acc); err != nil {
		return err
	}

	if err := c.archive.Update(); err != nil {
		return err
	}

	return nil
}
