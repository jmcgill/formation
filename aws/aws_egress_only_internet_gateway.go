package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/jmcgill/formation/core"
)

type AwsEgressOnlyInternetGatewayImporter struct {
}

// Lists all resources of this type
func (*AwsEgressOnlyInternetGatewayImporter) Describe(meta interface{}) ([]*core.Instance, error) {
	svc := meta.(*AWSClient).ec2conn

	// Add code to list resources here
	result, err := svc.DescribeEgressOnlyInternetGateways(nil)
	if err != nil {
		return nil, err
	}

	existingInstances := result.EgressOnlyInternetGateways
	instances := make([]*core.Instance, len(existingInstances))
	for i, existingInstance := range existingInstances {
		instances[i] = &core.Instance{
			Name: core.Format(aws.StringValue(existingInstance.EgressOnlyInternetGatewayId)),
			ID:   aws.StringValue(existingInstance.EgressOnlyInternetGatewayId),
		}
	}

	return instances, nil
}

// Describes which other resources this resource can reference
func (*AwsEgressOnlyInternetGatewayImporter) Links() map[string]string {
	return map[string]string{
		"vpc_id": "aws_vpc.id",
	}
}
