package ocm

Input: {
	type: string
	name: string
	labels?: {[string]: string}
	version?: string
}

#File: Input & {
	type: File
	path: string
}

#Image: Input & {
	type:  Image
	image: string
}

#Cue: Input & {
	type: CUE
	data: {
		args: {...}
		template: {...}
	}
}
