package aws

import (
	"github.com/jmcgill/formation/core"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/iam"
)

type AwsIamInstanceProfileImporter struct {
}

// Lists all resources of this type
func (*AwsIamInstanceProfileImporter) Describe(meta interface{}) ([]*core.Instance, error) {
	svc :=  meta.(*AWSClient).iamconn

	existingInstances := make([]*iam.InstanceProfile, 0)
	err := svc.ListInstanceProfilesPages(nil, func(o *iam.ListInstanceProfilesOutput, lastPage bool) bool {
		for _, i := range o.InstanceProfiles {
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
			Name: core.Format(aws.StringValue(existingInstance.InstanceProfileName)),
			ID:   aws.StringValue(existingInstance.InstanceProfileName),
		}
	}

	return instances, nil
}

// Describes which other resources this resource can reference
func (*AwsIamInstanceProfileImporter) Links() map[string]string {
	return map[string]string{
		"role": "aws_iam_role.name",
	}
}