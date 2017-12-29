package aws

import (
	"github.com/jmcgill/formation/core"
	"strings"
)

type AwsSqsQueueImporter struct {
}

// Lists all resources of this type
func (*AwsSqsQueueImporter) Describe(meta interface{}) ([]*core.Instance, error) {
	svc :=  meta.(*AWSClient).sqsconn

	// Add code to list resources here
	result, err := svc.ListQueues(nil)
	if err != nil {
	  return nil, err
	}

    existingInstances := result.QueueUrls
	instances := make([]*core.Instance, len(existingInstances))
	for i, existingInstance := range existingInstances {
		urlComponents := strings.Split(*existingInstance, "/")
		instances[i] = &core.Instance{
			Name: core.Format(urlComponents[len(urlComponents)-1]),
			ID:   *existingInstance,
		}
	}

	 return instances, nil
}

// Describes which other resources this resource can reference
func (*AwsSqsQueueImporter) Links() map[string]string {
	return map[string]string{
		"kms_master_key_id": "aws_kms_key.key_id",
	}
}