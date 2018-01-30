package aws

import (
	"fmt"
	"log"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAWSSsmResourceDataSync_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAWSSsmResourceDataSyncDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSsmResourceDataSyncConfig(acctest.RandInt(), acctest.RandString(5)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSSsmResourceDataSyncExists("aws_ssm_resource_data_sync.foo"),
				),
			},
		},
	})
}

func testAccCheckAWSSsmResourceDataSyncDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*AWSClient).ssmconn

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aws_ssm_resource_data_sync" {
			continue
		}
		syncItem, err := findResourceDataSyncItem(conn, rs.Primary.Attributes["name"])
		if err != nil {
			return err
		}
		if syncItem != nil {
			return fmt.Errorf("Resource Data Sync (%s) found", rs.Primary.Attributes["name"])
		}
	}
	return nil
}

func testAccCheckAWSSsmResourceDataSyncExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		log.Println(s.RootModule().Resources)
		_, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}
		return nil
	}
}

func testAccSsmResourceDataSyncConfig(rInt int, rName string) string {
	return fmt.Sprintf(`
    resource "aws_s3_bucket" "hoge" {
      bucket = "tf-test-bucket-%d"
      region = "us-west-2"
      force_destroy = true
    }

    resource "aws_s3_bucket_policy" "hoge" {
    	bucket = "${aws_s3_bucket.hoge.bucket}"
    	policy = <<EOF
{
        "Version": "2012-10-17",
        "Statement": [
            {
                "Sid": "SSMBucketPermissionsCheck",
                "Effect": "Allow",
                "Principal": {
                    "Service": "ssm.amazonaws.com"
                },
                "Action": "s3:GetBucketAcl",
                "Resource": "arn:aws:s3:::tf-test-bucket-%d"
            },
            {
                "Sid": " SSMBucketDelivery",
                "Effect": "Allow",
                "Principal": {
                    "Service": "ssm.amazonaws.com"
                },
                "Action": "s3:PutObject",
                "Resource": ["arn:aws:s3:::tf-test-bucket-%d/*"],
                "Condition": {
                    "StringEquals": {
                        "s3:x-amz-acl": "bucket-owner-full-control"
                    }
                }
            }
          ]
      }
      EOF
    }

    resource "aws_ssm_resource_data_sync" "foo" {
      name = "tf-test-ssm-%s"
      s3_destination = {
        bucket_name = "${aws_s3_bucket.hoge.bucket}"
        region = "${aws_s3_bucket.hoge.region}"
      }
    }
    `, rInt, rInt, rInt, rName)
}
