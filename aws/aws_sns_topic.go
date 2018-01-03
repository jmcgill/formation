package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/jmcgill/formation/core"
	"strings"
)

type AwsSnsTopicImporter struct {
}

// Lists all resources of this type
func (*AwsSnsTopicImporter) Describe(meta interface{}) ([]*core.Instance, error) {
	svc := meta.(*AWSClient).snsconn

	// Add code to list resources here
	existingInstances := make([]*sns.Topic, 0)
	err := svc.ListTopicsPages(nil, func(o *sns.ListTopicsOutput, lastPage bool) bool {
		for _, i := range o.Topics {
			existingInstances = append(existingInstances, i)
		}
		return true
	})

	if err != nil {
		return nil, err
	}

	instances := make([]*core.Instance, len(existingInstances))
	for i, existingInstance := range existingInstances {
		p := strings.Split(aws.StringValue(existingInstance.TopicArn), ":")
		instances[i] = &core.Instance{
			Name: core.Format(p[len(p)-1]),
			ID:   aws.StringValue(existingInstance.TopicArn),
		}
	}

	return instances, nil
}

// Describes which other resources this resource can reference
func (*AwsSnsTopicImporter) Links() map[string]string {
	return map[string]string{}
}
