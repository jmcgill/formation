package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/hashicorp/terraform/terraform"
	"github.com/jmcgill/formation/core"
)

type AwsIamUserPolicyImporter struct {
}

// Lists all resources of this type
func (*AwsIamUserPolicyImporter) Describe(meta interface{}) ([]*core.Instance, error) {
	svc := meta.(*AWSClient).iamconn

	// List all users
	users := make([]*iam.User, 0)
	err := svc.ListUsersPages(nil, func(o *iam.ListUsersOutput, lastPage bool) bool {
		for _, i := range o.Users {
			users = append(users, i)
		}
		return true // continue paging
	})

	if err != nil {
		return nil, err
	}

	// For each role, list all policies
	instances := make([]*core.Instance, 0)
	for _, user := range users {
		// Add code to list resources here
		input := &iam.ListUserPoliciesInput{
			UserName: user.UserName,
		}
		err := svc.ListUserPoliciesPages(input, func(o *iam.ListUserPoliciesOutput, lastPage bool) bool {
			for _, i := range o.PolicyNames {
				id := aws.StringValue(user.UserName) + ":" + aws.StringValue(i)
				instance := &core.Instance{
					Name: core.Format(id),
					ID:   id,
					CompositeID: map[string]string{
						"user_name":   aws.StringValue(user.UserName),
						"policy_name": aws.StringValue(i),
					},
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

func (*AwsIamUserPolicyImporter) Import(in *core.Instance, meta interface{}) ([]*terraform.InstanceState, bool, error) {
	state := &terraform.InstanceState{
		ID: in.ID,
		Attributes: map[string]string{
			"user": in.CompositeID["user_name"],
		},
	}

	return []*terraform.InstanceState{
		state,
	}, false, nil
}

func (*AwsIamUserPolicyImporter) Clean(in *terraform.InstanceState, meta interface{}) *terraform.InstanceState {
	return in
}

// Describes which other resources this resource can reference
func (*AwsIamUserPolicyImporter) Links() map[string]string {
	return map[string]string{
		"user": "aws_iam_user.name",
	}
}
