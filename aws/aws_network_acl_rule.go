package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/hashicorp/terraform/terraform"
	"github.com/jmcgill/formation/core"
	"strconv"
)

type AwsNetworkAclRuleImporter struct {
}

// Lists all resources of this type
func (*AwsNetworkAclRuleImporter) Describe(meta interface{}) ([]*core.Instance, error) {
	svc := meta.(*AWSClient).ec2conn

	// Add code to list resources here
	result, err := svc.DescribeNetworkAcls(nil)
	if err != nil {
		return nil, err
	}

	instances := make([]*core.Instance, 0)

	for _, acl := range result.NetworkAcls {
		for _, rule := range acl.Entries {
			number := strconv.FormatInt(*rule.RuleNumber, 10)
			instances = append(instances, &core.Instance{
				Name: core.Format(aws.StringValue(acl.NetworkAclId) + "-" + number),
				ID:   number,
				CompositeID: map[string]string{
					"rule_number":    number,
					"egress":         strconv.FormatBool(*rule.Egress),
					"network_acl_id": aws.StringValue(acl.NetworkAclId),
				},
			})
		}
	}

	return instances, nil
}

func (*AwsNetworkAclRuleImporter) Import(in *core.Instance, meta interface{}) ([]*terraform.InstanceState, bool, error) {
	state := &terraform.InstanceState{
		ID: in.ID,
		Attributes: map[string]string{
			"rule_number":    in.CompositeID["rule_number"],
			"egress":         in.CompositeID["egress"],
			"network_acl_id": in.CompositeID["network_acl_id"],
		},
	}
	return []*terraform.InstanceState{
		state,
	}, false, nil
}

func (*AwsNetworkAclRuleImporter) Clean(in *terraform.InstanceState, meta interface{}) *terraform.InstanceState {
	return in
}

// Describes which other resources this resource can reference
func (*AwsNetworkAclRuleImporter) Links() map[string]string {
	return map[string]string{
		"network_acl_id": "aws_network_acl.id",
	}
}
