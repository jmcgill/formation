package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/jmcgill/formation/core"
	"strings"
)

type AwsSnsTopicSubscriptionImporter struct {
}

// Lists all resources of this type
func (*AwsSnsTopicSubscriptionImporter) Describe(meta interface{}) ([]*core.Instance, error) {
	svc := meta.(*AWSClient).snsconn

	// Add code to list resources here
	existingInstances := make([]*sns.Subscription, 0)
	err := svc.ListSubscriptionsPages(nil, func(o *sns.ListSubscriptionsOutput, lastPage bool) bool {
		for _, i := range o.Subscriptions {
			existingInstances = append(existingInstances, i)
		}
		return true
	})

	if err != nil {
		return nil, err
	}

	instances := make([]*core.Instance, 0)
	for _, existingInstance := range existingInstances {
		// Skip SNS subscriptions that have not been validated
		if aws.StringValue(existingInstance.SubscriptionArn) == "PendingConfirmation" ||
			aws.StringValue(existingInstance.SubscriptionArn) == "Deleted" {
			continue
		}

		// Verify that the topic which we are subscribed to still exists
		input := sns.GetSubscriptionAttributesInput{
			SubscriptionArn: existingInstance.SubscriptionArn,
		}
		_, err := svc.GetSubscriptionAttributes(&input)
		if err != nil {
			continue
		}

		p := strings.Split(aws.StringValue(existingInstance.SubscriptionArn), ":")
		instances = append(instances, &core.Instance{
			Name: core.Format(p[len(p)-1]),
			ID:   aws.StringValue(existingInstance.SubscriptionArn),
		})
	}

	return instances, nil
}

// Describes which other resources this resource can reference
func (*AwsSnsTopicSubscriptionImporter) Links() map[string]string {
	return map[string]string{
		"topic_arn": "aws_sns_topic.arn",
		"endpoint":  "aws_sqs_queue.arn",
	}
}
