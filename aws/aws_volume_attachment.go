package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/hashicorp/terraform/terraform"
	"github.com/jmcgill/formation/core"
)

type AwsVolumeAttachmentImporter struct {
}

// Lists all resources of this type
func (*AwsVolumeAttachmentImporter) Describe(meta interface{}) ([]*core.Instance, error) {
	svc := meta.(*AWSClient).ec2conn

	// Add code to list resources here
	existingInstances := make([]*ec2.Volume, 0)
	err := svc.DescribeVolumesPages(nil, func(o *ec2.DescribeVolumesOutput, lastPage bool) bool {
		for _, i := range o.Volumes {
			existingInstances = append(existingInstances, i)
		}
		return true // continue paging
	})

	if err != nil {
		return nil, err
	}

	instances := make([]*core.Instance, 0)
	for _, existingInstance := range existingInstances {
		for _, attachment := range existingInstance.Attachments {
			name := aws.StringValue(existingInstance.VolumeId) + "_" + aws.StringValue(attachment.Device)

			instances = append(instances, &core.Instance{
				Name: core.Format(name),
				ID:   aws.StringValue(attachment.VolumeId),
				CompositeID: map[string]string{
					"volume_id":   aws.StringValue(attachment.VolumeId),
					"device_name": aws.StringValue(attachment.Device),
					"instance_id": aws.StringValue(attachment.InstanceId),
				},
			})
		}
	}

	return instances, nil
}

func (*AwsVolumeAttachmentImporter) Import(in *core.Instance, meta interface{}) ([]*terraform.InstanceState, bool, error) {
	state := &terraform.InstanceState{
		ID: in.CompositeID["volume_id"],
		Attributes: map[string]string{
			"volume_id":   in.CompositeID["volume_id"],
			"device_name": in.CompositeID["device_name"],
			"instance_id": in.CompositeID["instance_id"],
		},
	}

	return []*terraform.InstanceState{
		state,
	}, false, nil
}

func (*AwsVolumeAttachmentImporter) Clean(in *terraform.InstanceState, meta interface{}) *terraform.InstanceState {
	return in
}

// Describes which other resources this resource can reference
func (*AwsVolumeAttachmentImporter) Links() map[string]string {
	return map[string]string{
		"instance_id": "aws_instance.id",
		"volume_id":   "aws_ebs_volume.id",
	}
}
