package aws

import (
	"github.com/jmcgill/formation/core"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/autoscaling"
)

type AwsAutoscalingGroupImporter struct {
}

// Lists all resources of this type
func (*AwsAutoscalingGroupImporter) Describe(meta interface{}) ([]*core.Instance, error) {
	svc := meta.(*AWSClient).autoscalingconn

	// Add code to list resources here
	existingInstances := make([]*autoscaling.Group, 0)
	err := svc.DescribeAutoScalingGroupsPages(nil, func(o *autoscaling.DescribeAutoScalingGroupsOutput, lastPage bool) bool {
		for _, i := range o.AutoScalingGroups {
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
			Name: core.Format(aws.StringValue(existingInstance.AutoScalingGroupName)),
			ID:   aws.StringValue(existingInstance.AutoScalingGroupName),
		}
	}

	return instances, nil
}

// Describes which other resources this resource can reference
func (*AwsAutoscalingGroupImporter) Links() map[string]string {
	return map[string]string{
		"placement_group": "aws_placement_group.id",
		"launch_configuration": "aws_launch_configuration.name",
		"initial_lifecycle_hook.role_arn": "aws_role.arn",
	}
}
