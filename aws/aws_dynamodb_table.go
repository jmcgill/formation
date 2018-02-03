package aws

import (
	"github.com/jmcgill/formation/core"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type AwsDynamodbTableImporter struct {
}

// Lists all resources of this type
func (*AwsDynamodbTableImporter) Describe(meta interface{}) ([]*core.Instance, error) {
	svc :=  meta.(*AWSClient).dynamodbconn

	existingInstances := make([]*string, 0)
	err := svc.ListTablesPages(nil, func(o *dynamodb.ListTablesOutput, lastPage bool) bool {
		existingInstances = append(existingInstances, o.TableNames...)
		return true
	})

	if err != nil {
		return nil, err
	}

	instances := make([]*core.Instance, len(existingInstances))
	for i, table := range existingInstances {
		instances[i] = &core.Instance{
			Name: core.Format(aws.StringValue(table)),
			ID:   aws.StringValue(table),
		}
	}

	return instances, nil
}

// Describes which other resources this resource can reference
func (*AwsDynamodbTableImporter) Links() map[string]string {
	return map[string]string{}
}