package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/jmcgill/formation/core"
)

type AwsCustomerGatewayImporter struct {
}

// Lists all resources of this type
func (*AwsCustomerGatewayImporter) Describe(meta interface{}) ([]*core.Instance, error) {
	svc := meta.(*AWSClient).ec2conn

	// Add code to list resources here
	result, err := svc.DescribeCustomerGateways(nil)
	if err != nil {
		return nil, err
	}
	existingInstances := result.CustomerGateways

	namer := NewTagNamer()
	instances := make([]*core.Instance, len(existingInstances))
	for i, existingInstance := range existingInstances {
		instances[i] = &core.Instance{
			Name: namer.NameOrDefault(existingInstance.Tags, existingInstance.CustomerGatewayId),
			ID:   aws.StringValue(existingInstance.CustomerGatewayId),
		}
	}

	return instances, nil
}

// Describes which other resources this resource can reference
func (*AwsCustomerGatewayImporter) Links() map[string]string {
	return map[string]string{}
}
