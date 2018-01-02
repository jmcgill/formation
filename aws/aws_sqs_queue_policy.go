package aws

import (
	"github.com/jmcgill/formation/core"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"strings"
)

type AwsSqsQueuePolicyImporter struct {
}

// Lists all resources of this type
func (*AwsSqsQueuePolicyImporter) Describe(meta interface{}) ([]*core.Instance, error) {
	svc :=  meta.(*AWSClient).sqsconn

	// Add code to list resources here
	result, err := svc.ListQueues(nil)
	if err != nil {
		return nil, err
	}

	existingInstances := result.QueueUrls
	instances := make([]*core.Instance, 0)
	for _, existingInstance := range existingInstances {
		urlComponents := strings.Split(*existingInstance, "/")

		// Let's ensure that a queue exists before trying to import. This slows down Describing but avoids unexpected
		// import errors.
		r, err := svc.GetQueueAttributes(&sqs.GetQueueAttributesInput{
			QueueUrl:       existingInstance,
			AttributeNames: []*string{aws.String("Policy")},
		})
		if err != nil || r.Attributes["Policy"] == nil {
			continue
		}

		instances = append(instances, &core.Instance{
			Name: core.Format(urlComponents[len(urlComponents)-1]),
			ID:   *existingInstance,
		})
	}

	return instances, nil
}

// Describes which other resources this resource can reference
func (*AwsSqsQueuePolicyImporter) Links() map[string]string {
	return map[string]string{
		"queue_url": "aws_sqs_queue.id",
	}
}