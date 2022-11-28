import "ocm.software/ocm"

base: {
	repository: "ghcr.io/phoban01"
	component:  "github.com/phoban01/weave-gitops:v1.0.0"
}

wego: base & ocm.ResourceRequest & {
	resource: "image"
}

chart: base & ocm.ResourceRequest & {
	resource: "chart"
}

source: base & ocm.ResourceRequest & {
	resource: "source"
	data: args: repo: chart.image.repository
}

release: base & ocm.ResourceRequest & {
	resource: "helmrelease"
	data: args: values: {
		image: {
			repository: wego.image.name
			tag:        wego.image.tag
		}
		adminUser: {
			create:       true
			username:     "admin"
			passwordHash: "$2y$10$zTRdq9bLcEmGF27exGcKZ.LnSNIOpwV.n5H7tLP4/oyuSRGjTk7Ai"
		}
	}
}

out: [
	source.data.template,
	release.data.template,
]
