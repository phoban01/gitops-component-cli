## component CLI

** An experimental tool to manage OCM components using CUE **

Note this is very early stages work with many rough edges and performance is quite slow.

## Installation

To build the `component` CLI:

`make build`

To install in `/usr/local/bin`:

`make install`

## Getting started

The component CLI is a tool to enable building, shipping and deploying OCM components.

[CUE](https://cuelang.org) provides the frontend for building and rendering components.

## Component File

To package a **Component** create a `componentfile.cue`.

Here is a `componentfile` containing two resources:
- a container image (podinfo) 
- Kubernetes configuration (app) to deploy the container image

```cue
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

Build the component by passing the `componentfile.cue` to the build command:

`component build -f componentfile.cue acme.io/podinfo:v1.0.0`

Components can be stored in any OCI registry:

`component push acme.io/podinfo:v1.0.0 ghcr.io/octocat`

## Render applications

Using the Component CLI we can request resources using CUE and use the resource metadata in our configuration.

In the following example we request the podinfo image and deployment config.

Then we render the deployment configuration using parameters from the podinfo resource:

```cue
import "ocm.software/ocm"

podinfo: ocm.ResourceRequest & {
	repository: "ghcr.io/octocat"
	component:  "acme.io/podinfo:v1.0.0"
	resource:   "podinfo"
}

deployment: ocm.ResourceRequest & {
	repository: "ghcr.io/octocat"
	component:  "acme.io/podinfo:v1.0.0"
	resource:   "deployment"
}

out: (deployment.data & {
	args: {
		image:     podinfo.url
		replicas:  1
		namespace: "default"
	}
}).template
```

To render the output as `yaml` and apply it to the cluster, we can do the following:

`component render -f application.cue -oyaml | kubectl apply -f -`

## API

The following commands are available:

```shell
# build
component build github.com/acme/mycomponent:v1.0.0

# sign
component sign github.com/acme/mycomponent:v1.0.0 --key rsa.key

# verify
component verify github.com/acme/mycomponent:v1.0.0 --key rsa.pub

# push
component push github.com/acme/mycomponent:v1.0.0 ghcr.io/$GITHUB_USER

# describe -- show the component metadata
component describe github.com/acme/mycomponent:v1.0.0

# get resources -- print the component resources
component get resources github.com/acme/mycomponent:v1.0.0

## application file commands

# render
component render -f application.cue -oyaml

```
