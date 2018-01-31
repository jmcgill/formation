package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/hashicorp/terraform/terraform"
	"github.com/jmcgill/formation/core"
)

type AwsRouteTableAssociationImporter struct {
}

// Lists all resources of this type
func (*AwsRouteTableAssociationImporter) Describe(meta interface{}) ([]*core.Instance, error) {
	svc := meta.(*AWSClient).ec2conn

	result, err := svc.DescribeRouteTables(nil)
	if err != nil {
		return nil, err
	}

	namer := NewTagNamer()
	instances := make([]*core.Instance, 0)
	for _, table := range result.RouteTables {
		name := namer.NameOrDefault(table.Tags, table.RouteTableId)
		for _, association := range table.Associations {
			// In some rare instances it is possible for an association to exist to a subnet
			// that no longer exists.
			if association.SubnetId == nil || association.RouteTableId == nil {
				continue
			}

			instances = append(instances, &core.Instance{
				Name: name,
				ID:   aws.StringValue(association.RouteTableAssociationId),
				CompositeID: map[string]string{
					"route_table_id": aws.StringValue(association.RouteTableId),
				},
			})
		}
	}

	return instances, nil
}

func (*AwsRouteTableAssociationImporter) Import(in *core.Instance, meta interface{}) ([]*terraform.InstanceState, bool, error) {
	state := &terraform.InstanceState{
		ID: in.ID,
		Attributes: map[string]string{
			"route_table_id": in.CompositeID["route_table_id"],
		},
	}

	return []*terraform.InstanceState{
		state,
	}, false, nil
}

func (*AwsRouteTableAssociationImporter) Clean(in *terraform.InstanceState, meta interface{}) *terraform.InstanceState {
	return in
}

// Describes which other resources this resource can reference
func (*AwsRouteTableAssociationImporter) Links() map[string]string {
	return map[string]string{
		"route_table_id": "aws_route_table.id",
		"subnet_id":      "aws_subnet.id",
	}
}
