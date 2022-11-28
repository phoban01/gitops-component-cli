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
}

sourceOutput: (source.data & {
	args: repo: chart.image.repository
}).template

release: base & ocm.ResourceRequest & {
	resource: "helmrelease"
}

releaseOutput: (release.data & {
	args: {
		values: {
			image: {
				repository: wego.image.repository
				tag:        wego.image.tag
			}
		}
	}
}).template

out: [
	sourceOutput,
	releaseOutput,
]
