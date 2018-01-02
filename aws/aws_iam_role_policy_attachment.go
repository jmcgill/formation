package aws

import (
	"github.com/jmcgill/formation/core"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/hashicorp/terraform/terraform"
)

type AwsIamRolePolicyAttachmentImporter struct {
}

// Lists all resources of this type
func (*AwsIamRolePolicyAttachmentImporter) Describe(meta interface{}) ([]*core.Instance, error) {
	svc :=  meta.(*AWSClient).iamconn

	// List roles`
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

	instances := make([]*core.Instance, 0)
	for _, group := range roles {
		request := &iam.ListAttachedRolePoliciesInput{
			RoleName: group.RoleName,
		}
		err = svc.ListAttachedRolePoliciesPages(request, func(o *iam.ListAttachedRolePoliciesOutput, lastPage bool) bool {
			for _, policy := range o.AttachedPolicies {
				instance := &core.Instance{
					Name: core.Format(aws.StringValue(group.RoleName)) + "_" + aws.StringValue(policy.PolicyName),
					ID:   "unused",
					CompositeID: map[string]string {
						"role_name": aws.StringValue(group.RoleName),
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


func (*AwsIamRolePolicyAttachmentImporter) Import(in *core.Instance, meta interface{}) ([]*terraform.InstanceState, bool, error) {
	id := in.CompositeID["role_name"] + ":" + in.CompositeID["policy_arn"]
	state := &terraform.InstanceState{
		ID: id,
		Attributes: map[string]string {
			"role": in.CompositeID["role_name"],
			"policy_arn": in.CompositeID["policy_arn"],
		},
	}

	return []*terraform.InstanceState{
		state,
	}, false, nil
}

func (*AwsIamRolePolicyAttachmentImporter) Clean(in *terraform.InstanceState, meta interface{}) (*terraform.InstanceState) {
	return in
}

// Describes which other resources this resource can reference
func (*AwsIamRolePolicyAttachmentImporter) Links() map[string]string {
	return map[string]string{
		"role": "aws_iam_role.id",
		"policy_arn": "aws_iam_policy.arn",
	}
}