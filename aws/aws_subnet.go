package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/jmcgill/formation/core"
)

type AwsSubnetImporter struct {
}

// Lists all resources of this type
func (*AwsSubnetImporter) Describe(meta interface{}) ([]*core.Instance, error) {
	svc := meta.(*AWSClient).ec2conn

	// Add code to list resources here
	result, err := svc.DescribeSubnets(nil)
	if err != nil {
		return nil, err
	}

	existingInstances := result.Subnets
	instances := make([]*core.Instance, len(existingInstances))
	for i, existingInstance := range existingInstances {
		instances[i] = &core.Instance{
			Name: core.Format(aws.StringValue(existingInstance.SubnetId)),
			ID:   aws.StringValue(existingInstance.SubnetId),
		}
	}

	return instances, nil
}

// Describes which other resources this resource can reference
func (*AwsSubnetImporter) Links() map[string]string {
	return map[string]string{
		"availability_zone": "aws_availability_zone.name",
		"vpc_id":            "aws_vpc.id",
	}
}
