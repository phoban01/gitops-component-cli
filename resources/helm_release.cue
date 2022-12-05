package resources

import (
	"ocm.software/ocm"
)

HelmRelease: ocm.#Cue & {
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
