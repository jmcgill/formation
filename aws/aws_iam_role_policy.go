package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/jmcgill/formation/core"
)

type AwsIamRolePolicyImporter struct {
}

// Lists all resources of this type
func (*AwsIamRolePolicyImporter) Describe(meta interface{}) ([]*core.Instance, error) {
	svc := meta.(*AWSClient).iamconn

	// List all roles
	roles := make([]*iam.Role, 0)
	err := svc.ListRolesPages(nil, func(o *iam.ListRolesOutput, lastPage bool) bool {
		for _, i := range o.Roles {
			roles = append(roles, i)
		}
		return true // continue paging
	})

	if err != nil {
		return nil, err
	}

	// For each role, list all policies
	instances := make([]*core.Instance, 0)
	for _, role := range roles {
		// Add code to list resources here
		input := &iam.ListRolePoliciesInput{
			RoleName: role.RoleName,
		}
		err := svc.ListRolePoliciesPages(input, func(o *iam.ListRolePoliciesOutput, lastPage bool) bool {
			for _, i := range o.PolicyNames {
				id := aws.StringValue(role.RoleName) + ":" + aws.StringValue(i)
				instance := &core.Instance{
					Name: core.Format(id),
					ID:   id,
				}
				instances = append(instances, instance)
			}
			return true
		})

		if err != nil {
			return nil, err
		}
	}

	return instances, nil
}

// Describes which other resources this resource can reference
func (*AwsIamRolePolicyImporter) Links() map[string]string {
	return map[string]string{
		"role": "aws_iam_role.id",
	}
}
