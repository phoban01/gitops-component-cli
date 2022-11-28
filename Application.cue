import "ocm.software/ocm"

podinfo: ocm.ResourceRequest & {
	repository: "ghcr.io/phoban01"
	component:  "github.com/phoban01/test:v1.0.3"
	resource:   "podinfo"
}

deployment: ocm.ResourceRequest & {
	repository: "ghcr.io/phoban01"
	component:  "github.com/phoban01/test:v1.0.4"
	resource:   "deployment"
}

out: (deployment.data & {
	args: {
		image:     podinfo.image
		replicas:  1
		namespace: "default"
	}
}).template
