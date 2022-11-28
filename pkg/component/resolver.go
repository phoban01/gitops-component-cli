package component

import (
	"fmt"
	"os"
	"path"
	"path/filepath"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/load"
	"cuelang.org/go/encoding/yaml"
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

type ResolveOpts struct {
	Filename string
}

func (c *Context) Resolve(opts *ResolveOpts) (*cue.Value, error) {
	cmp, err := c.BuildInstance(&BuildOpts{Filename: opts.Filename})
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

	iter, _ := cmp.Fields()

	for iter.Next() {
		v := iter.Value()
		request := v.LookupPath(cue.ParsePath("$method"))
		if request.Err() != nil {
			continue
		}
		requestType, err := request.String()
		if err != nil {
			return nil, err
		}
		switch requestType {
		case "get-resource":
			repo, err := iter.Value().LookupPath(cue.ParsePath("repository")).String()
			if err != nil {
				fmt.Println(err)
				return nil, err
			}
			component, err := iter.Value().LookupPath(cue.ParsePath("component")).String()
			if err != nil {
				return nil, err
			}

			resource, err := iter.Value().LookupPath(cue.ParsePath("resource")).String()
			if err != nil {
				return nil, err
			}

			ocmRepo, err := octx.RepositoryForSpec(ocmreg.NewRepositorySpec(repo, nil))
			if err != nil {
				return nil, err
			}
			defer ocmRepo.Close()

			resp, err := c.resolveComponent(octx, ocmRepo, component)
			if err != nil {
				return nil, err
			}

			resources, err := resp.LookupPath(cue.MakePath(cue.Str("spec"), cue.Str("resources"))).List()
			if err != nil {
				return nil, err
			}

			for resources.Next() {
				item := resources.Value()
				match := item.LookupPath(cue.MakePath(cue.Str("name")))
				if match.Err() != nil {
					continue
				}
				matchValue, err := match.String()
				if err != nil {
					return nil, err
				}
				if matchValue != resource {
					continue
				}
				resType, err := item.LookupPath(cue.MakePath(cue.Str("type"))).String()
				if err != nil {
					return nil, err
				}
				if resType == "cuelang" {
					resData, err := c.resolveResourceData(octx, ocmRepo, component, resource)
					if err != nil {
						return nil, err
					}
					cueValue := c.context.CompileBytes(resData)
					item = item.FillPath(cue.ParsePath("data"), cueValue)
				}
				*cmp = cmp.FillPath(v.Path(), item)
				break
			}
		default:
			continue
		}
	}
	return cmp, err
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
