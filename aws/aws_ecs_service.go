package aws

import (
	"github.com/jmcgill/formation/core"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/hashicorp/terraform/terraform"
	"fmt"
)

// TODO
// Note: To prevent a race condition during service deletion, make sure to set depends_on to the related aws_iam_role_policy; otherwise, the policy may be destroyed too soon and the ECS service will then get stuck in the DRAINING state.

type AwsEcsServiceImporter struct {
}

// Lists all resources of this type
func (*AwsEcsServiceImporter) Describe(meta interface{}) ([]*core.Instance, error) {
	svc :=  meta.(*AWSClient).ecsconn

	// List all clusters
	clusters := make([]*string, 0)
	err := svc.ListClustersPages(nil, func(o *ecs.ListClustersOutput, lastPage bool) bool {
		for _, i := range o.ClusterArns {
			fmt.Printf("Imported cluster %s\n", aws.StringValue(i))
			clusters = append(clusters, i)
		}
		return true // continue paging
	})

	if err != nil {
		return nil, err
	}

	// List services within each cluster
	existingInstances := make([]*ecs.Service, 0)
	for _, cluster := range clusters {
		input := &ecs.ListServicesInput{
			Cluster: cluster,
		}
		err = svc.ListServicesPages(input, func(o *ecs.ListServicesOutput, lastPage bool) bool {
			for _, i := range o.ServiceArns {
				input := &ecs.DescribeServicesInput{
					Cluster: cluster,
					Services: []*string{i},
				}
				services, _ := svc.DescribeServices(input)
				for _, s := range services.Services {
					fmt.Printf("Imported ECS Service %s\n", aws.StringValue(s.ServiceArn))
					existingInstances = append(existingInstances, s)
				}
			}
			return true // continue paging
		})
	}

	if err != nil {
		return nil, err
	}

	instances := make([]*core.Instance, len(existingInstances))
	for i, existingInstance := range existingInstances {
		instances[i] = &core.Instance{
			Name: core.Format(aws.StringValue(existingInstance.ServiceName)),
			CompositeID: map[string]string{
				"cluster_arn": aws.StringValue(existingInstance.ClusterArn),
				"service_arn": aws.StringValue(existingInstance.ServiceArn),
			},
		}
	}

	return instances, nil
}

func (*AwsEcsServiceImporter) Import(in *core.Instance, meta interface{}) ([]*terraform.InstanceState, bool, error) {
	state := &terraform.InstanceState{
		ID: in.CompositeID["service_arn"],
		Attributes: map[string]string {
			"cluster": in.CompositeID["cluster_arn"],
		},
	}
	return []*terraform.InstanceState{
		state,
	}, false, nil
}

func (*AwsEcsServiceImporter) Clean(in *terraform.InstanceState, meta interface{}) (*terraform.InstanceState) {
	return in
}

// Describes which other resources this resource can reference
func (*AwsEcsServiceImporter) Links() map[string]string {
	return map[string]string{
		"task_definition": "aws_ecs_task_definition.arn",
		"cluster": "aws_ecs_cluster.arn",
		"iam_role": "aws_iam_role.arn",
		"load_balancer.elb_name": "aws_elb.name",
		"load_balancer.target_group_arn": "aws_alb_target_group.arn",
		"network_configuration.subnet": "aws_subnet.id",
		"network_configuration.security_groups": "aws_security_group.id",
	}
}