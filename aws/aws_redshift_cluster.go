package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/redshift"
	"github.com/jmcgill/formation/core"
)

type AwsRedshiftClusterImporter struct {
}

// Lists all resources of this type
func (*AwsRedshiftClusterImporter) Describe(meta interface{}) ([]*core.Instance, error) {
	svc := meta.(*AWSClient).redshiftconn

	existingInstances := make([]*redshift.Cluster, 0)
	err := svc.DescribeClustersPages(nil, func(o *redshift.DescribeClustersOutput, lastPage bool) bool {
		existingInstances = append(existingInstances, o.Clusters...)
		return true
	})

	if err != nil {
		return nil, err
	}

	instances := make([]*core.Instance, len(existingInstances))
	for i, cluster := range existingInstances {
		instances[i] = &core.Instance{
			Name: core.Format(aws.StringValue(cluster.ClusterIdentifier)),
			ID:   aws.StringValue(cluster.ClusterIdentifier),
		}
	}

	return instances, nil
}

// Describes which other resources this resource can reference
func (*AwsRedshiftClusterImporter) Links() map[string]string {
	return map[string]string{
		"vpc_security_group_ids": "aws_security_group.id",
		"elastic_ip":             "aws_eip.id",
		"iam_roles":              "aws_iam_role.arn",
	}
}
