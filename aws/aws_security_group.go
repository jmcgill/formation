package aws

import (
	"github.com/jmcgill/formation/core"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"strconv"
)

type AwsSecurityGroupImporter struct {
}

// Lists all resources of this type
func (*AwsSecurityGroupImporter) Describe(meta interface{}) ([]*core.Instance, error) {
	svc :=  meta.(*AWSClient).ec2conn

	names := make(map[string]int)

	// TODO(jimmy): Fold these into one
	securityGroups := make([]*ec2.SecurityGroup, 0)
	result, err := svc.DescribeSecurityGroups(nil)
	if err != nil {
	  return nil, err
	}
	securityGroups = append(securityGroups, result.SecurityGroups...)

	for result.NextToken != nil {
		input := &ec2.DescribeSecurityGroupsInput{
			NextToken: result.NextToken,
		}
		result, err := svc.DescribeSecurityGroups(input)
		if err != nil {
			return nil, err
		}
		securityGroups = append(securityGroups, result.SecurityGroups...)
	}

    existingInstances := securityGroups
	instances := make([]*core.Instance, len(existingInstances))
	for i, existingInstance := range existingInstances {
		defaultName := aws.StringValue(existingInstance.GroupId)
		name := core.Format(TagOrDefault(existingInstance.Tags, "Name", defaultName))

		if _, ok := names[name]; ok {
			names[name] += 1
			name = name + "-" + strconv.Itoa(names[name])
		} else {
			names[name] = 1
		}
		instances[i] = &core.Instance{
			Name: name,
			ID:   aws.StringValue(existingInstance.GroupId),
		}
	}

	return instances, nil
}

// Describes which other resources this resource can reference
func (*AwsSecurityGroupImporter) Links() map[string]string {
	return map[string]string{
		"vpc_id": "aws_vpc.id",
		"ingress.security_groups": "aws_security_group.id",
		"egress.security_groups": "aws_security_group.id",
	}
}