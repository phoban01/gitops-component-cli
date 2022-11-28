package ocm

ResourceRequest: {
	$method:    "get-resource"
	repository: string
	component:  string
	resource:   string
	access: {
		imageReference: string
	}
	image: access.imageReference
}