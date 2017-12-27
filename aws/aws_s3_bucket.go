package aws

import (
	"github.com/jmcgill/formation/core"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
)

type AwsS3BucketImporter struct {
}

func (*AwsS3BucketImporter) Describe(meta interface{}) ([]*core.Instance, error) {
	svc := meta.(*AWSClient).s3conn

	result, err := svc.ListBuckets(nil)
	if err != nil {
		return nil, err
	}

	r := make([]*core.Instance, len(result.Buckets))
	for i, b := range result.Buckets {
		r[i] = &core.Instance{
			Name: strings.Replace(aws.StringValue(b.Name), "-", "_", -1),
			ID:   aws.StringValue(b.Name),
		}
	}

	return r, nil
}

func (*AwsS3BucketImporter) Links() map[string]string {
	return map[string]string{
		"logging.target_bucket":                              "aws_s3_bucket.id",
		"replication_configuration.role":                     "aws_iam_role.arn",
		"replication_configuration.rules.destination.bucket": "aws_s3_bucket.arn",
	}
}