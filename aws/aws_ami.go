package aws

import (
	"github.com/jmcgill/formation/core"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/hashicorp/terraform/config/configschema"
	"github.com/hashicorp/terraform/terraform"
	"github.com/aws/aws-sdk-go/service/ec2"
)

type AwsAmiImporter struct {
}

func (*AwsAmiImporter) Describe(meta interface{}) ([]*core.Instance, error) {
	svc :=  meta.(*AWSClient).ec2conn

	input := ec2.DescribeImagesInput{
		Owners: []*string{aws.String("self")},
	}

	result, err := svc.DescribeImages(&input)
	if err != nil {
	  return nil, err
	}

	existingInstances := result.Images
	instances := make([]*core.Instance, len(existingInstances))
	for i, image := range existingInstances {
		instances[i] = &core.Instance{
			Name: core.Format(aws.StringValue(image.Name)),
			ID:   aws.StringValue(image.ImageId),
		}
	}

	 return instances, nil
}

func (*AwsAmiImporter) Import(in *core.Instance, meta interface{}) ([]*terraform.InstanceState, bool, error) {
	// the source too.
	state := &terraform.InstanceState{
		ID: in.ID,
	}
	return []*terraform.InstanceState{
		state,
	}, false, nil
}

func (*AwsAmiImporter) Clean(in *terraform.InstanceState, meta interface{}) *terraform.InstanceState {
	return in
}

func (*AwsAmiImporter) AdjustSchema(in *configschema.Block) *configschema.Block {
	return in
}

// Describes which other resources this resource can reference
func (*AwsAmiImporter) Links() map[string]string {
	return map[string]string{}
}


