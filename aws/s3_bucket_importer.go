package aws

import (
	"formation/core"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type S3BucketImporter struct {
}

func (i *S3BucketImporter) List() []*core.KnownResource {
	sess, err := session.NewSession()

	// Create S3 service client
	svc := s3.New(sess)
	result, err := svc.ListBuckets(nil)
	if err != nil {
		panic("Unable to list buckets")
	}

	r := make([]*core.KnownResource, len(result.Buckets))
	for i, b := range result.Buckets {
		r[i] = &core.KnownResource{
			Name: aws.StringValue(b.Name),
			ID:   aws.StringValue(b.Name),
		}
	}

	return r
}
