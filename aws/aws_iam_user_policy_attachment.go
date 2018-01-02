package aws

import (
	"github.com/jmcgill/formation/core"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/hashicorp/terraform/terraform"
)

type AwsIamUserPolicyAttachmentImporter struct {
}

// Lists all resources of this type
func (*AwsIamUserPolicyAttachmentImporter) Describe(meta interface{}) ([]*core.Instance, error) {
	svc :=  meta.(*AWSClient).iamconn

	// List users
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

	instances := make([]*core.Instance, 0)
	for _, user := range users {
		request := &iam.ListAttachedUserPoliciesInput{
			UserName: user.UserName,
		}
		err = svc.ListAttachedUserPoliciesPages(request, func(o *iam.ListAttachedUserPoliciesOutput, lastPage bool) bool {
			for _, policy := range o.AttachedPolicies {
				instance := &core.Instance{
					Name: core.Format(aws.StringValue(user.UserName)) + "_" + aws.StringValue(policy.PolicyName),
					ID:   "unused",
					CompositeID: map[string]string {
						"user_name": aws.StringValue(user.UserName),
						"policy_arn": aws.StringValue(policy.PolicyArn),
					},
				}
				instances = append(instances, instance)
			}
			return true // continue paging
		})

		if err != nil {
			return nil, err
		}
	}

	return instances, nil
}


func (*AwsIamUserPolicyAttachmentImporter) Import(in *core.Instance, meta interface{}) ([]*terraform.InstanceState, bool, error) {
	id := in.CompositeID["user_name"] + ":" + in.CompositeID["policy_arn"]
	state := &terraform.InstanceState{
		ID: id,
		Attributes: map[string]string {
			"user": in.CompositeID["user_name"],
			"policy_arn": in.CompositeID["policy_arn"],
		},
	}

	return []*terraform.InstanceState{
		state,
	}, false, nil
}

func (*AwsIamUserPolicyAttachmentImporter) Clean(in *terraform.InstanceState, meta interface{}) (*terraform.InstanceState) {
	return in
}

// Describes which other resources this resource can reference
func (*AwsIamUserPolicyAttachmentImporter) Links() map[string]string {
	return map[string]string{
		"user": "aws_iam_user.name",
		"policy_arn": "aws_iam_policy.arn",
	}
}