package aws

import (
	"github.com/jmcgill/formation/core"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/aws"
)

type AwsIamUserImporter struct {
}

// Lists all resources of this type
func (*AwsIamUserImporter) Describe(meta interface{}) ([]*core.Instance, error) {
	svc :=  meta.(*AWSClient).iamconn

	// Add code to list resources here
	existingInstances := make([]*iam.User, 0)
	err := svc.ListUsersPages(nil, func(o *iam.ListUsersOutput, lastPage bool) bool {
		for _, i := range o.Users {
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
			Name: core.Format(aws.StringValue(existingInstance.UserName)),
			ID:   aws.StringValue(existingInstance.UserName),
		}
	}

	return instances, nil
}

// Describes which other resources this resource can reference
func (*AwsIamUserImporter) Links() map[string]string {
	return map[string]string{
	}
}