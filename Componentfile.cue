import (
	"ocm.software/ocm"
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

	source: ocm.#Cue & {
		name: "source"
		data: {
			args: {
				repo:      string | *"ghcr.io/weaveworks/charts"
				name:      string | *"weave-gitops"
				namespace: string | *"default"
				interval:  string | *"10m0s"
			}
			template: {
				apiVersion: "source.toolkit.fluxcd.io/v1beta2"
				kind:       "HelmRepository"
				metadata: {
					name:      args.name
					namespace: args.namespace
				}
				spec: {
					interval: args.interval
					type:     "oci"
					url:      "oci://\(args.repo)"
				}
			}
		}
	}

	helmrelease: ocm.#Cue & {
		name: "helmrelease"
		data: {
			args: {
				version: string | *"^4.0"
				source: {
					kind:      string | *"HelmRepository"
					name:      string | *"weave-gitops"
					namespace: string | *"default"
				}
				name:      string | *"weave-gitops"
				namespace: string | *"default"
				interval:  string | *"10m0s"
				chart:     "weave-gitops"
				values: {...}
			}
			template: {
				apiVersion: "helm.toolkit.fluxcd.io/v2beta1"
				kind:       "HelmRelease"
				metadata: {
					name:      args.name
					namespace: args.namespace
				}
				spec: {
					interval: args.interval
					chart: spec: {
						chart:     args.chart
						version:   args.version
						sourceRef: args.source
					}
					values: args.values
				}
			}
		}
	}
}
