package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/hashicorp/terraform/terraform"
	"github.com/jmcgill/formation/core"
)

type AwsMainRouteTableAssociationImporter struct {
}

// Find the main route table association for each VPC
func (*AwsMainRouteTableAssociationImporter) Describe(meta interface{}) ([]*core.Instance, error) {
	svc := meta.(*AWSClient).ec2conn

	input := &ec2.DescribeRouteTablesInput{
		Filters: []*ec2.Filter{
			{
				Name: aws.String("association.main"),
				Values: []*string{
					aws.String("true"),
				},
			},
		},
	}
	result, err := svc.DescribeRouteTables(input)
	if err != nil {
		return nil, err
	}
	existingInstances := result.RouteTables // e.g. result.Buckets

	namer := NewTagNamer()
	instances := make([]*core.Instance, len(existingInstances))
	for i, existingInstance := range existingInstances {
		var associationId *string
		for _, a := range existingInstance.Associations {
			if *a.Main == true {
				associationId = a.RouteTableAssociationId
			}
		}

		// Fetch the VPC to get a better name for this association
		input := ec2.DescribeVpcsInput{
			VpcIds: []*string{
				existingInstance.VpcId,
			},
		}
		vpcs, err := svc.DescribeVpcs(&input)
		if err != nil || len(vpcs.Vpcs) != 1 {
			return nil, err
		}
		vpc := vpcs.Vpcs[0]

		vpcName := namer.NameOrDefault(vpc.Tags, vpc.VpcId)
		instances[i] = &core.Instance{
			Name: vpcName,
			ID:   aws.StringValue(associationId),
			CompositeID: map[string]string{
				"vpc_id":         aws.StringValue(existingInstance.VpcId),
				"route_table_id": aws.StringValue(existingInstance.RouteTableId),
			},
		}
	}

	return instances, nil
}

func (*AwsMainRouteTableAssociationImporter) Import(in *core.Instance, meta interface{}) ([]*terraform.InstanceState, bool, error) {
	state := &terraform.InstanceState{
		ID: in.ID,
		Attributes: map[string]string{
			"vpc_id":         in.CompositeID["vpc_id"],
			"route_table_id": in.CompositeID["route_table_id"],
		},
	}
	return []*terraform.InstanceState{
		state,
	}, false, nil
}

func (*AwsMainRouteTableAssociationImporter) Clean(in *terraform.InstanceState, meta interface{}) *terraform.InstanceState {
	return in
}

// Describes which other resources this resource can reference
func (*AwsMainRouteTableAssociationImporter) Links() map[string]string {
	return map[string]string{
		"vpc_id":         "aws_vpc.id",
		"route_table_id": "aws_route_table.id",
	}
}
