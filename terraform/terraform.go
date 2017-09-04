package terraform

type InstanceState struct {
	// Attributes are basic information about the resource. Any keys here
	// are accessible in variable format within Terraform configurations:
	// ${resourcetype.name.attribute}.
	Attributes map[string]string
}
