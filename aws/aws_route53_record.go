package aws

import (
	"github.com/jmcgill/formation/core"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/route53"
)

type AwsRoute53RecordImporter struct {
}

// Lists all resources of this type
func (*AwsRoute53RecordImporter) Describe(meta interface{}) ([]*core.Instance, error) {
	svc :=  meta.(*AWSClient).r53conn

	// Add code to list resources here
	zones := make([]*route53.HostedZone, 0)
	err := svc.ListHostedZonesPages(nil, func(o *route53.ListHostedZonesOutput, lastPage bool) bool {
		for _, i := range o.HostedZones {
			zones = append(zones, i)
		}
		return true // continue paging
	})

	if err != nil {
		return nil, err
	}

	// Add code to list resources here
	instances := make([]*core.Instance, 0)
	for _, zone := range zones {
		input := route53.ListResourceRecordSetsInput{
			HostedZoneId: zone.Id,
		}

		records := make([]*route53.ResourceRecordSet, 0)
		err := svc.ListResourceRecordSetsPages(&input, func(o *route53.ListResourceRecordSetsOutput, lastPage bool) bool {
			records = append(records, o.ResourceRecordSets...)
			return true
		})

		if err != nil {
			return nil, err
		}


		for _, record := range records {
			id := aws.StringValue(zone.Id) + "_" + aws.StringValue(record.Name) + "_" + aws.StringValue(record.Type)
			if record.SetIdentifier != nil {
				id = id + "_" + aws.StringValue(record.SetIdentifier)
			}

			instances = append(instances, &core.Instance{
				Name: id,
				ID: id,
			})

		}
	}
	return instances, nil
}

// Describes which other resources this resource can reference
func (*AwsRoute53RecordImporter) Links() map[string]string {
	return map[string]string{}
}
