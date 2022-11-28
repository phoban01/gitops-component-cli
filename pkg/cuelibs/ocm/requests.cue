package ocm

import (
	"strings"
	"path"

)

ResourceRequest: {
	$method:    "get-resource"
	repository: string
	component:  string
	resource:   string
	access: {
		imageReference: string
	}
	url: access.imageReference
	image: {
		_parts:     strings.Split(access.imageReference, ":")
		_base:      _parts[0]
		tag:        _parts[1]
		repository: path.Dir(_base)
		chart:      path.Base(_base)
	}
}
