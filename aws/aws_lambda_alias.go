package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/hashicorp/terraform/terraform"
	"github.com/jmcgill/formation/core"
)

type AwsLambdaAliasImporter struct {
}

// Lists all resources of this type
func (*AwsLambdaAliasImporter) Describe(meta interface{}) ([]*core.Instance, error) {
	svc := meta.(*AWSClient).lambdaconn

	// List Functions
	functions := make([]*lambda.FunctionConfiguration, 0)
	err := svc.ListFunctionsPages(nil, func(o *lambda.ListFunctionsOutput, lastPage bool) bool {
		for _, i := range o.Functions {
			functions = append(functions, i)
		}
		return true // continue paging
	})

	if err != nil {
		return nil, err
	}

	instances := make([]*core.Instance, 0)
	for _, function := range functions {
		input := lambda.ListAliasesInput{
			FunctionName: function.FunctionName,
		}
		result, err := svc.ListAliases(&input)
		if err != nil {
			return nil, err
		}

		existingInstances := result.Aliases
		for _, existingInstance := range existingInstances {
			instances = append(instances, &core.Instance{
				Name: core.Format(aws.StringValue(existingInstance.Name)),
				ID:   aws.StringValue(existingInstance.Name),
				CompositeID: map[string]string{
					"function_name": aws.StringValue(function.FunctionName),
					"name":          aws.StringValue(existingInstance.Name),
				},
			})
		}
	}

	return instances, nil
}

func (*AwsLambdaAliasImporter) Import(in *core.Instance, meta interface{}) ([]*terraform.InstanceState, bool, error) {
	state := &terraform.InstanceState{
		ID: in.ID,
		Attributes: map[string]string{
			"function_name": in.CompositeID["function_name"],
			"name":          in.CompositeID["name"],
		},
	}

	return []*terraform.InstanceState{
		state,
	}, false, nil
}

func (*AwsLambdaAliasImporter) Clean(in *terraform.InstanceState, meta interface{}) *terraform.InstanceState {
	return in
}

// Describes which other resources this resource can reference
func (*AwsLambdaAliasImporter) Links() map[string]string {
	return map[string]string{
		"function_name": "aws_lambda_function.arn",
	}
}
