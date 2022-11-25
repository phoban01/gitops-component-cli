import (
	"ocm.software/ocm"
)

provider: "weaveworks"

resources: [
	ocm.#File & {
		name: "deployment"
		path: "deploy.yaml"
	},

	ocm.#Image & {
		name:    "podinfo"
		image:   "ghcr.io/stefanprodan/podinfo:6.2.0"
		version: "6.2.0"
	},
]
