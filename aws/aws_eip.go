package aws

import (
	"github.com/jmcgill/formation/core"
	"github.com/aws/aws-sdk-go/aws"
	// "github.com/aws/aws-sdk-go/service/ec2"
	"github.com/davecgh/go-spew/spew"
)

type AwsEipImporter struct {
}

// Lists all resources of this type
func (*AwsEipImporter) Describe(meta interface{}) ([]*core.Instance, error) {
	svc :=  meta.(*AWSClient).ec2conn

	result, err := svc.DescribeAddresses(nil)
	if err != nil {
	  return nil, err
	}

	existingInstances :=  result.Addresses // e.g. result.Buckets

	namer := NewTagNamer()
	instances := make([]*core.Instance, len(existingInstances))
	for i, existingInstance := range existingInstances {
		instances[i] = &core.Instance{
			Name: namer.NameOrDefault(existingInstance.Tags, existingInstance.AllocationId),
			ID:   aws.StringValue(existingInstance.AllocationId),
		}
	}

	return instances, nil
}

// Describes which other resources this resource can reference
//func (*AwsEipImporter) ForceOutput() map[string]string {
//	return ["id", "instance"]
//}

// Describes which other resources this resource can reference
func (*AwsEipImporter) Links() map[string]string {
	return map[string]string{

	}
}
