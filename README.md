## gitopsx component CLI

** An experimental tool to manage components for GitOps **

## Getting started

The componet CLI is a tool to enable building, shipping and deploying OCM components.

(CUE)[https://cuelang.org] provides the frontend for building and rendering components.

```
# make
go install./cmd/gitopsx

# build
gitopsx component build github.com/acme/mycomponent:v1.0.0

# TODO: sign
# gitopsx component sign github.com/acme/mycomponent:v1.0.0

# push
gitopsx component push github.com/acme/mycomponent:v1.0.0 ghcr.io/$GITHUB_USER

# render (currently very slow...)
gitopsx component render -f Application.cue -oyaml
```

## Componentfile

To package a **Component** create a `Componentfile.cue`.

Here is a `Componentfile` that has two resources: a container(podinfo) and the Kubernetes configuration (app) that can be used to deploy it:
```
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
```

## Render applications

Using the Component CLI we can request resources directly via CUE. Here we request the podinfo and the deployment config. Then we render the deployment configuration using parameters from the podinfo resource:

```cuelang
import "ocm.software/ocm"

podinfo: ocm.ResourceRequest & {
	repository: "ghcr.io/phoban01"
	component:  "github.com/phoban01/test:v1.0.3"
	resource:   "podinfo"
}

deployment: ocm.ResourceRequest & {
	repository: "ghcr.io/phoban01"
	component:  "github.com/phoban01/test:v1.0.4"
	resource:   "deployment"
}

out: (deployment.data & {
	args: {
		image:     podinfo.image
		replicas:  1
		namespace: "default"
	}
}).template
```
