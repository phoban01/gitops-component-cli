package component

import (
	"errors"
	"os"
	"path/filepath"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/load"
)

func (c *Context) Build(opts *BuildOpts) (*cue.Value, error) {
	data, err := os.ReadFile(opts.Filename)
	if err != nil {
		return nil, err
	}

	c.overlays[filepath.Join(workingdir, opts.Filename)] = load.FromBytes(data)

	conf := &load.Config{
		Dir:     workingdir,
		Overlay: c.overlays,
	}

	inst := load.Instances([]string{opts.Filename}, conf)
	if len(inst) != 1 {
		return nil, errors.New("not vaild")
	}

	v := c.context.BuildInstance(inst[0])

	provider, err := v.LookupPath(cue.MakePath(cue.Str("provider"))).String()
	if err != nil {
		return nil, err
	}

	if err := c.configure(opts.Name, opts.Version, provider); err != nil {
		return nil, err
	}

	resources, err := v.LookupPath(cue.MakePath(cue.Str("resources"))).List()
	if err != nil {
		return nil, err
	}

	for {
		if !resources.Next() {
			break
		}
		res := resources.Value()

		name, err := res.LookupPath(cue.MakePath(cue.Str("name"))).String()
		if err != nil {
			return nil, err
		}

		fType, err := res.LookupPath(cue.MakePath(cue.Str("type"))).String()
		if err != nil {
			return nil, err
		}

		switch fType {
		case "file":
			path, err := res.LookupPath(cue.MakePath(cue.Str("path"))).String()
			if err != nil {
				return nil, err
			}
			o := &addFileOpts{
				name: name,
				path: path,
			}
			if err := c.fileHandler(o); err != nil {
				return nil, err
			}

		case "ociImage":
			image, err := res.LookupPath(cue.MakePath(cue.Str("image"))).String()
			if err != nil {
				return nil, err
			}
			version, err := res.LookupPath(cue.MakePath(cue.Str("version"))).String()
			if err != nil {
				return nil, err
			}
			o := &addImageOpts{
				name:    name,
				image:   image,
				version: version,
			}
			if err := c.imageHandler(o); err != nil {
				return nil, err
			}
		}
	}

	o := &addFileOpts{
		name: "componentfile",
		path: opts.Filename,
	}
	if err := c.fileHandler(o); err != nil {
		return nil, err
	}

	if err := c.archive.Close(); err != nil {
		return nil, err
	}

	return &v, nil
}
