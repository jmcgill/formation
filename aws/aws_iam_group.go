package aws

import (
	"github.com/jmcgill/formation/core"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/iam"
)

type AwsIamGroupImporter struct {
}

// Lists all resources of this type
func (*AwsIamGroupImporter) Describe(meta interface{}) ([]*core.Instance, error) {
	svc :=  meta.(*AWSClient).iamconn

	// Add code to list resources here
	existingInstances := make([]*iam.Group, 0)
	err := svc.ListGroupsPages(nil, func(o *iam.ListGroupsOutput, lastPage bool) bool {
		for _, i := range o.Groups {
			existingInstances = append(existingInstances, i)
		}
		return true // continue paging
	})

	if err != nil {
		return nil, err
	}

	instances := make([]*core.Instance, len(existingInstances))
	for i, existingInstance := range existingInstances {
		instances[i] = &core.Instance{
			Name: core.Format(aws.StringValue(existingInstance.GroupName)),
			ID:   aws.StringValue(existingInstance.GroupName),
		}
	}

	return instances, nil
}

// Describes which other resources this resource can reference
func (*AwsIamGroupImporter) Links() map[string]string {
	return map[string]string{
	}
}