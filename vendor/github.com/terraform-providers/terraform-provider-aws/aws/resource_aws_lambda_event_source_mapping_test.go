package aws

import (
	"fmt"
	"regexp"
	"strconv"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAWSLambdaEventSourceMapping_basic(t *testing.T) {
	var conf lambda.EventSourceMappingConfiguration

	rString := acctest.RandString(8)
	roleName := fmt.Sprintf("tf_acc_role_lambda_esm_basic_%s", rString)
	policyName := fmt.Sprintf("tf_acc_policy_lambda_esm_basic_%s", rString)
	attName := fmt.Sprintf("tf_acc_att_lambda_esm_basic_%s", rString)
	streamName := fmt.Sprintf("tf_acc_stream_lambda_esm_basic_%s", rString)
	funcName := fmt.Sprintf("tf_acc_lambda_esm_basic_%s", rString)
	uFuncName := fmt.Sprintf("tf_acc_lambda_esm_basic_updated_%s", rString)
	uFuncArnRe := regexp.MustCompile(":" + uFuncName + "$")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLambdaEventSourceMappingDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSLambdaEventSourceMappingConfig(roleName, policyName, attName, streamName, funcName, uFuncName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAwsLambdaEventSourceMappingExists("aws_lambda_event_source_mapping.lambda_event_source_mapping_test", &conf),
					testAccCheckAWSLambdaEventSourceMappingAttributes(&conf),
				),
			},
			{
				Config: testAccAWSLambdaEventSourceMappingConfigUpdate(roleName, policyName, attName, streamName, funcName, uFuncName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAwsLambdaEventSourceMappingExists("aws_lambda_event_source_mapping.lambda_event_source_mapping_test", &conf),
					resource.TestCheckResourceAttr("aws_lambda_event_source_mapping.lambda_event_source_mapping_test",
						"batch_size", strconv.Itoa(200)),
					resource.TestCheckResourceAttr("aws_lambda_event_source_mapping.lambda_event_source_mapping_test",
						"enabled", strconv.FormatBool(false)),
					resource.TestMatchResourceAttr("aws_lambda_event_source_mapping.lambda_event_source_mapping_test",
						"function_arn", uFuncArnRe),
				),
			},
		},
	})
}

func TestAccAWSLambdaEventSourceMapping_importBasic(t *testing.T) {
	resourceName := "aws_lambda_event_source_mapping.lambda_event_source_mapping_test"

	rString := acctest.RandString(8)
	roleName := fmt.Sprintf("tf_acc_role_lambda_esm_import_%s", rString)
	policyName := fmt.Sprintf("tf_acc_policy_lambda_esm_import_%s", rString)
	attName := fmt.Sprintf("tf_acc_att_lambda_esm_import_%s", rString)
	streamName := fmt.Sprintf("tf_acc_stream_lambda_esm_import_%s", rString)
	funcName := fmt.Sprintf("tf_acc_lambda_esm_import_%s", rString)
	uFuncName := fmt.Sprintf("tf_acc_lambda_esm_import_updated_%s", rString)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLambdaEventSourceMappingDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSLambdaEventSourceMappingConfig(roleName, policyName, attName, streamName, funcName, uFuncName),
			},

			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"enabled", "starting_position"},
			},
		},
	})
}

func TestAccAWSLambdaEventSourceMapping_disappears(t *testing.T) {
	var conf lambda.EventSourceMappingConfiguration

	rString := acctest.RandString(8)
	roleName := fmt.Sprintf("tf_acc_role_lambda_esm_import_%s", rString)
	policyName := fmt.Sprintf("tf_acc_policy_lambda_esm_import_%s", rString)
	attName := fmt.Sprintf("tf_acc_att_lambda_esm_import_%s", rString)
	streamName := fmt.Sprintf("tf_acc_stream_lambda_esm_import_%s", rString)
	funcName := fmt.Sprintf("tf_acc_lambda_esm_import_%s", rString)
	uFuncName := fmt.Sprintf("tf_acc_lambda_esm_import_updated_%s", rString)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLambdaEventSourceMappingDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSLambdaEventSourceMappingConfig(roleName, policyName, attName, streamName, funcName, uFuncName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAwsLambdaEventSourceMappingExists("aws_lambda_event_source_mapping.lambda_event_source_mapping_test", &conf),
					testAccCheckAWSLambdaEventSourceMappingDisappears(&conf),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccCheckAWSLambdaEventSourceMappingDisappears(conf *lambda.EventSourceMappingConfiguration) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := testAccProvider.Meta().(*AWSClient).lambdaconn

		params := &lambda.DeleteEventSourceMappingInput{
			UUID: conf.UUID,
		}

		_, err := conn.DeleteEventSourceMapping(params)
		if err != nil {
			if err != nil {
				return err
			}
		}

		return resource.Retry(10*time.Minute, func() *resource.RetryError {
			params := &lambda.GetEventSourceMappingInput{
				UUID: conf.UUID,
			}
			_, err := conn.GetEventSourceMapping(params)
			if err != nil {
				cgw, ok := err.(awserr.Error)
				if ok && cgw.Code() == "ResourceNotFoundException" {
					return nil
				}
				return resource.NonRetryableError(
					fmt.Errorf("Error retrieving Lambda Event Source Mapping: %s", err))
			}
			return resource.RetryableError(fmt.Errorf(
				"Waiting for Lambda Event Source Mapping: %v", conf.UUID))
		})
	}
}

func testAccCheckLambdaEventSourceMappingDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*AWSClient).lambdaconn

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aws_lambda_event_source_mapping" {
			continue
		}

		_, err := conn.GetEventSourceMapping(&lambda.GetEventSourceMappingInput{
			UUID: aws.String(rs.Primary.ID),
		})

		if err == nil {
			return fmt.Errorf("Lambda event source mapping was not deleted")
		}

	}

	return nil

}

func testAccCheckAwsLambdaEventSourceMappingExists(n string, mapping *lambda.EventSourceMappingConfiguration) resource.TestCheckFunc {
	// Wait for IAM role
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Lambda event source mapping not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Lambda event source mapping ID not set")
		}

		conn := testAccProvider.Meta().(*AWSClient).lambdaconn

		params := &lambda.GetEventSourceMappingInput{
			UUID: aws.String(rs.Primary.ID),
		}

		getSourceMappingConfiguration, err := conn.GetEventSourceMapping(params)
		if err != nil {
			return err
		}

		*mapping = *getSourceMappingConfiguration

		return nil
	}
}

func testAccCheckAWSLambdaEventSourceMappingAttributes(mapping *lambda.EventSourceMappingConfiguration) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		uuid := *mapping.UUID
		if uuid == "" {
			return fmt.Errorf("Could not read Lambda event source mapping's UUID")
		}

		return nil
	}
}

func testAccAWSLambdaEventSourceMappingConfig(roleName, policyName, attName, streamName,
	funcName, uFuncName string) string {
	return fmt.Sprintf(`
resource "aws_iam_role" "iam_for_lambda" {
    name = "%s"
    assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Principal": {
        "Service": "lambda.amazonaws.com"
      },
      "Effect": "Allow",
      "Sid": ""
    }
  ]
}
EOF
}

resource "aws_iam_policy" "policy_for_role" {
    name = "%s"
    path = "/"
    description = "IAM policy for for Lamda event mapping testing"
    policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
      {
          "Effect": "Allow",
          "Action": [
            "kinesis:GetRecords",
            "kinesis:GetShardIterator",
            "kinesis:DescribeStream"
          ],
          "Resource": "*"
      },
      {
          "Effect": "Allow",
          "Action": [
            "kinesis:ListStreams"
          ],
          "Resource": "*"
      }
  ]
}
EOF
}

resource "aws_iam_policy_attachment" "policy_attachment_for_role" {
    name = "%s"
    roles = ["${aws_iam_role.iam_for_lambda.name}"]
    policy_arn = "${aws_iam_policy.policy_for_role.arn}"
}

resource "aws_kinesis_stream" "kinesis_stream_test" {
    name = "%s"
    shard_count = 1
}

resource "aws_lambda_function" "lambda_function_test_create" {
    filename = "test-fixtures/lambdatest.zip"
    function_name = "%s"
    role = "${aws_iam_role.iam_for_lambda.arn}"
    handler = "exports.example"
    runtime = "nodejs4.3"
}

resource "aws_lambda_function" "lambda_function_test_update" {
    filename = "test-fixtures/lambdatest.zip"
    function_name = "%s"
    role = "${aws_iam_role.iam_for_lambda.arn}"
    handler = "exports.example"
    runtime = "nodejs4.3"
}

resource "aws_lambda_event_source_mapping" "lambda_event_source_mapping_test" {
		batch_size = 100
		event_source_arn = "${aws_kinesis_stream.kinesis_stream_test.arn}"
		enabled = true
		depends_on = ["aws_iam_policy_attachment.policy_attachment_for_role"]
		function_name = "${aws_lambda_function.lambda_function_test_create.arn}"
		starting_position = "TRIM_HORIZON"
}`, roleName, policyName, attName, streamName, funcName, uFuncName)
}

func testAccAWSLambdaEventSourceMappingConfigUpdate(roleName, policyName, attName, streamName,
	funcName, uFuncName string) string {
	return fmt.Sprintf(`
resource "aws_iam_role" "iam_for_lambda" {
    name = "%s"
    assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Principal": {
        "Service": "lambda.amazonaws.com"
      },
      "Effect": "Allow",
      "Sid": ""
    }
  ]
}
EOF
}

resource "aws_iam_policy" "policy_for_role" {
    name = "%s"
    path = "/"
    description = "IAM policy for for Lamda event mapping testing"
    policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
      {
          "Effect": "Allow",
          "Action": [
            "kinesis:GetRecords",
            "kinesis:GetShardIterator",
            "kinesis:DescribeStream"
          ],
          "Resource": "*"
      },
      {
          "Effect": "Allow",
          "Action": [
            "kinesis:ListStreams"
          ],
          "Resource": "*"
      }
  ]
}
EOF
}

resource "aws_iam_policy_attachment" "policy_attachment_for_role" {
    name = "%s"
    roles = ["${aws_iam_role.iam_for_lambda.name}"]
    policy_arn = "${aws_iam_policy.policy_for_role.arn}"
}

resource "aws_kinesis_stream" "kinesis_stream_test" {
    name = "%s"
    shard_count = 1
}

resource "aws_lambda_function" "lambda_function_test_create" {
    filename = "test-fixtures/lambdatest.zip"
    function_name = "%s"
    role = "${aws_iam_role.iam_for_lambda.arn}"
    handler = "exports.example"
    runtime = "nodejs4.3"
}

resource "aws_lambda_function" "lambda_function_test_update" {
    filename = "test-fixtures/lambdatest.zip"
    function_name = "%s"
    role = "${aws_iam_role.iam_for_lambda.arn}"
    handler = "exports.example"
    runtime = "nodejs4.3"
}

resource "aws_lambda_event_source_mapping" "lambda_event_source_mapping_test" {
		batch_size = 200
		event_source_arn = "${aws_kinesis_stream.kinesis_stream_test.arn}"
		enabled = false
		depends_on = ["aws_iam_policy_attachment.policy_attachment_for_role"]
		function_name = "${aws_lambda_function.lambda_function_test_update.arn}"
		starting_position = "TRIM_HORIZON"
}`, roleName, policyName, attName, streamName, funcName, uFuncName)
}
