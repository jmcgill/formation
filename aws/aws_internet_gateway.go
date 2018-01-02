package aws

import (
	"github.com/jmcgill/formation/core"
	"github.com/aws/aws-sdk-go/aws"
)

type AwsInternetGatewayImporter struct {
}

// Lists all resources of this type
func (*AwsInternetGatewayImporter) Describe(meta interface{}) ([]*core.Instance, error) {
	svc :=  meta.(*AWSClient).ec2conn

	// Add code to list resources here
	result, err := svc.DescribeInternetGateways(nil)
	if err != nil {
	  return nil, err
	}
	existingInstances := result.InternetGateways

	namer := NewTagNamer()
	instances := make([]*core.Instance, len(existingInstances))
	for i, existingInstance := range existingInstances {
		gatewayId := existingInstance.InternetGatewayId
		instances[i] = &core.Instance{
			Name: namer.NameOrDefault(existingInstance.Tags, gatewayId),
			ID:   aws.StringValue(existingInstance.InternetGatewayId),
		}
	}

	return instances, nil
}

// Describes which other resources this resource can reference
func (*AwsInternetGatewayImporter) Links() map[string]string {
	return map[string]string{
		"vpc_id": "aws_vpc.id",
	}
}