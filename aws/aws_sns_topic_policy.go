package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/hashicorp/terraform/terraform"
	"github.com/jmcgill/formation/core"
	"strings"
)

type AwsSnsTopicPolicyImporter struct {
}

// Lists all resources of this type
func (*AwsSnsTopicPolicyImporter) Describe(meta interface{}) ([]*core.Instance, error) {
	svc := meta.(*AWSClient).snsconn

	// List topics
	topics := make([]*sns.Topic, 0)
	err := svc.ListTopicsPages(nil, func(o *sns.ListTopicsOutput, lastPage bool) bool {
		for _, i := range o.Topics {
			topics = append(topics, i)
		}
		return true
	})

	if err != nil {
		return nil, err
	}

	instances := make([]*core.Instance, 0)
	for _, topic := range topics {
		// Get policies for this topic
		input := sns.GetTopicAttributesInput{
			TopicArn: topic.TopicArn,
		}
		attributes, err := svc.GetTopicAttributes(&input)
		if err != nil {
			continue
		}

		// Only create a resource for SNS Topics that have a policy
		if _, ok := attributes.Attributes["Policy"]; ok {
			p := strings.Split(aws.StringValue(topic.TopicArn), ":")
			instances = append(instances, &core.Instance{
				Name: core.Format(p[len(p)-1]),
				ID:   aws.StringValue(topic.TopicArn),
			})
		}
	}

	return instances, nil
}

func (*AwsSnsTopicPolicyImporter) Import(in *core.Instance, meta interface{}) ([]*terraform.InstanceState, bool, error) {
	state := &terraform.InstanceState{
		ID: in.ID,
		Attributes: map[string]string{
			"arn": in.ID,
		},
	}

	return []*terraform.InstanceState{
		state,
	}, false, nil
}

func (*AwsSnsTopicPolicyImporter) Clean(in *terraform.InstanceState, meta interface{}) *terraform.InstanceState {
	if policy, ok := in.Attributes["policy"]; ok {
		in.Attributes["policy"] = SafeCleanPolicy(policy)
	}
	return in
}

// Describes which other resources this resource can reference
func (*AwsSnsTopicPolicyImporter) Links() map[string]string {
	return map[string]string{}
}
