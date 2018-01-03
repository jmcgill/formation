package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/jmcgill/formation/core"
)

type AwsIamPolicyImporter struct {
}

// Lists all resources of this type
func (*AwsIamPolicyImporter) Describe(meta interface{}) ([]*core.Instance, error) {
	svc := meta.(*AWSClient).iamconn

	// Add code to list resources here
	existingInstances := make([]*iam.Policy, 0)
	scope := "Local"
	input := &iam.ListPoliciesInput{
		Scope: &scope,
	}
	err := svc.ListPoliciesPages(input, func(o *iam.ListPoliciesOutput, lastPage bool) bool {
		for _, i := range o.Policies {
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
			Name: core.Format(aws.StringValue(existingInstance.PolicyName)),
			ID:   aws.StringValue(existingInstance.Arn),
		}
	}

	return instances, nil
}

// Describes which other resources this resource can reference
func (*AwsIamPolicyImporter) Links() map[string]string {
	return map[string]string{}
}
