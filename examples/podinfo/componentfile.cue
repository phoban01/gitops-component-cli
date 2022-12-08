import (
	"ocm.software/ocm"
)

provider: "weaveworks"

resources: {
	podinfo: ocm.#Image & {
		name:    "podinfo"
		image:   "ghcr.io/stefanprodan/podinfo:6.2.0"
		version: "6.2.0"
	}

	app: ocm.#Cue & {
		name: "deployment"
		data: {
			args: {
				image:     string
				namespace: string | *"default"
				replicas:  int | *2
			}
			template: {
				apiVersion: "apps/v1"
				kind:       "Deployment"
				metadata: {
					name:      "app"
					namespace: args.namespace
				}
				spec: {
					replicas: args.replicas
					template: spec: containers: [{
						image: args.image
					}]
				}
			}
		}
	}
}
