package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/hashicorp/terraform/terraform"
	"github.com/jmcgill/formation/core"
)

type AwsIamGroupPolicyImporter struct {
}

// Lists all resources of this type
func (*AwsIamGroupPolicyImporter) Describe(meta interface{}) ([]*core.Instance, error) {
	svc := meta.(*AWSClient).iamconn

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
		request := &iam.ListGroupPoliciesInput{
			GroupName: group.GroupName,
		}
		err = svc.ListGroupPoliciesPages(request, func(o *iam.ListGroupPoliciesOutput, lastPage bool) bool {
			for _, policy := range o.PolicyNames {
				instance := &core.Instance{
					Name: core.Format(aws.StringValue(policy)),
					ID:   aws.StringValue(policy),
					CompositeID: map[string]string{
						"group_name":  aws.StringValue(group.GroupName),
						"policy_name": aws.StringValue(policy),
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

func (*AwsIamGroupPolicyImporter) Import(in *core.Instance, meta interface{}) ([]*terraform.InstanceState, bool, error) {
	id := in.CompositeID["group_name"] + ":" + in.CompositeID["policy_name"]
	state := &terraform.InstanceState{
		ID: id,
		Attributes: map[string]string{
			"group": in.CompositeID["group_name"],
			"name":  in.CompositeID["policy_name"],
		},
	}

	return []*terraform.InstanceState{
		state,
	}, false, nil
}

func (*AwsIamGroupPolicyImporter) Clean(in *terraform.InstanceState, meta interface{}) *terraform.InstanceState {
	return in
}

// Describes which other resources this resource can reference
func (*AwsIamGroupPolicyImporter) Links() map[string]string {
	return map[string]string{
		"group": "aws_iam_group.name",
	}
}
