package aws

import (
	"github.com/jmcgill/formation/core"
	"github.com/aws/aws-sdk-go/aws"
	"strings"
	"github.com/aws/aws-sdk-go/service/elb"
)

type AwsLoadBalancerPolicyImporter struct {
}

// Lists all resources of this type
func (*AwsLoadBalancerPolicyImporter) Describe(meta interface{}) ([]*core.Instance, error) {
	svc :=  meta.(*AWSClient).elbconn

	elbs := make([]*elb.LoadBalancerDescription, 0)
	err := svc.DescribeLoadBalancersPages(nil, func(o *elb.DescribeLoadBalancersOutput, lastPage bool) bool {
		elbs = append(elbs, o.LoadBalancerDescriptions...)
		return true // continue paging
	})
	if err != nil {
		return nil, err
	}

	instances := make([]*core.Instance, 0)
	for _, instance := range elbs {
		input := &elb.DescribeLoadBalancerPoliciesInput{
			LoadBalancerName: instance.LoadBalancerName,
		}
		result, err := svc.DescribeLoadBalancerPolicies(input)
		if err != nil {
			return nil, err
		}
		for _, policy := range result.PolicyDescriptions {
			id := aws.StringValue(instance.LoadBalancerName) + ":" + aws.StringValue(policy.PolicyName)
			i := &core.Instance{
				Name: core.Format(id),
				ID:   id,
			}
			instances = append(instances, i)
		}
	}

	return instances, nil
}

func (*AwsLoadBalancerPolicyImporter) Import(in *core.Instance, meta interface{}) ([]*terraform.InstanceState, bool, error) {
	state := &terraform.InstanceState{
		ID: in.ID,
	}

	return []*terraform.InstanceState{
		state,
	}, false, nil
}

func (*AwsLoadBalancerPolicyImporter) Clean(in *terraform.InstanceState, meta interface{}) *terraform.InstanceState {
	return in
}

// Describes which other resources this resource can reference
func (*AwsLoadBalancerPolicyImporter) Links() map[string]string {
	return map[string]string{
		"load_balancer": "aws_elb.name",
	}
}
