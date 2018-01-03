package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/hashicorp/terraform/terraform"
	"github.com/jmcgill/formation/core"
)

type AwsIamGroupMembershipImporter struct {
}

// Lists all resources of this type
func (*AwsIamGroupMembershipImporter) Describe(meta interface{}) ([]*core.Instance, error) {
	svc := meta.(*AWSClient).iamconn

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

func (*AwsIamGroupMembershipImporter) Import(in *core.Instance, meta interface{}) ([]*terraform.InstanceState, bool, error) {
	state := &terraform.InstanceState{
		ID: in.ID,
		Attributes: map[string]string{
			"group": in.ID,
			"name":  in.Name,
		},
	}
	return []*terraform.InstanceState{
		state,
	}, false, nil
}

func (*AwsIamGroupMembershipImporter) Clean(in *terraform.InstanceState, meta interface{}) *terraform.InstanceState {
	return in
}

// Describes which other resources this resource can reference
func (*AwsIamGroupMembershipImporter) Links() map[string]string {
	return map[string]string{
		"group": "aws_iam_group.group.name",
		"users": "aws_iam_user.name",
	}
}
