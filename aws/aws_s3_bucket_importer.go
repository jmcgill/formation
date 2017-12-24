package aws

import (
	"fmt"
	"formation/core"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type S3BucketImporter struct {
}

func (i *S3BucketImporter) Describe() core.ResourceDescription {
	sess, err := session.NewSession()

	// Create S3 service client
	// TODO(jimmy): Get this from config.meta
	svc := s3.New(sess)
	result, err := svc.ListBuckets(nil)
	if err != nil {
		fmt.Printf("Error in AWS: %s\n", err)
		panic("Unable to list buckets")
	}

	r := make([]*core.Instance, len(result.Buckets))
	for i, b := range result.Buckets {
		r[i] = &core.Instance{
			Name: strings.Replace(aws.StringValue(b.Name), "-", "_", -1),
			ID:   aws.StringValue(b.Name),
		}
	}

	var links = map[string]string{
		"logging.target_bucket":                              "aws_s3_bucket.id",
		"replication_configuration.role":                     "aws_iam_role.arn",
		"replication_configuration.rules.destination.bucket": "aws_s3_bucket.arn",
	}

	// This information is declared in the provider Schema, but not exposed through the public interface
	// Consider writing a script to automatically populate this by parsing the existing code
	var defaults = map[string]core.Default{
		"acl":                {Value: "private"},
		"force_destroy":      {Value: "false", IsBool: true},
		"version.enabled":    {Value: "false", IsBool: true},
		"version.mfa_delete": {Value: "false", IsBool: true},
	}

	return core.ResourceDescription{
		Links:     links,
		Instances: r,
		Defaults:  defaults,
	}
}
