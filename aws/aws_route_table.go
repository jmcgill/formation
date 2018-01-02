package aws

import (
	"github.com/jmcgill/formation/core"
	"github.com/aws/aws-sdk-go/aws"
)

type AwsRouteTableImporter struct {
}

// Lists all resources of this type
func (*AwsRouteTableImporter) Describe(meta interface{}) ([]*core.Instance, error) {
	svc :=  meta.(*AWSClient).ec2conn

	result, err := svc.DescribeRouteTables(nil)
	if err != nil {
	  return nil, err
	}
    existingInstances := result.RouteTables // e.g. result.Buckets

    namer := NewTagNamer()
	instances := make([]*core.Instance, len(existingInstances))
	for i, existingInstance := range existingInstances {
		instances[i] = &core.Instance{
			Name: namer.NameOrDefault(existingInstance.Tags, existingInstance.RouteTableId),
			ID:   aws.StringValue(existingInstance.RouteTableId),
		}
	}

	 return instances, nil
}

// Describes which other resources this resource can reference
func (*AwsRouteTableImporter) Links() map[string]string {
	return map[string]string{
		"vpc_id": "aws_vpc.id",
		"route.instance_id": "aws_instance.id",
		"route.gateway_id": "aws_gateway.id",
		"route.nat_gateway_id": "aws_nat_gateway.id",
		"route.egress_only_gateway_id": "aws_gateway.id",
		"route.vpc_peering_connection_id": "aws_vpc_peering_connection.id",
		"route.network_interface_id": "aws_network_interface.id",
	}
}