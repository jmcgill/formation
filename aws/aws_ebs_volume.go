package aws

import (
	"github.com/jmcgill/formation/core"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
)

type AwsEbsVolumeImporter struct {
}

// Lists all resources of this type
func (*AwsEbsVolumeImporter) Describe(meta interface{}) ([]*core.Instance, error) {
	svc :=  meta.(*AWSClient).ec2conn

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

	instances := make([]*core.Instance, len(existingInstances))
	for i, existingInstance := range existingInstances {
		name := ""
		tag := GetTag(existingInstance.Tags, "Name")
		if tag != nil && aws.StringValue(tag.Value) != "" {
			// Tag names are useful, but not necessarily unique
			name = aws.StringValue(tag.Value) + "-"
		}

		name += aws.StringValue(existingInstance.VolumeId)

		instances[i] = &core.Instance{
			Name: core.Format(name),
			ID:   aws.StringValue(existingInstance.VolumeId),
		}
	}

	return instances, nil
}

// Describes which other resources this resource can reference
func (*AwsEbsVolumeImporter) Links() map[string]string {
	return map[string]string{
		"snapshot_id": "aws_ebs_snapshot.id",
		"kms_key_id": "aws_kms_key.id",
	}
}