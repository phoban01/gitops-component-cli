package resources

import (
	"ocm.software/ocm"
)

HelmSource: ocm.#Cue & {
	name: "source"
	data: {
		args: {
			repository: string | *"ghcr.io/weaveworks/charts"
			name:       string | *"weave-gitops"
			namespace:  string | *"default"
			interval:   string | *"10m0s"
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
				url:      "oci://\(args.repository)"
			}
		}
	}
}
