package core

import (
	"github.com/hashicorp/terraform/terraform"
	"github.com/hashicorp/terraform/config/configschema"
)

type Instance struct {
	Name string

	// One of ID or CompositeID must be set
	ID          string
	CompositeID map[string]string
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
	Import(in *Instance, meta interface{}) ([]*terraform.InstanceState, bool, error)
	Clean(in *terraform.InstanceState, meta interface{}) *terraform.InstanceState
	AdjustSchema(in *configschema.Block) *configschema.Block
}
