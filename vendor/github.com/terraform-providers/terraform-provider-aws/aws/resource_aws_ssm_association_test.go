package aws

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAWSSSMAssociation_basic(t *testing.T) {
	name := acctest.RandString(10)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAWSSSMAssociationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSSSMAssociationBasicConfig(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSSSMAssociationExists("aws_ssm_association.foo"),
				),
			},
		},
	})
}

func TestAccAWSSSMAssociation_withTargets(t *testing.T) {
	name := acctest.RandString(10)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAWSSSMAssociationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSSSMAssociationBasicConfigWithTargets(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSSSMAssociationExists("aws_ssm_association.foo"),
				),
			},
		},
	})
}

func TestAccAWSSSMAssociation_withParameters(t *testing.T) {
	name := acctest.RandString(10)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAWSSSMAssociationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSSSMAssociationBasicConfigWithParameters(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSSSMAssociationExists("aws_ssm_association.foo"),
					resource.TestCheckResourceAttr(
						"aws_ssm_association.foo", "parameters.Directory", "myWorkSpace"),
				),
			},
			{
				Config: testAccAWSSSMAssociationBasicConfigWithParametersUpdated(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSSSMAssociationExists("aws_ssm_association.foo"),
					resource.TestCheckResourceAttr(
						"aws_ssm_association.foo", "parameters.Directory", "myWorkSpaceUpdated"),
				),
			},
		},
	})
}

func TestAccAWSSSMAssociation_withAssociationName(t *testing.T) {
	assocName1 := acctest.RandString(10)
	assocName2 := acctest.RandString(10)
	rName := acctest.RandString(5)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAWSSSMAssociationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSSSMAssociationBasicConfigWithAssociationName(rName, assocName1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSSSMAssociationExists("aws_ssm_association.foo"),
					resource.TestCheckResourceAttr(
						"aws_ssm_association.foo", "association_name", assocName1),
				),
			},
			{
				Config: testAccAWSSSMAssociationBasicConfigWithAssociationName(rName, assocName2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSSSMAssociationExists("aws_ssm_association.foo"),
					resource.TestCheckResourceAttr(
						"aws_ssm_association.foo", "association_name", assocName2),
				),
			},
		},
	})
}

func TestAccAWSSSMAssociation_withDocumentVersion(t *testing.T) {
	name := acctest.RandString(10)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAWSSSMAssociationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSSSMAssociationBasicConfigWithDocumentVersion(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSSSMAssociationExists("aws_ssm_association.foo"),
					resource.TestCheckResourceAttr(
						"aws_ssm_association.foo", "document_version", "1"),
				),
			},
		},
	})
}

func TestAccAWSSSMAssociation_withOutputLocation(t *testing.T) {
	name := acctest.RandString(10)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAWSSSMAssociationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSSSMAssociationBasicConfigWithOutPutLocation(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSSSMAssociationExists("aws_ssm_association.foo"),
					resource.TestCheckResourceAttr(
						"aws_ssm_association.foo", "output_location.0.s3_bucket_name", fmt.Sprintf("tf-acc-test-ssmoutput-%s", name)),
					resource.TestCheckResourceAttr(
						"aws_ssm_association.foo", "output_location.0.s3_key_prefix", "SSMAssociation"),
				),
			},
			{
				Config: testAccAWSSSMAssociationBasicConfigWithOutPutLocationUpdateBucketName(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSSSMAssociationExists("aws_ssm_association.foo"),
					resource.TestCheckResourceAttr(
						"aws_ssm_association.foo", "output_location.0.s3_bucket_name", fmt.Sprintf("tf-acc-test-ssmoutput-updated-%s", name)),
					resource.TestCheckResourceAttr(
						"aws_ssm_association.foo", "output_location.0.s3_key_prefix", "SSMAssociation"),
				),
			},
			{
				Config: testAccAWSSSMAssociationBasicConfigWithOutPutLocationUpdateKeyPrefix(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSSSMAssociationExists("aws_ssm_association.foo"),
					resource.TestCheckResourceAttr(
						"aws_ssm_association.foo", "output_location.0.s3_bucket_name", fmt.Sprintf("tf-acc-test-ssmoutput-updated-%s", name)),
					resource.TestCheckResourceAttr(
						"aws_ssm_association.foo", "output_location.0.s3_key_prefix", "UpdatedAssociation"),
				),
			},
		},
	})
}

func TestAccAWSSSMAssociation_withScheduleExpression(t *testing.T) {
	name := acctest.RandString(10)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAWSSSMAssociationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSSSMAssociationBasicConfigWithScheduleExpression(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSSSMAssociationExists("aws_ssm_association.foo"),
					resource.TestCheckResourceAttr(
						"aws_ssm_association.foo", "schedule_expression", "cron(0 16 ? * TUE *)"),
				),
			},
			{
				Config: testAccAWSSSMAssociationBasicConfigWithScheduleExpressionUpdated(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSSSMAssociationExists("aws_ssm_association.foo"),
					resource.TestCheckResourceAttr(
						"aws_ssm_association.foo", "schedule_expression", "cron(0 16 ? * WED *)"),
				),
			},
		},
	})
}

func testAccCheckAWSSSMAssociationExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No SSM Assosciation ID is set")
		}

		conn := testAccProvider.Meta().(*AWSClient).ssmconn

		_, err := conn.DescribeAssociation(&ssm.DescribeAssociationInput{
			AssociationId: aws.String(rs.Primary.Attributes["association_id"]),
		})

		if err != nil {
			if wserr, ok := err.(awserr.Error); ok && wserr.Code() == "AssociationDoesNotExist" {
				return nil
			}
			return err
		}

		return nil
	}
}

func testAccCheckAWSSSMAssociationDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*AWSClient).ssmconn

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aws_ssm_association" {
			continue
		}

		out, err := conn.DescribeAssociation(&ssm.DescribeAssociationInput{
			AssociationId: aws.String(rs.Primary.Attributes["association_id"]),
		})

		if err != nil {
			if wserr, ok := err.(awserr.Error); ok && wserr.Code() == "AssociationDoesNotExist" {
				return nil
			}
			return err
		}

		if out != nil {
			return fmt.Errorf("Expected AWS SSM Association to be gone, but was still found")
		}
	}

	return fmt.Errorf("Default error in SSM Association Test")
}

func testAccAWSSSMAssociationBasicConfigWithParameters(rName string) string {
	return fmt.Sprintf(`
resource "aws_ssm_document" "foo_document" {
  name = "test_document_association-%s",
  document_type = "Command"
  content = <<-DOC
  {
    "schemaVersion": "1.2",
    "description": "Check ip configuration of a Linux instance.",
    "parameters": {
	  "Directory": {
		"description":"(Optional) The path to the working directory on your instance.",
		"default":"",
		"type": "String",
		"maxChars": 4096
	  }
	},
    "runtimeConfig": {
      "aws:runShellScript": {
        "properties": [
          {
            "id": "0.aws:runShellScript",
            "runCommand": ["ifconfig"]
          }
        ]
      }
    }
  }
  DOC
}

resource "aws_ssm_association" "foo" {
  name = "${aws_ssm_document.foo_document.name}",
  parameters {
  	Directory = "myWorkSpace"
  }
  targets {
    key = "tag:Name"
    values = ["acceptanceTest"]
  }
}`, rName)
}

func testAccAWSSSMAssociationBasicConfigWithParametersUpdated(rName string) string {
	return fmt.Sprintf(`
resource "aws_ssm_document" "foo_document" {
  name = "test_document_association-%s",
  document_type = "Command"
  content = <<-DOC
  {
    "schemaVersion": "1.2",
    "description": "Check ip configuration of a Linux instance.",
    "parameters": {
	  "Directory": {
		"description":"(Optional) The path to the working directory on your instance.",
		"default":"",
		"type": "String",
		"maxChars": 4096
	  }
	},
    "runtimeConfig": {
      "aws:runShellScript": {
        "properties": [
          {
            "id": "0.aws:runShellScript",
            "runCommand": ["ifconfig"]
          }
        ]
      }
    }
  }
  DOC
}

resource "aws_ssm_association" "foo" {
  name = "${aws_ssm_document.foo_document.name}",
  parameters {
  	Directory = "myWorkSpaceUpdated"
  }
  targets {
    key = "tag:Name"
    values = ["acceptanceTest"]
  }
}`, rName)
}

func testAccAWSSSMAssociationBasicConfigWithTargets(rName string) string {
	return fmt.Sprintf(`
resource "aws_ssm_document" "foo_document" {
  name = "test_document_association-%s",
  document_type = "Command"
  content = <<DOC
  {
    "schemaVersion": "1.2",
    "description": "Check ip configuration of a Linux instance.",
    "parameters": {

    },
    "runtimeConfig": {
      "aws:runShellScript": {
        "properties": [
          {
            "id": "0.aws:runShellScript",
            "runCommand": ["ifconfig"]
          }
        ]
      }
    }
  }
DOC
}

resource "aws_ssm_association" "foo" {
  name = "${aws_ssm_document.foo_document.name}",
  targets {
    key = "tag:Name"
    values = ["acceptanceTest"]
  }
}`, rName)
}

func testAccAWSSSMAssociationBasicConfig(rName string) string {
	return fmt.Sprintf(`
resource "aws_security_group" "tf_test_foo" {
  name = "tf_test_foo-%s"
  description = "foo"
  ingress {
    protocol = "icmp"
    from_port = -1
    to_port = -1
    cidr_blocks = ["0.0.0.0/0"]
  }
}

resource "aws_instance" "foo" {
  ami = "ami-4fccb37f"
  availability_zone = "us-west-2a"
  instance_type = "m1.small"
  security_groups = ["${aws_security_group.tf_test_foo.name}"]
}

resource "aws_ssm_document" "foo_document" {
  name    = "test_document_association-%s",
	document_type = "Command"
  content = <<DOC
  {
    "schemaVersion": "1.2",
    "description": "Check ip configuration of a Linux instance.",
    "parameters": {

    },
    "runtimeConfig": {
      "aws:runShellScript": {
        "properties": [
          {
            "id": "0.aws:runShellScript",
            "runCommand": ["ifconfig"]
          }
        ]
      }
    }
  }
DOC
}

resource "aws_ssm_association" "foo" {
  name        = "test_document_association-%s",
  instance_id = "${aws_instance.foo.id}"
}
`, rName, rName, rName)
}

func testAccAWSSSMAssociationBasicConfigWithDocumentVersion(rName string) string {
	return fmt.Sprintf(`
resource "aws_security_group" "tf_test_foo" {
  name = "tf_test_foo-%s"
  description = "foo"
  ingress {
    protocol = "icmp"
    from_port = -1
    to_port = -1
    cidr_blocks = ["0.0.0.0/0"]
  }
}

resource "aws_ssm_document" "foo_document" {
  name    = "test_document_association-%s",
	document_type = "Command"
  content = <<DOC
  {
    "schemaVersion": "1.2",
    "description": "Check ip configuration of a Linux instance.",
    "parameters": {

    },
    "runtimeConfig": {
      "aws:runShellScript": {
        "properties": [
          {
            "id": "0.aws:runShellScript",
            "runCommand": ["ifconfig"]
          }
        ]
      }
    }
  }
DOC
}

resource "aws_ssm_association" "foo" {
  name        = "test_document_association-%s",
  document_version = "${aws_ssm_document.foo_document.latest_version}"
  targets {
    key = "tag:Name"
    values = ["acceptanceTest"]
  }
}
`, rName, rName, rName)
}

func testAccAWSSSMAssociationBasicConfigWithScheduleExpression(rName string) string {
	return fmt.Sprintf(`
resource "aws_ssm_document" "foo_document" {
  name = "test_document_association-%s",
  document_type = "Command"
  content = <<DOC
  {
    "schemaVersion": "1.2",
    "description": "Check ip configuration of a Linux instance.",
    "parameters": {

    },
    "runtimeConfig": {
      "aws:runShellScript": {
        "properties": [
          {
            "id": "0.aws:runShellScript",
            "runCommand": ["ifconfig"]
          }
        ]
      }
    }
  }
DOC
}

resource "aws_ssm_association" "foo" {
  name = "${aws_ssm_document.foo_document.name}",
  schedule_expression = "cron(0 16 ? * TUE *)"
  targets {
    key = "tag:Name"
    values = ["acceptanceTest"]
  }
}`, rName)
}

func testAccAWSSSMAssociationBasicConfigWithScheduleExpressionUpdated(rName string) string {
	return fmt.Sprintf(`
resource "aws_ssm_document" "foo_document" {
  name = "test_document_association-%s",
  document_type = "Command"
  content = <<DOC
  {
    "schemaVersion": "1.2",
    "description": "Check ip configuration of a Linux instance.",
    "parameters": {

    },
    "runtimeConfig": {
      "aws:runShellScript": {
        "properties": [
          {
            "id": "0.aws:runShellScript",
            "runCommand": ["ifconfig"]
          }
        ]
      }
    }
  }
DOC
}

resource "aws_ssm_association" "foo" {
  name = "${aws_ssm_document.foo_document.name}",
  schedule_expression = "cron(0 16 ? * WED *)"
  targets {
    key = "tag:Name"
    values = ["acceptanceTest"]
  }
}`, rName)
}

func testAccAWSSSMAssociationBasicConfigWithOutPutLocation(rName string) string {
	return fmt.Sprintf(`
resource "aws_s3_bucket" "output_location" {
  bucket = "tf-acc-test-ssmoutput-%s"
  force_destroy = true
}

resource "aws_ssm_document" "foo_document" {
  name = "test_document_association-%s",
  document_type = "Command"
  content = <<DOC
  {
    "schemaVersion": "1.2",
    "description": "Check ip configuration of a Linux instance.",
    "parameters": {

    },
    "runtimeConfig": {
      "aws:runShellScript": {
        "properties": [
          {
            "id": "0.aws:runShellScript",
            "runCommand": ["ifconfig"]
          }
        ]
      }
    }
  }
DOC
}

resource "aws_ssm_association" "foo" {
  name = "${aws_ssm_document.foo_document.name}",
  targets {
    key = "tag:Name"
    values = ["acceptanceTest"]
  }
  output_location {
    s3_bucket_name = "${aws_s3_bucket.output_location.id}"
    s3_key_prefix = "SSMAssociation"
  }
}`, rName, rName)
}

func testAccAWSSSMAssociationBasicConfigWithOutPutLocationUpdateBucketName(rName string) string {
	return fmt.Sprintf(`
resource "aws_s3_bucket" "output_location" {
  bucket = "tf-acc-test-ssmoutput-%s"
  force_destroy = true
}

resource "aws_s3_bucket" "output_location_updated" {
  bucket = "tf-acc-test-ssmoutput-updated-%s"
  force_destroy = true
}

resource "aws_ssm_document" "foo_document" {
  name = "test_document_association-%s",
  document_type = "Command"
  content = <<DOC
  {
    "schemaVersion": "1.2",
    "description": "Check ip configuration of a Linux instance.",
    "parameters": {

    },
    "runtimeConfig": {
      "aws:runShellScript": {
        "properties": [
          {
            "id": "0.aws:runShellScript",
            "runCommand": ["ifconfig"]
          }
        ]
      }
    }
  }
DOC
}

resource "aws_ssm_association" "foo" {
  name = "${aws_ssm_document.foo_document.name}",
  targets {
    key = "tag:Name"
    values = ["acceptanceTest"]
  }
  output_location {
    s3_bucket_name = "${aws_s3_bucket.output_location_updated.id}"
    s3_key_prefix = "SSMAssociation"
  }
}`, rName, rName, rName)
}

func testAccAWSSSMAssociationBasicConfigWithOutPutLocationUpdateKeyPrefix(rName string) string {
	return fmt.Sprintf(`
resource "aws_s3_bucket" "output_location" {
  bucket = "tf-acc-test-ssmoutput-%s"
  force_destroy = true
}

resource "aws_s3_bucket" "output_location_updated" {
  bucket = "tf-acc-test-ssmoutput-updated-%s"
  force_destroy = true
}

resource "aws_ssm_document" "foo_document" {
  name = "test_document_association-%s",
  document_type = "Command"
  content = <<DOC
  {
    "schemaVersion": "1.2",
    "description": "Check ip configuration of a Linux instance.",
    "parameters": {

    },
    "runtimeConfig": {
      "aws:runShellScript": {
        "properties": [
          {
            "id": "0.aws:runShellScript",
            "runCommand": ["ifconfig"]
          }
        ]
      }
    }
  }
DOC
}

resource "aws_ssm_association" "foo" {
  name = "${aws_ssm_document.foo_document.name}",
  targets {
    key = "tag:Name"
    values = ["acceptanceTest"]
  }
  output_location {
    s3_bucket_name = "${aws_s3_bucket.output_location_updated.id}"
    s3_key_prefix = "UpdatedAssociation"
  }
}`, rName, rName, rName)
}

func testAccAWSSSMAssociationBasicConfigWithAssociationName(rName, assocName string) string {
	return fmt.Sprintf(`
resource "aws_ssm_document" "foo_document" {
  name = "test_document_association-%s",
  document_type = "Command"
  content = <<DOC
  {
    "schemaVersion": "1.2",
    "description": "Check ip configuration of a Linux instance.",
    "parameters": {
    },
    "runtimeConfig": {
      "aws:runShellScript": {
        "properties": [
          {
            "id": "0.aws:runShellScript",
            "runCommand": ["ifconfig"]
          }
        ]
      }
    }
  }
DOC
}

resource "aws_ssm_association" "foo" {
  name = "${aws_ssm_document.foo_document.name}",
  association_name = "%s"
  targets {
    key = "tag:Name"
    values = ["acceptanceTest"]
  }
}
`, rName, assocName)
}
