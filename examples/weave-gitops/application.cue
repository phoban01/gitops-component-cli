package main

import "ocm.software/ocm"

config: {
	ns: "wego-dev"
}

base: {
	repository: "ghcr.io/phoban01"
	component:  "weave.works/weave-gitops:v1.0.0"
}

wego: base & ocm.ResourceRequest & {
	resource: "image"
}

chart: base & ocm.ResourceRequest & {
	resource: "chart"
}

source: base & ocm.ResourceRequest & {
	resource: "source"
	data: args: {
		namespace:  config.ns
		repository: chart.image.repository
	}
}

release: base & ocm.ResourceRequest & {
	resource: "helmrelease"
	data: args: {
		namespace: config.ns
		source: namespace: config.ns
		values: {
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
}

out: [...{
	kind: "HelmRepository" | "HelmRelease"
	metadata: namespace: =~"^[a-z]+-dev$"
}]

out: [
	source.output,
	release.output,
]
