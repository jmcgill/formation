package aws

import (
	"github.com/jmcgill/formation/core"
	"github.com/aws/aws-sdk-go/aws"
)

type AwsNetworkAclImporter struct {
}

// Lists all resources of this type
func (*AwsNetworkAclImporter) Describe(meta interface{}) ([]*core.Instance, error) {
	svc :=  meta.(*AWSClient).ec2conn

	// Add code to list resources here
	result, err := svc.DescribeNetworkAcls(nil)
	if err != nil {
	  return nil, err
	}

    existingInstances := result.NetworkAcls
    names := make(map[string]int)
	instances := make([]*core.Instance, len(existingInstances))
	for i, existingInstance := range existingInstances {
		name := NameTagOrDefault(existingInstance.Tags, existingInstance.NetworkAclId, names)
		instances[i] = &core.Instance{
			Name: name,
			ID:   aws.StringValue(existingInstance.NetworkAclId),
		}
	}

	 return instances, nil
}

// Describes which other resources this resource can reference
func (*AwsNetworkAclImporter) Links() map[string]string {
	return map[string]string{
		"vpc_id": "aws_vpc.id",
		"subnet_ids": "aws_subnet.id",
	}
}