package main

import (
	"fmt"
	"formation/aws"
	"formation/core"
	"formation/terraform_helpers"

	"github.com/hashicorp/terraform/terraform"
	aws2 "github.com/terraform-providers/terraform-provider-aws/aws"
)

// TODO(jimmy): Decide exactly where this belongs
func RegisterImporters() map[string]core.Importer {
	return map[string]core.Importer{
		"aws_s3_bucket": &aws.S3BucketImporter{},
	}
}

type UIInput struct {
}

func (u *UIInput) Input(opts *terraform.InputOpts) (string, error) {
	fmt.Printf("Asking for input %s\n", opts.Query)
	return "us-west-2", nil
}

func main() {
	fmt.Println("Welcome to formation")

	provider := aws2.Provider()
	c := terraform.NewResourceConfig(nil)
	provider.Input(&UIInput{}, c)
	err := provider.Configure(c)
	if err != nil {
		fmt.Printf("Error in configuration: %s", err)
	}

	x := &terraform.InstanceInfo{
		// Id is a unique name to represent this instance. This is not related
		// to InstanceState.ID in any way.
		Id: "formation-test-bucket",

		// Type is the resource type of this instance
		Type: "aws_s3_bucket",
	}

	y, err := provider.ImportState(x, "formation-test-bucket")
	if err != nil {
		fmt.Printf("Error: %s", err)
	}
	fmt.Printf("Output: %s", y)

	z, err := provider.Refresh(x, y[0])
	if err != nil {
		fmt.Printf("Error: %s", err)
	}
	fmt.Printf("Output: %s", z)

	//spew.Config.DisableCapacities = true
	//spew.Config.Indent = "   "
	//spew.Config.DisablePointerMethods = true
	//spew.Config.DisablePointerAddresses = true
	//
	//spew.Dump(z.Attributes)

	for k, v := range z.Attributes {
		fmt.Printf("\"%s\": \"%s\"\n", k, v)
	}

	parser := core.InstanceStateParser{}
	printer := core.Printer{}
	var zz terraform_helpers.InstanceState
	zz = terraform_helpers.InstanceState(*z)
	r := printer.Print(parser.Parse(&zz))
	fmt.Printf("Resource: %s\n", r)

	// y should now have attribute info that I can extract and pretty print

	//// provider.ImportState(*InstanceInfo, string) ([]*InstanceState, error)
	//importers := RegisterImporters()
	//
	//// For each importer
	//for key, importer := range importers {
	//	p := provider.Resources
	//	resources := importer.List()
	//	for _, resource := range resources {
	//
	//	}
	//}
	// List the objects
	// For each object listed
	// Pull in the tfstate from terraform_helpers
}
