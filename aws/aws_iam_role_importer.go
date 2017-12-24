package aws

import (
	"fmt"
	"formation/core"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
)

type IAmRoleImporter struct {
}

func (this *IAmRoleImporter) Describe() core.ResourceDescription {
	sess, err := session.NewSession()

	svc := iam.New(sess)
	result, err := svc.ListRoles(nil)
	if err != nil {
		fmt.Printf("Error in AWS: %s\n", err)
		panic("Unable to list IAM Roles")
	}

	instances := make([]*core.Instance, len(result.Roles))
	for i, r := range result.Roles {
		instances[i] = &core.Instance{
			Name: aws.StringValue(r.RoleName),
			ID:   aws.StringValue(r.RoleName),
		}
	}

	var links = map[string]string{}

	// This information is declared in the provider Schema, but not exposed through the public interface
	// Consider writing a script to automatically populate this by parsing the existing code
	var defaults = map[string]core.Default{
		"path":                  {Value: "/"},
		"force_detach_policies": {Value: "false", IsBool: true},
	}

	return core.ResourceDescription{
		Links:     links,
		Instances: instances,
		Defaults:  defaults,
	}
}
