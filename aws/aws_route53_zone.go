package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/jmcgill/formation/core"
)

type AwsRoute53ZoneImporter struct {
}

// Lists all resources of this type
func (*AwsRoute53ZoneImporter) Describe(meta interface{}) ([]*core.Instance, error) {
	svc := meta.(*AWSClient).r53conn

	// Add code to list resources here
	existingInstances := make([]*route53.HostedZone, 0)
	err := svc.ListHostedZonesPages(nil, func(o *route53.ListHostedZonesOutput, lastPage bool) bool {
		for _, i := range o.HostedZones {
			existingInstances = append(existingInstances, i)
		}
		return true // continue paging
	})

	if err != nil {
		return nil, err
	}

	instances := make([]*core.Instance, len(existingInstances))
	for i, existingInstance := range existingInstances {
		instances[i] = &core.Instance{
			Name: core.Format(aws.StringValue(existingInstance.Name)),
			ID:   aws.StringValue(existingInstance.Id),
		}
	}

	return instances, nil
}

// Describes which other resources this resource can reference
func (*AwsRoute53ZoneImporter) Links() map[string]string {
	return map[string]string{}
}
