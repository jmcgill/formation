package aws

import (
	"github.com/jmcgill/formation/core"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/hashicorp/terraform/terraform"
	"github.com/hashicorp/terraform/config/configschema"
)

type AwsEipImporter struct {
}

// Lists all resources of this type
func (*AwsEipImporter) Describe(meta interface{}) ([]*core.Instance, error) {
	svc :=  meta.(*AWSClient).ec2conn

	// Add code to list resources here
	result, err := svc.DescribeAddresses(nil)
	if err != nil {
	  return nil, err
	}

	existingInstances := result.Addresses
	instances := make([]*core.Instance, len(existingInstances))
	for i, address := range existingInstances {
		instances[i] = &core.Instance{
			Name: core.Format(aws.StringValue(address.AllocationId)),
			ID:   aws.StringValue(address.AllocationId),
		}
	}

	 return instances, nil
}

func (*AwsEipImporter) Import(in *core.Instance, meta interface{}) ([]*terraform.InstanceState, bool, error) {
	return nil, true, nil
}

func (*AwsEipImporter) Clean(in *terraform.InstanceState, meta interface{}) *terraform.InstanceState {
	return in
}

func (*AwsEipImporter) AdjustSchema(in *configschema.Block) *configschema.Block {
	in.Attributes["vpc"].Computed = false
	in.Attributes["instance"].Computed = false
	in.Attributes["network_interface"].Computed = false
	return in
}

// Describes which other resources this resource can reference
func (*AwsEipImporter) Links() map[string]string {
	return map[string]string{
		"instance": "aws_instance.id",
		"subnet_id": "aws_subnet.id",
		"network_interface": "aws_network_interface.id",
	}
}
