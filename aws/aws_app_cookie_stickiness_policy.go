package aws

import (
	"github.com/jmcgill/formation/core"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/elb"
	"github.com/hashicorp/terraform/config/configschema"
	"github.com/hashicorp/terraform/terraform"
	"strconv"
)

type AwsAppCookieStickinessPolicyImporter struct {
}

func containsPolicy(s []string, b string) bool {
	for _, a := range s {
		if a == b {
			return true
		}
	}
	return false
}

// Lists all resources of this type
func (*AwsAppCookieStickinessPolicyImporter) Describe(meta interface{}) ([]*core.Instance, error) {
	svc := meta.(*AWSClient).elbconn

	// Add code to list resources here
	existingInstances := make([]*elb.LoadBalancerDescription, 0)
	err := svc.DescribeLoadBalancersPages(nil, func(o *elb.DescribeLoadBalancersOutput, lastPage bool) bool {
		existingInstances = append(existingInstances, o.LoadBalancerDescriptions...)
		return true // continue paging
	})

	if err != nil {
		return nil, err
	}

	instances := make([]*core.Instance, 0)
	for _, loadbalancer := range existingInstances {
		// Policies can be app listeners or ELB listeners. We need to build up a list of the former so that we can
		// ensure that we're importing only those policies.
		appListenerPolicies := make([]string, 0)

		for _, policy := range loadbalancer.Policies.AppCookieStickinessPolicies {
			appListenerPolicies = append(appListenerPolicies, aws.StringValue(policy.PolicyName))
		}

		for _, listener := range loadbalancer.ListenerDescriptions {
			for _, policy := range listener.PolicyNames {
				if containsPolicy(appListenerPolicies, aws.StringValue(policy)) {
					id := aws.StringValue(loadbalancer.LoadBalancerName) + ":" +
						strconv.FormatInt(*listener.Listener.LoadBalancerPort, 10) + ":"
						aws.StringValue(policy)

					instance := &core.Instance{
						Name: core.Format(id),
						ID: id,
					}
					instances = append(instances, instance)
				}
			}
		}
	}

	return instances, nil
}

// Describes which other resources this resource can reference
func (*AwsAppCookieStickinessPolicyImporter) Links() map[string]string {
	return map[string]string{}
}

func (*AwsAppCookieStickinessPolicyImporter) Import(in *core.Instance, meta interface{}) ([]*terraform.InstanceState, bool, error) {
	state := &terraform.InstanceState{
		ID: in.ID,
	}
	return []*terraform.InstanceState{
		state,
	}, false, nil
}

func (*AwsAppCookieStickinessPolicyImporter) Clean(in *terraform.InstanceState, meta interface{}) *terraform.InstanceState {
	return in
}

func (*AwsAppCookieStickinessPolicyImporter) AdjustSchema(in *configschema.Block) *configschema.Block {
	return in
}