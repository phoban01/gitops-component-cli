import (
	"ocm.software/ocm"
	defs "weave.works/resources"
)

provider: "weaveworks"

resources: {
	image: ocm.#Image & {
		name:    "image"
		image:   "ghcr.io/weaveworks/wego-app:v0.11.0"
		version: "v0.11.0"
	}

	chart: ocm.#Image & {
		name:    "chart"
		image:   "ghcr.io/weaveworks/charts/weave-gitops:4.0.8"
		version: "v0.11.0"
	}

	source: defs.HelmSource

	release: defs.HelmRelease
}
