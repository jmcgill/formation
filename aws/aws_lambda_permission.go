package aws

import (
	"encoding/json"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/davecgh/go-spew/spew"
	"github.com/hashicorp/terraform/terraform"
	"github.com/jmcgill/formation/core"
	aws2 "github.com/terraform-providers/terraform-provider-aws/aws"
)

type AwsLambdaPermissionImporter struct {
}

func getSidFromPolicy(in string) (string, error) {
	policyInBytes := []byte(in)
	policy := aws2.LambdaPolicy{}
	err := json.Unmarshal(policyInBytes, &policy)
	if err != nil {
		return "", err
	}

	// TODO(jimmy): Can there be multiple policy statements?
	if len(policy.Statement) > 1 {
		panic("More than one policy")
	}

	return policy.Statement[0].Sid, nil
}

func getPolicySid(svc *lambda.Lambda, functionName *string, qualifier *string) (string, bool) {
	input := lambda.GetPolicyInput{
		FunctionName: functionName,
		Qualifier:    qualifier,
	}
	p, err := svc.GetPolicy(&input)
	if err != nil {
		return "", false
	}

	sid, err := getSidFromPolicy(aws.StringValue(p.Policy))
	if err != nil {
		return "", false
	}

	return sid, true
}

// Lists all resources of this type
func (*AwsLambdaPermissionImporter) Describe(meta interface{}) ([]*core.Instance, error) {
	svc := meta.(*AWSClient).lambdaconn

	// List Functions
	functions := make([]*lambda.FunctionConfiguration, 0)
	err := svc.ListFunctionsPages(nil, func(o *lambda.ListFunctionsOutput, lastPage bool) bool {
		for _, i := range o.Functions {
			functions = append(functions, i)
		}
		return true
	})

	if err != nil {
		return nil, err
	}

	// TODO(jimmy): Should also scan all qualifiers (versions, aliases)
	instances := make([]*core.Instance, 0)
	for _, function := range functions {
		// Each function may have a policy on its root instance
		if sid, ok := getPolicySid(svc, function.FunctionName, nil); ok {
			instances = append(instances, &core.Instance{
				Name: core.Format(aws.StringValue(function.FunctionName)),
				ID:   sid,
				CompositeID: map[string]string{
					"function_name": aws.StringValue(function.FunctionName),
				},
			})
		}

		// There may also be a policy for each alias
		aliasesInput := lambda.ListAliasesInput{
			FunctionName: function.FunctionName,
		}
		result, err := svc.ListAliases(&aliasesInput)
		if err != nil {
			return nil, err
		}

		aliases := result.Aliases
		spew.Dump(result.Aliases)
		for _, alias := range aliases {
			if sid, ok := getPolicySid(svc, function.FunctionName, alias.AliasArn); ok {
				instances = append(instances, &core.Instance{
					Name: core.Format(aws.StringValue(function.FunctionName) + "-" + aws.StringValue(alias.AliasArn)),
					ID:   sid,
					CompositeID: map[string]string{
						"function_name": aws.StringValue(function.FunctionName),
						"qualifier":     aws.StringValue(alias.AliasArn),
					},
				})
			}
		}
	}

	return instances, nil
}

func (*AwsLambdaPermissionImporter) Import(in *core.Instance, meta interface{}) ([]*terraform.InstanceState, bool, error) {
	// Terraform's existing import does not import source code. We add our own importer so that we can download
	// the source too.
	state := &terraform.InstanceState{
		ID: in.ID,
		Attributes: map[string]string{
			"function_name": in.CompositeID["function_name"],
			"statement_id":  in.ID,
		},
	}

	if v, ok := in.CompositeID["qualifier"]; ok {
		state.Attributes["qualifier"] = v
	}

	return []*terraform.InstanceState{
		state,
	}, false, nil
}

func (*AwsLambdaPermissionImporter) Clean(in *terraform.InstanceState, meta interface{}) *terraform.InstanceState {
	return in
}

// Describes which other resources this resource can reference
func (*AwsLambdaPermissionImporter) Links() map[string]string {
	return map[string]string{
		"function_name": "aws_lambda_function.name",
		// Source ARN can be an s3 bucker OR a cloud watch event rule - need to support multiple
		// "source_arn": ""
	}
}
