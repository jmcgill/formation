package aws

import (
	"github.com/jmcgill/formation/core"
	"github.com/aws/aws-sdk-go/aws"
)

type AwsIamAccountAliasImporter struct {
}

// Lists all resources of this type
func (*AwsIamAccountAliasImporter) Describe(meta interface{}) ([]*core.Instance, error) {
	svc :=  meta.(*AWSClient).iamconn

	result, err := svc.ListAccountAliases(nil)
	if err != nil {
	  return nil, err
	}

    existingInstances := result.AccountAliases
	instances := make([]*core.Instance, len(existingInstances))
	for i, existingInstance := range existingInstances {
		instances[i] = &core.Instance{
			Name: core.Format(aws.StringValue(existingInstance)),
			ID:   aws.StringValue(existingInstance),
		}
	}

	return instances, nil
}

// Describes which other resources this resource can reference
func (*AwsIamAccountAliasImporter) Links() map[string]string {
	return map[string]string{
	}
}