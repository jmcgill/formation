package aws

import (
	"github.com/jmcgill/formation/core"
	"github.com/aws/aws-sdk-go/aws"
)

type AwsVpcImporter struct {
}

// Lists all resources of this type
func (*AwsVpcImporter) Describe(meta interface{}) ([]*core.Instance, error) {
	svc :=  meta.(*AWSClient).ec2conn

	result, err := svc.DescribeVpcs(nil)
	if err != nil {
	  return nil, err
	}

    existingInstances := result.Vpcs
	instances := make([]*core.Instance, len(existingInstances))
	for i, existingInstance := range existingInstances {
		vpcId := aws.StringValue(existingInstance.VpcId)
		instances[i] = &core.Instance{
			Name: core.Format(TagOrDefault(existingInstance.Tags, "Name", vpcId)),
			ID:   aws.StringValue(existingInstance.VpcId),
		}
	}

	 return instances, nil
}

// Describes which other resources this resource can reference
func (*AwsVpcImporter) Links() map[string]string {
	return map[string]string{
	}
}