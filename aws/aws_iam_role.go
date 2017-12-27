package aws

import (
	"github.com/jmcgill/formation/core"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/iam"
)

type AwsIamRoleImporter struct {
}

// Lists all resources of this type
func (*AwsIamRoleImporter) Describe(meta interface{}) ([]*core.Instance, error) {
	svc :=  meta.(*AWSClient).iamconn

	// Add code to list resources here
	existingInstances := make([]*iam.Role, 0)
	err := svc.ListRolesPages(nil, func(o *iam.ListRolesOutput, lastPage bool) bool {
		for _, i := range o.Roles {
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
			Name: core.Format(aws.StringValue(existingInstance.RoleName)),
			ID:   aws.StringValue(existingInstance.RoleName),
		}
	}

	return instances, nil
}

// Describes which other resources this resource can reference
func (*AwsIamRoleImporter) Links() map[string]string {
	return map[string]string{
	}
}