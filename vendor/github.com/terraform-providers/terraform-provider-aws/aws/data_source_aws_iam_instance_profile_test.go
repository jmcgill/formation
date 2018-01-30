package aws

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAWSDataSourceIAMInstanceProfile_basic(t *testing.T) {
	roleName := fmt.Sprintf("test-datasource-user-%d", acctest.RandInt())
	profileName := fmt.Sprintf("test-datasource-user-%d", acctest.RandInt())

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceAwsIamInstanceProfileConfig(roleName, profileName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.aws_iam_instance_profile.test", "role_id"),
					resource.TestCheckResourceAttr("data.aws_iam_instance_profile.test", "path", "/testpath/"),
					resource.TestMatchResourceAttr("data.aws_iam_instance_profile.test", "arn",
						regexp.MustCompile("^arn:aws:iam::[0-9]{12}:instance-profile/testpath/"+profileName+"$")),
				),
			},
		},
	})
}

func testAccDatasourceAwsIamInstanceProfileConfig(roleName, profileName string) string {
	return fmt.Sprintf(`
resource "aws_iam_role" "test" {
	name = "%s"
	assume_role_policy = "{\"Version\":\"2012-10-17\",\"Statement\":[{\"Effect\":\"Allow\",\"Principal\":{\"Service\":[\"ec2.amazonaws.com\"]},\"Action\":[\"sts:AssumeRole\"]}]}"
}

resource "aws_iam_instance_profile" "test" {
	name = "%s"
	role = "${aws_iam_role.test.name}"
	path = "/testpath/"
}

data "aws_iam_instance_profile" "test" {
	name = "${aws_iam_instance_profile.test.name}"
}
`, roleName, profileName)
}
