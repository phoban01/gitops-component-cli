package ocm

Image:              "ociImage"
Artifact:           "ociArtefact"
HelmChart:          "helmChart"
HelmRepository:     "helmRepository"
KubernetesManifest: "cue"
GitRepository:      "repository"
File:               "file"

ResourceType: Image | Artifact | KubernetesManifest | HelmChart | HelmRepository | GitRepository | File

#Component: {
	version:   string
	name:      string
	namespace: string
	provider:  string
	resources: [string]:  #Resource
	references: [string]: #Reference
	sources: [string]:    #Source
}

#Resource: {
	type:    ResourceType
	version: string
	url:     string
}

#Reference: {
	type:    ResourceType
	version: string
	url:     string
}

#Source: {
	type:    ResourceType
	version: string
	url:     string
}
