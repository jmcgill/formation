package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/hashicorp/terraform/terraform"
	"github.com/jmcgill/formation/core"
)

type AwsIamUserSshKeyImporter struct {
}

// Lists all resources of this type
func (*AwsIamUserSshKeyImporter) Describe(meta interface{}) ([]*core.Instance, error) {
	svc := meta.(*AWSClient).iamconn

	// Add code to list resources here
	existingInstances := make([]*iam.SSHPublicKeyMetadata, 0)
	err := svc.ListSSHPublicKeysPages(nil, func(o *iam.ListSSHPublicKeysOutput, lastPage bool) bool {
		for _, i := range o.SSHPublicKeys {
			existingInstances = append(existingInstances, i)
		}
		return true
	})

	if err != nil {
		return nil, err
	}

	instances := make([]*core.Instance, len(existingInstances))
	for i, existingInstance := range existingInstances {
		instances[i] = &core.Instance{
			Name: core.Format(aws.StringValue(existingInstance.UserName)),
			ID:   aws.StringValue(existingInstance.UserName),
			CompositeID: map[string]string{
				"key_id":    aws.StringValue(existingInstance.SSHPublicKeyId),
				"user_name": aws.StringValue(existingInstance.UserName),
			},
		}
	}

	return instances, nil
}

func (*AwsIamUserSshKeyImporter) Import(in *core.Instance, meta interface{}) ([]*terraform.InstanceState, bool, error) {
	state := &terraform.InstanceState{
		ID: in.CompositeID["key_id"],
		Attributes: map[string]string{
			"username": in.CompositeID["user_name"],
			"encoding": "SSH",
		},
	}

	return []*terraform.InstanceState{
		state,
	}, false, nil
}

func (*AwsIamUserSshKeyImporter) Clean(in *terraform.InstanceState, meta interface{}) *terraform.InstanceState {
	return in
}

// Describes which other resources this resource can reference
func (*AwsIamUserSshKeyImporter) Links() map[string]string {
	return map[string]string{}
}
