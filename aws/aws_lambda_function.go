package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/cavaliercoder/grab"
	"github.com/hashicorp/terraform/terraform"
	"github.com/jmcgill/formation/core"
	"log"
)

type AwsLambdaFunctionImporter struct {
}

// Lists all resources of this type
func (*AwsLambdaFunctionImporter) Describe(meta interface{}) ([]*core.Instance, error) {
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

	instances := make([]*core.Instance, len(functions))
	for i, function := range functions {
		input := lambda.GetFunctionInput{
			FunctionName: function.FunctionName,
		}
		f, err := svc.GetFunction(&input)
		if err != nil {
			return nil, err
		}

		url := f.Code.Location
		filename := aws.StringValue(function.FunctionName) + ".zip"
		resp, err := grab.Get(filename, *url)
		if err != nil {
			log.Printf("[ERROR] Error downloading Lambda Function from %s\n", *url)
			return nil, err
		}

		instances[i] = &core.Instance{
			Name: core.Format(aws.StringValue(function.FunctionName)),
			ID:   aws.StringValue(function.FunctionName),
			CompositeID: map[string]string{
				"filename": resp.Filename,
			},
		}
	}

	return instances, nil
}

func (*AwsLambdaFunctionImporter) Import(in *core.Instance, meta interface{}) ([]*terraform.InstanceState, bool, error) {
	// Terraform's existing import does not import source code. We add our own importer so that we can download
	// the source too.
	state := &terraform.InstanceState{
		ID: in.ID,
		Attributes: map[string]string{
			"filename":      in.CompositeID["filename"],
			"function_name": in.ID,
		},
	}
	return []*terraform.InstanceState{
		state,
	}, false, nil
}

func (*AwsLambdaFunctionImporter) Clean(in *terraform.InstanceState, meta interface{}) *terraform.InstanceState {
	return in
}

// Describes which other resources this resource can reference
func (*AwsLambdaFunctionImporter) Links() map[string]string {
	return map[string]string{
		"role":                          "aws_iam_role.arn",
		"s3_bucket":                     "aws_s3_bucket.bucket",
		"vpc_config.security_group_ids": "aws_security_group.id",
		"vpc_config.subnet_ids":         "aws_subnet.id",
	}
}
