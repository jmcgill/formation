package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/jmcgill/formation/core"
)

type AwsNatGatewayImporter struct {
}

// Lists all resources of this type
func (*AwsNatGatewayImporter) Describe(meta interface{}) ([]*core.Instance, error) {
	svc := meta.(*AWSClient).ec2conn

	// Add code to list resources here
	existingInstances := make([]*ec2.NatGateway, 0)
	err := svc.DescribeNatGatewaysPages(nil, func(o *ec2.DescribeNatGatewaysOutput, lastPage bool) bool {
		for _, i := range o.NatGateways {
			existingInstances = append(existingInstances, i)
		}
		return true
	})

	if err != nil {
		return nil, err
	}

	namer := NewTagNamer()
	instances := make([]*core.Instance, len(existingInstances))
	for i, existingInstance := range existingInstances {
		instances[i] = &core.Instance{
			Name: namer.NameOrDefault(existingInstance.Tags, existingInstance.NatGatewayId),
			ID:   aws.StringValue(existingInstance.NatGatewayId),
		}
	}

	return instances, nil
}

// Describes which other resources this resource can reference
func (*AwsNatGatewayImporter) Links() map[string]string {
	return map[string]string{
		"allocation_id": "aws_eip.id",
		"subnet_id":     "aws_subnet.id",
	}
}
