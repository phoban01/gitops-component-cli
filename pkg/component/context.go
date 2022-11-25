package component

import (
	"io/fs"
	"os"
	"path/filepath"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/cuecontext"
	"cuelang.org/go/cue/load"
	"github.com/open-component-model/ocm/pkg/common/accessobj"
	"github.com/open-component-model/ocm/pkg/contexts/ocm/compdesc"
	metav1 "github.com/open-component-model/ocm/pkg/contexts/ocm/compdesc/meta/v1"
	"github.com/open-component-model/ocm/pkg/contexts/ocm/repositories/comparch"
	"github.com/open-component-model/ocm/pkg/env"
	"github.com/open-component-model/ocm/pkg/env/builder"
)

var (
	workingdir string = "/_workspace_/"
)

type Context struct {
	context  *cue.Context
	builder  *builder.Builder
	overlays map[string]load.Source
	archive  *comparch.Object
}

type PushOpts struct {
	Name       string
	Repository string
	Vendor     bool
}

type BuildOpts struct {
	Name     string
	Version  string
	Filename string
}

func New() *Context {
	return &Context{
		context:  cuecontext.New(),
		builder:  builder.NewBuilder(env.NewEnvironment()),
		overlays: map[string]load.Source{},
	}
}

// Inspired by https://github.com/acorn-io/acorn/blob/a936079406945dc8344f9a7f07dd1ad4a90c655c/pkg/cue/instance.go
// AddFS adds filesystem to the context
func (c *Context) WithFS(files fs.FS) error {
	return fs.WalkDir(files, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}

		data, err := fs.ReadFile(files, path)
		if err != nil {
			return err
		}

		c.overlays[filepath.Join(workingdir, path)] = load.FromBytes(data)
		return nil
	})
}

func (c *Context) configure(name, version, provider string) error {
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		return err
	}

	// TODO: prune CD with resources that are no longer required, rather than recreating everytime
	dir := filepath.Join(cacheDir, "ocm", name, version)
	if _, err := os.Stat(dir); err == nil {
		if err := os.RemoveAll(dir); err != nil {
			return err
		}
	}
	if err := os.MkdirAll(dir, fs.ModePerm); err != nil {
		return err
	}

	c.archive, err = comparch.Open(c.builder.OCMContext(), accessobj.ACC_CREATE, dir, os.ModePerm)
	if err != nil {
		return err
	}

	desc := c.archive.GetDescriptor()
	desc.Metadata.ConfiguredVersion = "ocm.software/v3alpha1"
	desc.Name = name
	desc.Version = version
	desc.Provider.Name = metav1.ProviderName(provider)
	if err := compdesc.Validate(desc); err != nil {
		c.archive.Close()
		return err
	}
	return nil
}
