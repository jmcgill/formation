package aws

import (
	"github.com/jmcgill/formation/core"
	"github.com/aws/aws-sdk-go/aws"
)

type AwsCustomerGatewayImporter struct {
}

// Lists all resources of this type
func (*AwsCustomerGatewayImporter) Describe(meta interface{}) ([]*core.Instance, error) {
	svc :=  meta.(*AWSClient).ec2conn

	// Add code to list resources here
	result, err := svc.DescribeCustomerGateways(nil)
	if err != nil {
	  return nil, err
	}

    names := make(map[string]int)
	existingInstances := result.CustomerGateways
	instances := make([]*core.Instance, len(existingInstances))
	for i, existingInstance := range existingInstances {
		instances[i] = &core.Instance{
			Name: NameTagOrDefault(existingInstance.Tags, existingInstance.CustomerGatewayId, names),
			ID:   aws.StringValue(existingInstance.CustomerGatewayId),
		}
	}

	 return instances, nil
}

// Describes which other resources this resource can reference
func (*AwsCustomerGatewayImporter) Links() map[string]string {
	return map[string]string{
	}
}