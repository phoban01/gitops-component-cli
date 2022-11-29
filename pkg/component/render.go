package component

import (
	"context"
	"os"
	"path"
	"path/filepath"
	"sync"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/load"
	"cuelang.org/go/encoding/yaml"
	"cuelang.org/go/tools/flow"
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/open-component-model/ocm/pkg/common/accessio"
	"github.com/open-component-model/ocm/pkg/contexts/credentials/repositories/dockerconfig"
	"github.com/open-component-model/ocm/pkg/contexts/oci/attrs/cacheattr"
	"github.com/open-component-model/ocm/pkg/contexts/ocm"
	"github.com/open-component-model/ocm/pkg/contexts/ocm/compdesc"
	metav1 "github.com/open-component-model/ocm/pkg/contexts/ocm/compdesc/meta/v1"
	ocmreg "github.com/open-component-model/ocm/pkg/contexts/ocm/repositories/ocireg"

	"github.com/pkg/errors"
)

type safeMap struct {
	mu     sync.Mutex
	output map[string]cue.Value
}

func newSafeMap() safeMap {
	return safeMap{
		output: make(map[string]cue.Value),
	}
}

func (o *safeMap) Add(key string, v cue.Value) {
	o.mu.Lock()
	defer o.mu.Unlock()
	o.output[key] = v
}

func (o *safeMap) Get(key string) cue.Value {
	o.mu.Lock()
	defer o.mu.Unlock()
	return o.output[key]
}

func (o *safeMap) Has(key string) bool {
	o.mu.Lock()
	defer o.mu.Unlock()
	_, found := o.output[key]
	return found
}

func (o *safeMap) Copy() map[string]cue.Value {
	return o.output
}

type RenderOpts struct {
	Filename string
}

func (c *Context) Render(opts *RenderOpts) (*cue.Value, error) {
	v, err := c.BuildInstance(&BuildOpts{Filename: opts.Filename})
	if err != nil {
		return nil, err
	}

	octx := ocm.DefaultContext()

	cacheDir, err := os.UserCacheDir()
	if err != nil {
		return nil, err
	}

	cache, err := accessio.NewStaticBlobCache(path.Join(cacheDir, "ocm", "blobs"))
	if err != nil {
		return nil, err
	}
	if err := cacheattr.Set(octx.AttributesContext(), cache); err != nil {
		return nil, err
	}

	spec := dockerconfig.NewRepositorySpec("~/.docker/config.json", true)
	if _, err := octx.CredentialsContext().RepositoryForSpec(spec); err != nil {
		return nil, errors.Wrapf(err, "cannot access default docker config")
	}

	componentCache := newSafeMap()

	cfg := &flow.Config{}

	workflow := flow.New(cfg, v, newTaskFactory(c, octx, &componentCache))

	if err := workflow.Run(context.Background()); err != nil {
		return nil, err
	}

	result := workflow.Value()

	return &result, nil
}

type RequestResourceTask struct {
	ctx    *Context
	ocmctx ocm.Context
	cache  *safeMap
}

func newTaskFactory(c *Context, octx ocm.Context, componentCache *safeMap) func(val cue.Value) (flow.Runner, error) {
	return func(val cue.Value) (flow.Runner, error) {
		request := val.LookupPath(cue.ParsePath("$method"))
		if !request.Exists() {
			return nil, nil
		}

		requestType, err := request.String()
		if err != nil {
			return nil, err
		}

		if requestType != "get-resource" {
			return nil, nil
		}

		return &RequestResourceTask{
			ctx:    c,
			ocmctx: octx,
			cache:  componentCache,
		}, nil
	}
}

func (r *RequestResourceTask) Run(t *flow.Task, pErr error) error {
	// not sure this is OK, but the value which was used for this task
	val := t.Value()

	repo, err := val.LookupPath(cue.ParsePath("repository")).String()
	if err != nil {
		return err
	}

	component, err := val.LookupPath(cue.ParsePath("component")).String()
	if err != nil {
		return err
	}

	key := path.Join(repo, component)
	if ok := r.cache.Has(key); !ok {
		ocmRepo, err := r.ocmctx.RepositoryForSpec(ocmreg.NewRepositorySpec(repo, nil))
		if err != nil {
			return err
		}
		defer ocmRepo.Close()

		res, err := r.ctx.resolveComponent(r.ocmctx, ocmRepo, component)
		if err != nil {
			return err
		}

		r.cache.Add(key, *res)
	}

	resource, err := val.LookupPath(cue.ParsePath("resource")).String()
	if err != nil {
		return err
	}

	resources, err := r.cache.Get(key).LookupPath(cue.MakePath(cue.Str("spec"), cue.Str("resources"))).List()
	if err != nil {
		return err
	}

	// iterate until we find the matching resource
	for resources.Next() {
		item := resources.Value()
		match := item.LookupPath(cue.MakePath(cue.Str("name")))
		if !match.Exists() {
			continue
		}

		matchValue, err := match.String()
		if err != nil {
			return err
		}

		if matchValue != resource {
			continue
		}

		resType, err := item.LookupPath(cue.MakePath(cue.Str("type"))).String()
		if err != nil {
			return err
		}

		if resType == "cuelang" {
			ocmRepo, err := r.ocmctx.RepositoryForSpec(ocmreg.NewRepositorySpec(repo, nil))
			if err != nil {
				return err
			}
			defer ocmRepo.Close()

			resData, err := r.ctx.resolveResourceData(r.ocmctx, ocmRepo, component, resource)
			if err != nil {
				return err
			}

			cueValue := r.ctx.context.CompileBytes(resData)

			item = item.FillPath(cue.ParsePath("data"), cueValue)
		}

		t.Fill(item)

		break
	}

	return nil
}

func (c *Context) resolveComponent(ctx ocm.Context, repo ocm.Repository, component string) (*cue.Value, error) {
	tag, err := name.NewTag(component)
	if err != nil {
		return nil, err
	}

	compvers, err := repo.LookupComponentVersion(tag.Context().Name(), tag.Identifier())
	if err != nil {
		return nil, err
	}
	defer compvers.Close()

	cd := compvers.GetDescriptor()
	data, err := compdesc.Encode(cd)
	if err != nil {
		return nil, err
	}

	cdv, err := yaml.Extract("cd", data)
	if err != nil {
		return nil, err
	}

	sources := map[string]load.Source{
		filepath.Join(workingdir, "cd.cue"): load.FromFile(cdv),
	}

	return parse(c.context, sources)
}

func (c *Context) resolveResourceData(ctx ocm.Context, repo ocm.Repository, component, resource string) ([]byte, error) {
	tag, err := name.NewTag(component)
	if err != nil {
		return nil, err
	}

	compvers, err := repo.LookupComponentVersion(tag.Context().Name(), tag.Identifier())
	if err != nil {
		return nil, err
	}
	defer compvers.Close()

	res, err := compvers.GetResource(metav1.NewIdentity(resource))
	if err != nil {
		return nil, err
	}

	acc, err := res.AccessMethod()
	if err != nil {
		return nil, err
	}

	return acc.Get()
}

func parse(ctx *cue.Context, s map[string]load.Source) (*cue.Value, error) {
	bis := load.Instances([]string{}, &load.Config{
		Dir:     workingdir,
		Package: "*",
		Overlay: s,
	})
	if len(bis) != 1 {
		return &cue.Value{}, errors.New("not vaild")
	}
	v := ctx.BuildInstance(bis[0])
	return &v, nil
}
