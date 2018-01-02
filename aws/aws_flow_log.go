package aws

import (
	"github.com/jmcgill/formation/core"
	"github.com/aws/aws-sdk-go/aws"
)

type AwsFlowLogImporter struct {
}

// Lists all resources of this type
func (*AwsFlowLogImporter) Describe(meta interface{}) ([]*core.Instance, error) {
	svc :=  meta.(*AWSClient).ec2conn

	// Add code to list resources here
	result, err := svc.DescribeFlowLogs(nil)
	if err != nil {
	  return nil, err
	}

    existingInstances := result.FlowLogs
	instances := make([]*core.Instance, len(existingInstances))
	for i, existingInstance := range existingInstances {
		instances[i] = &core.Instance{
			Name: core.Format(aws.StringValue(existingInstance.FlowLogId)),
			ID:   aws.StringValue(existingInstance.FlowLogId),
		}
	}

	 return instances, nil
}

// Describes which other resources this resource can reference
func (*AwsFlowLogImporter) Links() map[string]string {
	return map[string]string{
		"log_group_name": "aws_cloudwatch_log_group.name",
		"iam_role_arn": "aws_iam_role.arn",
		"vpc_id": "aws_vpc.id",
		"subnet_id": "aws_subnet.id",
		"eni_id": "aws_network_interface.id",
	}
}