package aws

import (
	"github.com/jmcgill/formation/core"
)

type AwsIamAccountPasswordPolicyImporter struct {
}

// Lists all resources of this type
func (*AwsIamAccountPasswordPolicyImporter) Describe(meta interface{}) ([]*core.Instance, error) {
    // You can only have a single password policy per account, and it is always known to Terraform as
    // iam-account-password-policy
	instance :=  &core.Instance{
		Name: "policy",
		ID:   "iam-account-password-policy",
	}

	return []*core.Instance{
		instance,
	}, nil
}

// Describes which other resources this resource can reference
func (*AwsIamAccountPasswordPolicyImporter) Links() map[string]string {
	return map[string]string{
	}
}