package aws

import (
	"github.com/jmcgill/formation/core"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
)

type AwsInstanceImporter struct {
}

func GetTag(tags []*ec2.Tag, key string) *ec2.Tag {
	for _, t := range tags {
		if aws.StringValue(t.Key) == key {
			return t
		}
	}
	return nil
}

// Lists all resources of this type
func (*AwsInstanceImporter) Describe(meta interface{}) ([]*core.Instance, error) {
	svc :=  meta.(*AWSClient).ec2conn

	existingInstances := make([]*ec2.Instance, 0)
	err := svc.DescribeInstancesPages(nil, func(o *ec2.DescribeInstancesOutput, lastPage bool) bool {
		for _, i := range o.Reservations {
			for _, j := range i.Instances {
				existingInstances = append(existingInstances, j)
			}
		}
		return true // continue paging
	})

	if err != nil {
		return nil, err
	}

	instances := make([]*core.Instance, len(existingInstances))
	for i, existingInstance := range existingInstances {
		name := ""
		tag := GetTag(existingInstance.Tags, "Name")
		if tag != nil && aws.StringValue(tag.Value) != "" {
			// Tag names are useful, but not necessarily unique
			name = aws.StringValue(tag.Value) + "-"
		}

		name += aws.StringValue(existingInstance.InstanceId)

		instances[i] = &core.Instance{
			Name: core.Format(name),
			ID:   aws.StringValue(existingInstance.InstanceId),
		}
	}

	return instances, nil
}

// Describes which other resources this resource can reference
// TODO(jimmy): The documentation for aws_instance is a bit vague, but
// I should have enough case studies to determine what should link.
// Come back and check links once all resources are imported.
func (*AwsInstanceImporter) Links() map[string]string {
	return map[string]string{
		"ami": "aws_ami.id",
		"availability_zone": "aws_availability_zone.name",
		"placement_group": "aws_placement_group.id",

		// Should this be ARN or name? The documentation suggests
		// name, but that is optional.
		"security_groups": "aws_security_group.arn",
		"vpc_security_group_ids": "aws_security_group.id",
		"subnet_id": "aws_subnet_id.id",
		"iam_instance_profile": "aws_iam_instance_profile.name",
		"ebs_block_device.device_name": "aws_ebs_volume.id",
		"ebs_block_device.snapshot_id": "aws_ebs_snapshot.id",
		"network_interface.network_interface_id": "aws_network_interface.id",
	}
}