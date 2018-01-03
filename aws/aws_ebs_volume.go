package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/jmcgill/formation/core"
)

type AwsEbsVolumeImporter struct {
}

// Lists all resources of this type
func (*AwsEbsVolumeImporter) Describe(meta interface{}) ([]*core.Instance, error) {
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

	namer := NewTagNamer()
	instances := make([]*core.Instance, len(existingInstances))
	for i, existingInstance := range existingInstances {
		instances[i] = &core.Instance{
			Name: namer.NameOrDefault(existingInstance.Tags, existingInstance.VolumeId),
			ID:   aws.StringValue(existingInstance.VolumeId),
		}
	}

	return instances, nil
}

// Describes which other resources this resource can reference
func (*AwsEbsVolumeImporter) Links() map[string]string {
	return map[string]string{
		"snapshot_id": "aws_ebs_snapshot.id",
		"kms_key_id":  "aws_kms_key.id",
	}
}
