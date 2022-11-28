package ocm

Image:              "ociImage"
Artifact:           "ociArtefact"
HelmChart:          "helmChart"
HelmRepository:     "helmRepository"
KubernetesManifest: "cue"
GitRepository:      "repository"
File:               "file"
CUE:                "cuelang"

ResourceType: Image | Artifact | KubernetesManifest | HelmChart | HelmRepository | GitRepository | File

Component: {
	apiVersion: string
	kind:       string
	metadata: {
		name:    string
		version: string
		provider: {
			name: string
		}
	}
	repositoryContexts: [...]
	spec: {
		resources: [...]
		references: [...]
		sources: [...]
	}
}

Resource: {
	name: string
	type: string
	access: {
		imageReference: string
	}
	image: access.imageReference
}

Reference: {
	type:    ResourceType
	version: string
	url:     string
}

Source: {
	type:    ResourceType
	version: string
	url:     string
}
