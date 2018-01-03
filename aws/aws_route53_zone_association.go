package aws

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/hashicorp/terraform/terraform"
	"github.com/jmcgill/formation/core"
)

type AwsRoute53ZoneAssociationImporter struct {
}

// Lists all resources of this type
func (*AwsRoute53ZoneAssociationImporter) Describe(meta interface{}) ([]*core.Instance, error) {
	svc := meta.(*AWSClient).r53conn

	// List hosted zones
	zones := make([]*route53.HostedZone, 0)
	err := svc.ListHostedZonesPages(nil, func(o *route53.ListHostedZonesOutput, lastPage bool) bool {
		zones = append(zones, o.HostedZones...)
		return true
	})

	if err != nil {
		return nil, err
	}

	instances := make([]*core.Instance, 0)
	for _, zone := range zones {
		// TODO(jimmy): Support NextToken
		input := route53.ListVPCAssociationAuthorizationsInput{
			HostedZoneId: zone.Id,
		}
		r, err := svc.ListVPCAssociationAuthorizations(&input)
		if err != nil {
			return nil, err
		}

		for _, vpc := range r.VPCs {
			name := aws.StringValue(zone.Name) + "_" + aws.StringValue(vpc.VPCId)
			id := aws.StringValue(zone.Id) + ":" + aws.StringValue(vpc.VPCId)
			instance := &core.Instance{
				Name: core.Format(name),
				ID:   id,
			}
			instances = append(instances, instance)
		}
	}

	return instances, nil
}

func (*AwsRoute53ZoneAssociationImporter) Import(in *core.Instance, meta interface{}) ([]*terraform.InstanceState, bool, error) {
	state := &terraform.InstanceState{
		ID:         in.ID,
		Attributes: map[string]string{},
	}
	return []*terraform.InstanceState{
		state,
	}, false, nil
}

func (*AwsRoute53ZoneAssociationImporter) Clean(in *terraform.InstanceState, meta interface{}) *terraform.InstanceState {
	return in
}

// Describes which other resources this resource can reference
func (*AwsRoute53ZoneAssociationImporter) Links() map[string]string {
	return map[string]string{
		"zone_id": "aws_route53_zone.zone_id",
		"vpc_id":  "aws_vpc.id",
	}
}
