import "ocm.software/ocm"

podinfo: ocm.ResourceRequest & {
	repository: "ghcr.io/phoban01"
	component:  "github.com/phoban01/test:v1.0.3"
	resource:   "podinfo"
}

deployment: ocm.ResourceRequest & {
	repository: "ghcr.io/phoban01"
	component:  "github.com/phoban01/test:v1.0.3"
	resource:   "deployment"
	data: {
		args: {
			image:     podinfo.image
			namespace: "test"
			replicas:  1
		}
	}
}

out: deployment.data.template
