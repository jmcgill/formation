package aws

import (
	"github.com/jmcgill/formation/core"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/hashicorp/terraform/terraform"
)

type AwsIamGroupPolicyAttachmentImporter struct {
}

// Lists all resources of this type
func (*AwsIamGroupPolicyAttachmentImporter) Describe(meta interface{}) ([]*core.Instance, error) {
	svc :=  meta.(*AWSClient).iamconn

	// List groups
	groups := make([]*iam.Group, 0)
	err := svc.ListGroupsPages(nil, func(o *iam.ListGroupsOutput, lastPage bool) bool {
		for _, i := range o.Groups {
			groups = append(groups, i)
		}
		return true // continue paging
	})

	if err != nil {
		return nil, err
	}

	instances := make([]*core.Instance, 0)
	for _, group := range groups {
		request := &iam.ListAttachedGroupPoliciesInput{
			GroupName: group.GroupName,
		}
		err = svc.ListAttachedGroupPoliciesPages(request, func(o *iam.ListAttachedGroupPoliciesOutput, lastPage bool) bool {
			for _, policy := range o.AttachedPolicies {
				instance := &core.Instance{
					Name: core.Format(aws.StringValue(group.GroupName)) + "_" + aws.StringValue(policy.PolicyName),
					ID:   "unused",
					CompositeID: map[string]string {
						"group_name": aws.StringValue(group.GroupName),
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


func (*AwsIamGroupPolicyAttachmentImporter) Import(in *core.Instance, meta interface{}) ([]*terraform.InstanceState, bool, error) {
	id := in.CompositeID["group_name"] + ":" + in.CompositeID["policy_arn"]
	state := &terraform.InstanceState{
		ID: id,
		Attributes: map[string]string {
			"group": in.CompositeID["group_name"],
			"policy_arn": in.CompositeID["policy_arn"],
		},
	}

	return []*terraform.InstanceState{
		state,
	}, false, nil
}

func (*AwsIamGroupPolicyAttachmentImporter) Clean(in *terraform.InstanceState, meta interface{}) (*terraform.InstanceState) {
	return in
}

// Describes which other resources this resource can reference
func (*AwsIamGroupPolicyAttachmentImporter) Links() map[string]string {
	return map[string]string{
		"group": "aws_iam_group.name",
		"policy_arn": "aws_iam_policy.arn",
	}
}