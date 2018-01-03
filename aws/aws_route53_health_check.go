package aws

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/davecgh/go-spew/spew"
	"github.com/jmcgill/formation/core"
)

type AwsRoute53HealthCheckImporter struct {
}

// Lists all resources of this type
func (*AwsRoute53HealthCheckImporter) Describe(meta interface{}) ([]*core.Instance, error) {
	svc := meta.(*AWSClient).r53conn
	existingInstances := make([]*route53.HealthCheck, 0)
	err := svc.ListHealthChecksPages(nil, func(o *route53.ListHealthChecksOutput, lastPage bool) bool {
		existingInstances = append(existingInstances, o.HealthChecks...)
		return true
	})

	if err != nil {
		return nil, err
	}

	fmt.Printf("Running health check importer")
	spew.Dump(existingInstances)

	instances := make([]*core.Instance, len(existingInstances))
	for i, existingInstance := range existingInstances {
		instances[i] = &core.Instance{
			Name: core.Format(aws.StringValue(existingInstance.Id)),
			ID:   aws.StringValue(existingInstance.Id),
		}
	}

	return instances, nil
}

// Describes which other resources this resource can reference
func (*AwsRoute53HealthCheckImporter) Links() map[string]string {
	return map[string]string{
		"child_healthchecks":    "aws_route53_health_check.id",
		"cloudwatch_alarm_name": "aws_cloudwatch_metric_alarm.name",
	}
}
