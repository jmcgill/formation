package aws

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/athena"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAWSAthenaDatabase_basic(t *testing.T) {
	rInt := acctest.RandInt()
	dbName := acctest.RandString(8)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAWSAthenaDatabaseDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAthenaDatabaseConfig(rInt, dbName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSAthenaDatabaseExists("aws_athena_database.hoge"),
				),
			},
		},
	})
}

func TestAccAWSAthenaDatabase_destroyFailsIfTablesExist(t *testing.T) {
	rInt := acctest.RandInt()
	dbName := acctest.RandString(8)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAWSAthenaDatabaseDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAthenaDatabaseConfig(rInt, dbName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSAthenaDatabaseExists("aws_athena_database.hoge"),
					testAccAWSAthenaDatabaseCreateTables(dbName),
					testAccCheckAWSAthenaDatabaseDropFails(dbName),
					testAccAWSAthenaDatabaseDestroyTables(dbName),
				),
			},
		},
	})
}

func TestAccAWSAthenaDatabase_forceDestroyAlwaysSucceeds(t *testing.T) {
	rInt := acctest.RandInt()
	dbName := acctest.RandString(8)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAWSAthenaDatabaseDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAthenaDatabaseConfigForceDestroy(rInt, dbName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSAthenaDatabaseExists("aws_athena_database.hoge"),
					testAccAWSAthenaDatabaseCreateTables(dbName),
				),
			},
		},
	})
}

// StartQueryExecution requires OutputLocation but terraform destroy deleted S3 bucket as well.
// So temporary S3 bucket as OutputLocation is created to confirm whether the database is actually deleted.
func testAccCheckAWSAthenaDatabaseDestroy(s *terraform.State) error {
	athenaconn := testAccProvider.Meta().(*AWSClient).athenaconn
	s3conn := testAccProvider.Meta().(*AWSClient).s3conn
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aws_athena_database" {
			continue
		}

		rInt := acctest.RandInt()
		bucketName := fmt.Sprintf("tf-athena-db-%s-%d", rs.Primary.Attributes["name"], rInt)
		_, err := s3conn.CreateBucket(&s3.CreateBucketInput{
			Bucket: aws.String(bucketName),
		})
		if err != nil {
			return err
		}

		input := &athena.StartQueryExecutionInput{
			QueryString: aws.String(fmt.Sprint("show databases;")),
			ResultConfiguration: &athena.ResultConfiguration{
				OutputLocation: aws.String("s3://" + bucketName),
			},
		}

		resp, err := athenaconn.StartQueryExecution(input)
		if err != nil {
			return err
		}

		ers, err := queryExecutionResult(*resp.QueryExecutionId, athenaconn)
		if err != nil {
			return err
		}
		found := false
		dbName := rs.Primary.Attributes["name"]
		for _, row := range ers.Rows {
			for _, datum := range row.Data {
				if *datum.VarCharValue == dbName {
					found = true
				}
			}
		}
		if found {
			return fmt.Errorf("[DELETE ERROR] Athena failed to drop database: %s", dbName)
		}

		loresp, err := s3conn.ListObjectsV2(
			&s3.ListObjectsV2Input{
				Bucket: aws.String(bucketName),
			},
		)
		if err != nil {
			return fmt.Errorf("[DELETE ERROR] S3 Bucket list Objects err: %s", err)
		}

		objectsToDelete := make([]*s3.ObjectIdentifier, 0)

		if len(loresp.Contents) != 0 {
			for _, v := range loresp.Contents {
				objectsToDelete = append(objectsToDelete, &s3.ObjectIdentifier{
					Key: v.Key,
				})
			}
		}

		_, err = s3conn.DeleteObjects(&s3.DeleteObjectsInput{
			Bucket: aws.String(bucketName),
			Delete: &s3.Delete{
				Objects: objectsToDelete,
			},
		})
		if err != nil {
			return fmt.Errorf("[DELETE ERROR] S3 Bucket delete Objects err: %s", err)
		}

		_, err = s3conn.DeleteBucket(&s3.DeleteBucketInput{
			Bucket: aws.String(bucketName),
		})
		if err != nil {
			return fmt.Errorf("[DELETE ERROR] S3 Bucket delete Bucket err: %s", err)
		}

	}
	return nil
}

func testAccCheckAWSAthenaDatabaseExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("not found: %s, %v", name, s.RootModule().Resources)
		}
		return nil
	}
}

func testAccAWSAthenaDatabaseCreateTables(dbName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		bucketName, err := testAccAthenaDatabaseFindBucketName(s, dbName)
		if err != nil {
			return err
		}

		athenaconn := testAccProvider.Meta().(*AWSClient).athenaconn

		input := &athena.StartQueryExecutionInput{
			QueryExecutionContext: &athena.QueryExecutionContext{
				Database: aws.String(dbName),
			},
			QueryString: aws.String(fmt.Sprintf(
				"create external table foo (bar int) location 's3://%s/';", bucketName)),
			ResultConfiguration: &athena.ResultConfiguration{
				OutputLocation: aws.String("s3://" + bucketName),
			},
		}

		resp, err := athenaconn.StartQueryExecution(input)
		if err != nil {
			return err
		}

		_, err = queryExecutionResult(*resp.QueryExecutionId, athenaconn)
		if err != nil {
			return err
		}

		return nil
	}
}

func testAccAWSAthenaDatabaseDestroyTables(dbName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		bucketName, err := testAccAthenaDatabaseFindBucketName(s, dbName)
		if err != nil {
			return err
		}

		athenaconn := testAccProvider.Meta().(*AWSClient).athenaconn

		input := &athena.StartQueryExecutionInput{
			QueryExecutionContext: &athena.QueryExecutionContext{
				Database: aws.String(dbName),
			},
			QueryString: aws.String("drop table foo;"),
			ResultConfiguration: &athena.ResultConfiguration{
				OutputLocation: aws.String("s3://" + bucketName),
			},
		}

		resp, err := athenaconn.StartQueryExecution(input)
		if err != nil {
			return err
		}

		_, err = queryExecutionResult(*resp.QueryExecutionId, athenaconn)
		if err != nil {
			return err
		}

		return nil
	}
}

func testAccCheckAWSAthenaDatabaseDropFails(dbName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		bucketName, err := testAccAthenaDatabaseFindBucketName(s, dbName)
		if err != nil {
			return err
		}

		athenaconn := testAccProvider.Meta().(*AWSClient).athenaconn

		input := &athena.StartQueryExecutionInput{
			QueryExecutionContext: &athena.QueryExecutionContext{
				Database: aws.String(dbName),
			},
			QueryString: aws.String(fmt.Sprintf("drop database %s;", dbName)),
			ResultConfiguration: &athena.ResultConfiguration{
				OutputLocation: aws.String("s3://" + bucketName),
			},
		}

		resp, err := athenaconn.StartQueryExecution(input)
		if err != nil {
			return err
		}

		_, err = queryExecutionResult(*resp.QueryExecutionId, athenaconn)
		if err == nil {
			return fmt.Errorf("drop database unexpectedly succeeded for a database with tables")
		}

		return nil
	}
}

func testAccAthenaDatabaseFindBucketName(s *terraform.State, dbName string) (bucket string, err error) {
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aws_athena_database" && rs.Primary.Attributes["name"] == dbName {
			bucket = rs.Primary.Attributes["bucket"]
			break
		}
	}

	if bucket == "" {
		err = fmt.Errorf("cannot find database %s", dbName)
	}

	return bucket, err
}

func testAccAthenaDatabaseConfig(randInt int, dbName string) string {
	return fmt.Sprintf(`
    resource "aws_s3_bucket" "hoge" {
      bucket = "tf-athena-db-%s-%d"
      force_destroy = true
    }

    resource "aws_athena_database" "hoge" {
      name = "%s"
      bucket = "${aws_s3_bucket.hoge.bucket}"
    }
    `, dbName, randInt, dbName)
}

func testAccAthenaDatabaseConfigForceDestroy(randInt int, dbName string) string {
	return fmt.Sprintf(`
    resource "aws_s3_bucket" "hoge" {
      bucket = "tf-athena-db-%s-%d"
      force_destroy = true
    }

    resource "aws_athena_database" "hoge" {
      name = "%s"
      bucket = "${aws_s3_bucket.hoge.bucket}"
	  force_destroy = true
    }
    `, dbName, randInt, dbName)
}
