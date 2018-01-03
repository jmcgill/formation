package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/jmcgill/formation/core"
)

type AwsVpcImporter struct {
}

// Lists all resources of this type
func (*AwsVpcImporter) Describe(meta interface{}) ([]*core.Instance, error) {
	svc := meta.(*AWSClient).ec2conn

	result, err := svc.DescribeVpcs(nil)
	if err != nil {
		return nil, err
	}
	existingInstances := result.Vpcs

	namer := NewTagNamer()
	instances := make([]*core.Instance, len(existingInstances))
	for i, existingInstance := range existingInstances {
		instances[i] = &core.Instance{
			Name: namer.NameOrDefault(existingInstance.Tags, existingInstance.VpcId),
			ID:   aws.StringValue(existingInstance.VpcId),
		}
	}

	return instances, nil
}

// Describes which other resources this resource can reference
func (*AwsVpcImporter) Links() map[string]string {
	return map[string]string{}
}
