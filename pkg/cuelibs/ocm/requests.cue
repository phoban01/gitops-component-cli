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
		_name:      _parts[0]
		tag:        _parts[1]
		repository: path.Dir(_name)
		name:       _name
		nameOnly:   path.Base(_name)
	}
	data: {
		args: {...}
		template: {...}
	}
	output: data.template
}
