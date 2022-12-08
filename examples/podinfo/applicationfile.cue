package main

import "ocm.software/ocm"

podinfo: ocm.ResourceRequest & {
	repository: "ghcr.io/phoban01"
	component:  "github.com/phoban01/podinfo:v1.0.0"
	resource:   "podinfo"
}

deployment: ocm.ResourceRequest & {
	repository: "ghcr.io/phoban01"
	component:  "github.com/phoban01/podinfo:v1.0.0"
	resource:   "deployment"
}

out: [(deployment.data & {
	args: {
		image:     podinfo.image.name
		replicas:  1
		namespace: "default"
	}
}).template]
