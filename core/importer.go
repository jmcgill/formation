package core

import "github.com/hashicorp/terraform/terraform"

type Instance struct {
	ID   string
	Name string
}

type Importer interface {
	Describe(meta interface{}) ([]*Instance, error)
	Links() map[string]string
}

// This interface should only be implemented for resources where we need to hack around bugs or limitations in
// Terraform
type PatchyImporter interface {
	Describe(meta interface{}) ([]*Instance, error)
	Links() map[string]string
	Import(*terraform.InstanceInfo, interface{}) ([]*terraform.InstanceState, bool, error)
	Clean(*terraform.InstanceState, interface{}) (*terraform.InstanceState)
}