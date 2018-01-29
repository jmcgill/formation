package aws

import (
	"github.com/jmcgill/formation/core"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/elb"
)

type AwsElbImporter struct {
}

// Lists all resources of this type
func (*AwsElbImporter) Describe(meta interface{}) ([]*core.Instance, error) {
	svc :=  meta.(*AWSClient).elbconn

	// Add code to list resources here
	existingInstances := make([]*elb.LoadBalancerDescription, 0)
	err := svc.DescribeLoadBalancersPages(nil, func(o *elb.DescribeLoadBalancersOutput, lastPage bool) bool {
		for _, i := range o.LoadBalancerDescriptions {
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
			Name: aws.StringValue(existingInstance.LoadBalancerName),
			ID:   aws.StringValue(existingInstance.LoadBalancerName),
		}
	}

	return instances, nil
}

// Describes which other resources this resource can reference
func (*AwsElbImporter) Links() map[string]string {
	return map[string]string{
		"instances": "aws_instance.id",
		"access_logs.bucket": "aws_s3_bucket.bucket",
		// "listener.ssl_certificate_id": "???"
	}
}
