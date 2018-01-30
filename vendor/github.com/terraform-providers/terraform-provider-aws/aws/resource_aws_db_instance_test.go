package aws

import (
	"fmt"
	"log"
	"math/rand"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/rds"
)

func init() {
	resource.AddTestSweepers("aws_db_instance", &resource.Sweeper{
		Name: "aws_db_instance",
		F:    testSweepDbInstances,
	})
}

func testSweepDbInstances(region string) error {
	client, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting client: %s", err)
	}
	conn := client.(*AWSClient).rdsconn

	prefixes := []string{
		"foobarbaz-test-terraform-",
		"foobarbaz-enhanced-monitoring-",
		"mydb-rds-",
		"terraform-",
		"tf-",
	}

	err = conn.DescribeDBInstancesPages(&rds.DescribeDBInstancesInput{}, func(out *rds.DescribeDBInstancesOutput, lastPage bool) bool {
		for _, dbi := range out.DBInstances {
			hasPrefix := false
			for _, prefix := range prefixes {
				if strings.HasPrefix(*dbi.DBInstanceIdentifier, prefix) {
					hasPrefix = true
				}
			}
			if !hasPrefix {
				continue
			}
			log.Printf("[INFO] Deleting DB instance: %s", *dbi.DBInstanceIdentifier)

			_, err := conn.DeleteDBInstance(&rds.DeleteDBInstanceInput{
				DBInstanceIdentifier: dbi.DBInstanceIdentifier,
				SkipFinalSnapshot:    aws.Bool(true),
			})
			if err != nil {
				log.Printf("[ERROR] Failed to delete DB instance %s: %s",
					*dbi.DBInstanceIdentifier, err)
				continue
			}

			err = waitUntilAwsDbInstanceIsDeleted(*dbi.DBInstanceIdentifier, conn, 40*time.Minute)
			if err != nil {
				log.Printf("[ERROR] Failure while waiting for DB instance %s to be deleted: %s",
					*dbi.DBInstanceIdentifier, err)
			}
		}
		return !lastPage
	})
	if err != nil {
		return fmt.Errorf("Error retrieving DB instances: %s", err)
	}

	return nil
}

func TestAccAWSDBInstance_basic(t *testing.T) {
	var v rds.DBInstance

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAWSDBInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSDBInstanceConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSDBInstanceExists("aws_db_instance.bar", &v),
					testAccCheckAWSDBInstanceAttributes(&v),
					resource.TestCheckResourceAttr(
						"aws_db_instance.bar", "allocated_storage", "10"),
					resource.TestCheckResourceAttr(
						"aws_db_instance.bar", "engine", "mysql"),
					resource.TestCheckResourceAttr(
						"aws_db_instance.bar", "license_model", "general-public-license"),
					resource.TestCheckResourceAttr(
						"aws_db_instance.bar", "instance_class", "db.t2.micro"),
					resource.TestCheckResourceAttr(
						"aws_db_instance.bar", "name", "baz"),
					resource.TestCheckResourceAttr(
						"aws_db_instance.bar", "username", "foo"),
					resource.TestCheckResourceAttr(
						"aws_db_instance.bar", "parameter_group_name", "default.mysql5.6"),
					resource.TestCheckResourceAttrSet("aws_db_instance.bar", "hosted_zone_id"),
					resource.TestCheckResourceAttrSet("aws_db_instance.bar", "ca_cert_identifier"),
					resource.TestCheckResourceAttrSet(
						"aws_db_instance.bar", "resource_id"),
				),
			},
		},
	})
}

func TestAccAWSDBInstance_namePrefix(t *testing.T) {
	var v rds.DBInstance

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAWSDBInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSDBInstanceConfig_namePrefix,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSDBInstanceExists("aws_db_instance.test", &v),
					testAccCheckAWSDBInstanceAttributes(&v),
					resource.TestMatchResourceAttr(
						"aws_db_instance.test", "identifier", regexp.MustCompile("^tf-test-")),
				),
			},
		},
	})
}

func TestAccAWSDBInstance_generatedName(t *testing.T) {
	var v rds.DBInstance

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAWSDBInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSDBInstanceConfig_generatedName,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSDBInstanceExists("aws_db_instance.test", &v),
					testAccCheckAWSDBInstanceAttributes(&v),
				),
			},
		},
	})
}

func TestAccAWSDBInstance_kmsKey(t *testing.T) {
	var v rds.DBInstance
	keyRegex := regexp.MustCompile("^arn:aws:kms:")

	ri := rand.New(rand.NewSource(time.Now().UnixNano())).Int()
	config := fmt.Sprintf(testAccAWSDBInstanceConfigKmsKeyId, ri)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAWSDBInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSDBInstanceExists("aws_db_instance.bar", &v),
					testAccCheckAWSDBInstanceAttributes(&v),
					resource.TestMatchResourceAttr(
						"aws_db_instance.bar", "kms_key_id", keyRegex),
				),
			},
		},
	})
}

func TestAccAWSDBInstance_subnetGroup(t *testing.T) {
	var v rds.DBInstance
	rName := acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAWSDBInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSDBInstanceConfigWithSubnetGroup(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSDBInstanceExists("aws_db_instance.bar", &v),
					resource.TestCheckResourceAttr(
						"aws_db_instance.bar", "db_subnet_group_name", "foo-"+rName),
				),
			},
			{
				Config: testAccAWSDBInstanceConfigWithSubnetGroupUpdated(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSDBInstanceExists("aws_db_instance.bar", &v),
					resource.TestCheckResourceAttr(
						"aws_db_instance.bar", "db_subnet_group_name", "bar-"+rName),
				),
			},
		},
	})
}

func TestAccAWSDBInstance_optionGroup(t *testing.T) {
	var v rds.DBInstance

	rName := fmt.Sprintf("tf-option-test-%d", acctest.RandInt())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAWSDBInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSDBInstanceConfigWithOptionGroup(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSDBInstanceExists("aws_db_instance.bar", &v),
					testAccCheckAWSDBInstanceAttributes(&v),
					resource.TestCheckResourceAttr(
						"aws_db_instance.bar", "option_group_name", rName),
				),
			},
		},
	})
}

func TestAccAWSDBInstance_iamAuth(t *testing.T) {
	var v rds.DBInstance

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAWSDBInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAWSDBIAMAuth(acctest.RandInt()),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSDBInstanceExists("aws_db_instance.bar", &v),
					testAccCheckAWSDBInstanceAttributes(&v),
					resource.TestCheckResourceAttr(
						"aws_db_instance.bar", "iam_database_authentication_enabled", "true"),
				),
			},
		},
	})
}

func TestAccAWSDBInstance_replica(t *testing.T) {
	var s, r rds.DBInstance

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAWSDBInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccReplicaInstanceConfig(rand.New(rand.NewSource(time.Now().UnixNano())).Int()),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSDBInstanceExists("aws_db_instance.bar", &s),
					testAccCheckAWSDBInstanceExists("aws_db_instance.replica", &r),
					testAccCheckAWSDBInstanceReplicaAttributes(&s, &r),
				),
			},
		},
	})
}

func TestAccAWSDBInstance_noSnapshot(t *testing.T) {
	var snap rds.DBInstance

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAWSDBInstanceNoSnapshot,
		Steps: []resource.TestStep{
			{
				Config: testAccSnapshotInstanceConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSDBInstanceExists("aws_db_instance.snapshot", &snap),
				),
			},
		},
	})
}

func TestAccAWSDBInstance_snapshot(t *testing.T) {
	var snap rds.DBInstance
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		// testAccCheckAWSDBInstanceSnapshot verifies a database snapshot is
		// created, and subequently deletes it
		CheckDestroy: testAccCheckAWSDBInstanceSnapshot(rInt),
		Steps: []resource.TestStep{
			{
				Config: testAccSnapshotInstanceConfigWithSnapshot(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSDBInstanceExists("aws_db_instance.snapshot", &snap),
				),
			},
		},
	})
}

func TestAccAWSDBInstance_enhancedMonitoring(t *testing.T) {
	var dbInstance rds.DBInstance
	rName := acctest.RandString(5)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAWSDBInstanceNoSnapshot,
		Steps: []resource.TestStep{
			{
				Config: testAccSnapshotInstanceConfig_enhancedMonitoring(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSDBInstanceExists("aws_db_instance.enhanced_monitoring", &dbInstance),
					resource.TestCheckResourceAttr(
						"aws_db_instance.enhanced_monitoring", "monitoring_interval", "5"),
				),
			},
		},
	})
}

// Regression test for https://github.com/hashicorp/terraform/issues/3760 .
// We apply a plan, then change just the iops. If the apply succeeds, we
// consider this a pass, as before in 3760 the request would fail
func TestAccAWS_separate_DBInstance_iops_update(t *testing.T) {
	var v rds.DBInstance

	rName := acctest.RandString(5)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAWSDBInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSnapshotInstanceConfig_iopsUpdate(rName, 1000),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSDBInstanceExists("aws_db_instance.bar", &v),
					testAccCheckAWSDBInstanceAttributes(&v),
				),
			},

			{
				Config: testAccSnapshotInstanceConfig_iopsUpdate(rName, 2000),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSDBInstanceExists("aws_db_instance.bar", &v),
					testAccCheckAWSDBInstanceAttributes(&v),
				),
			},
		},
	})
}

func TestAccAWSDBInstance_portUpdate(t *testing.T) {
	var v rds.DBInstance

	rName := acctest.RandString(5)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAWSDBInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSnapshotInstanceConfig_mysqlPort(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSDBInstanceExists("aws_db_instance.bar", &v),
					resource.TestCheckResourceAttr(
						"aws_db_instance.bar", "port", "3306"),
				),
			},

			{
				Config: testAccSnapshotInstanceConfig_updateMysqlPort(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSDBInstanceExists("aws_db_instance.bar", &v),
					resource.TestCheckResourceAttr(
						"aws_db_instance.bar", "port", "3305"),
				),
			},
		},
	})
}

func TestAccAWSDBInstance_MSSQL_TZ(t *testing.T) {
	var v rds.DBInstance
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAWSDBInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSDBMSSQL_timezone(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSDBInstanceExists("aws_db_instance.mssql", &v),
					testAccCheckAWSDBInstanceAttributes_MSSQL(&v, ""),
					resource.TestCheckResourceAttr(
						"aws_db_instance.mssql", "allocated_storage", "20"),
					resource.TestCheckResourceAttr(
						"aws_db_instance.mssql", "engine", "sqlserver-ex"),
				),
			},

			{
				Config: testAccAWSDBMSSQL_timezone_AKST(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSDBInstanceExists("aws_db_instance.mssql", &v),
					testAccCheckAWSDBInstanceAttributes_MSSQL(&v, "Alaskan Standard Time"),
					resource.TestCheckResourceAttr(
						"aws_db_instance.mssql", "allocated_storage", "20"),
					resource.TestCheckResourceAttr(
						"aws_db_instance.mssql", "engine", "sqlserver-ex"),
				),
			},
		},
	})
}

func TestAccAWSDBInstance_MinorVersion(t *testing.T) {
	var v rds.DBInstance

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAWSDBInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSDBInstanceConfigAutoMinorVersion,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSDBInstanceExists("aws_db_instance.bar", &v),
				),
			},
		},
	})
}

// See https://github.com/hashicorp/terraform/issues/11881
func TestAccAWSDBInstance_diffSuppressInitialState(t *testing.T) {
	var v rds.DBInstance
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAWSDBInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSDBInstanceConfigSuppressInitialState(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSDBInstanceExists("aws_db_instance.bar", &v),
				),
			},
		},
	})
}

func TestAccAWSDBInstance_ec2Classic(t *testing.T) {
	var v rds.DBInstance
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t); testAccEC2ClassicPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAWSDBInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSDBInstanceConfigEc2Classic(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSDBInstanceExists("aws_db_instance.bar", &v),
				),
			},
		},
	})
}

func testAccCheckAWSDBInstanceDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*AWSClient).rdsconn

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aws_db_instance" {
			continue
		}

		// Try to find the Group
		var err error
		resp, err := conn.DescribeDBInstances(
			&rds.DescribeDBInstancesInput{
				DBInstanceIdentifier: aws.String(rs.Primary.ID),
			})

		if ae, ok := err.(awserr.Error); ok && ae.Code() == "DBInstanceNotFound" {
			continue
		}

		if err == nil {
			if len(resp.DBInstances) != 0 &&
				*resp.DBInstances[0].DBInstanceIdentifier == rs.Primary.ID {
				return fmt.Errorf("DB Instance still exists")
			}
		}

		// Verify the error
		newerr, ok := err.(awserr.Error)
		if !ok {
			return err
		}
		if newerr.Code() != "DBInstanceNotFound" {
			return err
		}
	}

	return nil
}

func testAccCheckAWSDBInstanceAttributes(v *rds.DBInstance) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if *v.Engine != "mysql" {
			return fmt.Errorf("bad engine: %#v", *v.Engine)
		}

		if *v.EngineVersion == "" {
			return fmt.Errorf("bad engine_version: %#v", *v.EngineVersion)
		}

		if *v.BackupRetentionPeriod != 0 {
			return fmt.Errorf("bad backup_retention_period: %#v", *v.BackupRetentionPeriod)
		}

		return nil
	}
}

func testAccCheckAWSDBInstanceAttributes_MSSQL(v *rds.DBInstance, tz string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if *v.Engine != "sqlserver-ex" {
			return fmt.Errorf("bad engine: %#v", *v.Engine)
		}

		rtz := ""
		if v.Timezone != nil {
			rtz = *v.Timezone
		}

		if tz != rtz {
			return fmt.Errorf("Expected (%s) Timezone for MSSQL test, got (%s)", tz, rtz)
		}

		return nil
	}
}

func testAccCheckAWSDBInstanceReplicaAttributes(source, replica *rds.DBInstance) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if replica.ReadReplicaSourceDBInstanceIdentifier != nil && *replica.ReadReplicaSourceDBInstanceIdentifier != *source.DBInstanceIdentifier {
			return fmt.Errorf("bad source identifier for replica, expected: '%s', got: '%s'", *source.DBInstanceIdentifier, *replica.ReadReplicaSourceDBInstanceIdentifier)
		}

		return nil
	}
}

func testAccCheckAWSDBInstanceSnapshot(rInt int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "aws_db_instance" {
				continue
			}

			awsClient := testAccProvider.Meta().(*AWSClient)
			conn := awsClient.rdsconn

			var err error
			log.Printf("[INFO] Trying to locate the DBInstance Final Snapshot")
			snapshot_identifier := fmt.Sprintf("foobarbaz-test-terraform-final-snapshot-%d", rInt)
			_, snapErr := conn.DescribeDBSnapshots(
				&rds.DescribeDBSnapshotsInput{
					DBSnapshotIdentifier: aws.String(snapshot_identifier),
				})

			if snapErr != nil {
				newerr, _ := snapErr.(awserr.Error)
				if newerr.Code() == "DBSnapshotNotFound" {
					return fmt.Errorf("Snapshot %s not found", snapshot_identifier)
				}
			} else { // snapshot was found,
				// verify we have the tags copied to the snapshot
				instanceARN, err := buildRDSARN(snapshot_identifier, testAccProvider.Meta().(*AWSClient).partition, testAccProvider.Meta().(*AWSClient).accountid, testAccProvider.Meta().(*AWSClient).region)
				// tags have a different ARN, just swapping :db: for :snapshot:
				tagsARN := strings.Replace(instanceARN, ":db:", ":snapshot:", 1)
				if err != nil {
					return fmt.Errorf("Error building ARN for tags check with ARN (%s): %s", tagsARN, err)
				}
				resp, err := conn.ListTagsForResource(&rds.ListTagsForResourceInput{
					ResourceName: aws.String(tagsARN),
				})
				if err != nil {
					return fmt.Errorf("Error retrieving tags for ARN (%s): %s", tagsARN, err)
				}

				if resp.TagList == nil || len(resp.TagList) == 0 {
					return fmt.Errorf("Tag list is nil or zero: %s", resp.TagList)
				}

				var found bool
				for _, t := range resp.TagList {
					if *t.Key == "Name" && *t.Value == "tf-tags-db" {
						found = true
					}
				}
				if !found {
					return fmt.Errorf("Expected to find tag Name (%s), but wasn't found. Tags: %s", "tf-tags-db", resp.TagList)
				}
				// end tag search

				log.Printf("[INFO] Deleting the Snapshot %s", snapshot_identifier)
				_, snapDeleteErr := conn.DeleteDBSnapshot(
					&rds.DeleteDBSnapshotInput{
						DBSnapshotIdentifier: aws.String(snapshot_identifier),
					})
				if snapDeleteErr != nil {
					return err
				}
			} // end snapshot was found

			resp, err := conn.DescribeDBInstances(
				&rds.DescribeDBInstancesInput{
					DBInstanceIdentifier: aws.String(rs.Primary.ID),
				})

			if err != nil {
				newerr, _ := err.(awserr.Error)
				if newerr.Code() != "DBInstanceNotFound" {
					return err
				}

			} else {
				if len(resp.DBInstances) != 0 &&
					*resp.DBInstances[0].DBInstanceIdentifier == rs.Primary.ID {
					return fmt.Errorf("DB Instance still exists")
				}
			}
		}

		return nil
	}
}

func testAccCheckAWSDBInstanceNoSnapshot(s *terraform.State) error {
	conn := testAccProvider.Meta().(*AWSClient).rdsconn

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aws_db_instance" {
			continue
		}

		var err error
		resp, err := conn.DescribeDBInstances(
			&rds.DescribeDBInstancesInput{
				DBInstanceIdentifier: aws.String(rs.Primary.ID),
			})

		if err != nil {
			newerr, _ := err.(awserr.Error)
			if newerr.Code() != "DBInstanceNotFound" {
				return err
			}

		} else {
			if len(resp.DBInstances) != 0 &&
				*resp.DBInstances[0].DBInstanceIdentifier == rs.Primary.ID {
				return fmt.Errorf("DB Instance still exists")
			}
		}

		snapshot_identifier := "foobarbaz-test-terraform-final-snapshot-2"
		_, snapErr := conn.DescribeDBSnapshots(
			&rds.DescribeDBSnapshotsInput{
				DBSnapshotIdentifier: aws.String(snapshot_identifier),
			})

		if snapErr != nil {
			newerr, _ := snapErr.(awserr.Error)
			if newerr.Code() != "DBSnapshotNotFound" {
				return fmt.Errorf("Snapshot %s found and it shouldn't have been", snapshot_identifier)
			}
		}
	}

	return nil
}

func testAccCheckAWSDBInstanceExists(n string, v *rds.DBInstance) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No DB Instance ID is set")
		}

		conn := testAccProvider.Meta().(*AWSClient).rdsconn

		opts := rds.DescribeDBInstancesInput{
			DBInstanceIdentifier: aws.String(rs.Primary.ID),
		}

		resp, err := conn.DescribeDBInstances(&opts)

		if err != nil {
			return err
		}

		if len(resp.DBInstances) != 1 ||
			*resp.DBInstances[0].DBInstanceIdentifier != rs.Primary.ID {
			return fmt.Errorf("DB Instance not found")
		}

		*v = *resp.DBInstances[0]

		return nil
	}
}

// Database names cannot collide, and deletion takes so long, that making the
// name a bit random helps so able we can kill a test that's just waiting for a
// delete and not be blocked on kicking off another one.
var testAccAWSDBInstanceConfig = `
resource "aws_db_instance" "bar" {
	allocated_storage = 10
	engine = "MySQL"
	engine_version = "5.6.35"
	instance_class = "db.t2.micro"
	name = "baz"
	password = "barbarbarbar"
	username = "foo"


	# Maintenance Window is stored in lower case in the API, though not strictly
	# documented. Terraform will downcase this to match (as opposed to throw a
	# validation error).
	maintenance_window = "Fri:09:00-Fri:09:30"
	skip_final_snapshot = true

	backup_retention_period = 0

	parameter_group_name = "default.mysql5.6"

	timeouts {
		create = "30m"
	}
}`

const testAccAWSDBInstanceConfig_namePrefix = `
resource "aws_db_instance" "test" {
	allocated_storage = 10
	engine = "MySQL"
	identifier_prefix = "tf-test-"
	instance_class = "db.t2.micro"
	password = "password"
	username = "root"
	publicly_accessible = true
	skip_final_snapshot = true

	timeouts {
		create = "30m"
	}
}`

const testAccAWSDBInstanceConfig_generatedName = `
resource "aws_db_instance" "test" {
	allocated_storage = 10
	engine = "MySQL"
	instance_class = "db.t2.micro"
	password = "password"
	username = "root"
	publicly_accessible = true
	skip_final_snapshot = true

	timeouts {
		create = "30m"
	}
}`

var testAccAWSDBInstanceConfigKmsKeyId = `
resource "aws_kms_key" "foo" {
    description = "Terraform acc test %s"
    policy = <<POLICY
{
  "Version": "2012-10-17",
  "Id": "kms-tf-1",
  "Statement": [
    {
      "Sid": "Enable IAM User Permissions",
      "Effect": "Allow",
      "Principal": {
        "AWS": "*"
      },
      "Action": "kms:*",
      "Resource": "*"
    }
  ]
}
POLICY
}

resource "aws_db_instance" "bar" {
	allocated_storage = 10
	engine = "MySQL"
	engine_version = "5.6.35"
	instance_class = "db.t2.small"
	name = "baz"
	password = "barbarbarbar"
	username = "foo"


	# Maintenance Window is stored in lower case in the API, though not strictly
	# documented. Terraform will downcase this to match (as opposed to throw a
	# validation error).
	maintenance_window = "Fri:09:00-Fri:09:30"

	backup_retention_period = 0
	storage_encrypted = true
	kms_key_id = "${aws_kms_key.foo.arn}"

	skip_final_snapshot = true

	parameter_group_name = "default.mysql5.6"
}
`

func testAccAWSDBInstanceConfigWithOptionGroup(rName string) string {
	return fmt.Sprintf(`
resource "aws_db_option_group" "bar" {
	name = "%s"
	option_group_description = "Test option group for terraform"
	engine_name = "mysql"
	major_engine_version = "5.6"
}

resource "aws_db_instance" "bar" {
	identifier = "foobarbaz-test-terraform-%d"

	allocated_storage = 10
	engine = "MySQL"
	instance_class = "db.t2.micro"
	name = "baz"
	password = "barbarbarbar"
	username = "foo"

	backup_retention_period = 0
	skip_final_snapshot = true

	parameter_group_name = "default.mysql5.6"
	option_group_name = "${aws_db_option_group.bar.name}"
}`, rName, acctest.RandInt())
}

func testAccCheckAWSDBIAMAuth(n int) string {
	return fmt.Sprintf(`
resource "aws_db_instance" "bar" {
	identifier = "foobarbaz-test-terraform-%d"
	allocated_storage = 10
	engine = "mysql"
	engine_version = "5.6.34"
	instance_class = "db.t2.micro"
	name = "baz"
	password = "barbarbarbar"
	username = "foo"
	backup_retention_period = 0
	skip_final_snapshot = true
	parameter_group_name = "default.mysql5.6"
	iam_database_authentication_enabled = true
}`, n)
}

func testAccReplicaInstanceConfig(val int) string {
	return fmt.Sprintf(`
	resource "aws_db_instance" "bar" {
		identifier = "foobarbaz-test-terraform-%d"

		allocated_storage = 5
		engine = "mysql"
		engine_version = "5.6.35"
		instance_class = "db.t2.micro"
		name = "baz"
		password = "barbarbarbar"
		username = "foo"

		backup_retention_period = 1
		skip_final_snapshot = true

		parameter_group_name = "default.mysql5.6"
	}

	resource "aws_db_instance" "replica" {
		identifier = "tf-replica-db-%d"
		backup_retention_period = 0
		replicate_source_db = "${aws_db_instance.bar.identifier}"
		allocated_storage = "${aws_db_instance.bar.allocated_storage}"
		engine = "${aws_db_instance.bar.engine}"
		engine_version = "${aws_db_instance.bar.engine_version}"
		instance_class = "${aws_db_instance.bar.instance_class}"
		password = "${aws_db_instance.bar.password}"
		username = "${aws_db_instance.bar.username}"
		skip_final_snapshot = true
		tags {
			Name = "tf-replica-db"
		}
	}
	`, val, val)
}

func testAccSnapshotInstanceConfig() string {
	return fmt.Sprintf(`
resource "aws_db_instance" "snapshot" {
	identifier = "tf-test-%d"

	allocated_storage = 5
	engine = "mysql"
	engine_version = "5.6.35"
	instance_class = "db.t2.micro"
	name = "baz"
	password = "barbarbarbar"
	username = "foo"
	backup_retention_period = 1

	publicly_accessible = true

	parameter_group_name = "default.mysql5.6"

	skip_final_snapshot = true
	final_snapshot_identifier = "foobarbaz-test-terraform-final-snapshot-1"
}`, acctest.RandInt())
}

func testAccSnapshotInstanceConfigWithSnapshot(rInt int) string {
	return fmt.Sprintf(`
resource "aws_db_instance" "snapshot" {
	identifier = "tf-snapshot-%d"

	allocated_storage = 5
	engine = "mysql"
	engine_version = "5.6.35"
	instance_class = "db.t2.micro"
	name = "baz"
	password = "barbarbarbar"
	publicly_accessible = true
	username = "foo"
	backup_retention_period = 1

	parameter_group_name = "default.mysql5.6"

	copy_tags_to_snapshot = true
	final_snapshot_identifier = "foobarbaz-test-terraform-final-snapshot-%d"
	tags {
		Name = "tf-tags-db"
	}
}
`, rInt, rInt)
}

func testAccSnapshotInstanceConfig_enhancedMonitoring(rName string) string {
	return fmt.Sprintf(`
resource "aws_iam_role" "enhanced_policy_role" {
    name = "enhanced-monitoring-role-%s"
    assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Sid": "",
      "Effect": "Allow",
      "Principal": {
        "Service": "monitoring.rds.amazonaws.com"
      },
      "Action": "sts:AssumeRole"
    }
  ]
}
EOF

}

resource "aws_iam_policy_attachment" "test-attach" {
    name = "enhanced-monitoring-attachment-%s"
    roles = [
        "${aws_iam_role.enhanced_policy_role.name}",
    ]

    policy_arn = "${aws_iam_policy.test.arn}"
}

resource "aws_iam_policy" "test" {
  name   = "tf-enhanced-monitoring-policy-%s"
  policy = <<POLICY
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Sid": "EnableCreationAndManagementOfRDSCloudwatchLogGroups",
            "Effect": "Allow",
            "Action": [
                "logs:CreateLogGroup",
                "logs:PutRetentionPolicy"
            ],
            "Resource": [
                "arn:aws:logs:*:*:log-group:RDS*"
            ]
        },
        {
            "Sid": "EnableCreationAndManagementOfRDSCloudwatchLogStreams",
            "Effect": "Allow",
            "Action": [
                "logs:CreateLogStream",
                "logs:PutLogEvents",
                "logs:DescribeLogStreams",
                "logs:GetLogEvents"
            ],
            "Resource": [
                "arn:aws:logs:*:*:log-group:RDS*:log-stream:*"
            ]
        }
    ]
}
POLICY
}

resource "aws_db_instance" "enhanced_monitoring" {
	identifier = "foobarbaz-enhanced-monitoring-%s"
	depends_on = ["aws_iam_policy_attachment.test-attach"]

	allocated_storage = 5
	engine = "mysql"
	engine_version = "5.6.35"
	instance_class = "db.t2.micro"
	name = "baz"
	password = "barbarbarbar"
	username = "foo"
	backup_retention_period = 1

	parameter_group_name = "default.mysql5.6"

	monitoring_role_arn = "${aws_iam_role.enhanced_policy_role.arn}"
	monitoring_interval = "5"

	skip_final_snapshot = true
}`, rName, rName, rName, rName)
}

func testAccSnapshotInstanceConfig_iopsUpdate(rName string, iops int) string {
	return fmt.Sprintf(`
resource "aws_db_instance" "bar" {
  identifier           = "mydb-rds-%s"
  engine               = "mysql"
  engine_version       = "5.6.35"
  instance_class       = "db.t2.micro"
  name                 = "mydb"
  username             = "foo"
  password             = "barbarbar"
  parameter_group_name = "default.mysql5.6"
  skip_final_snapshot = true

  apply_immediately = true

  storage_type      = "io1"
  allocated_storage = 200
  iops              = %d
}`, rName, iops)
}

func testAccSnapshotInstanceConfig_mysqlPort(rName string) string {
	return fmt.Sprintf(`
resource "aws_db_instance" "bar" {
  identifier           = "mydb-rds-%s"
  engine               = "mysql"
  engine_version       = "5.6.35"
  instance_class       = "db.t2.micro"
  name                 = "mydb"
  username             = "foo"
  password             = "barbarbar"
  parameter_group_name = "default.mysql5.6"
  port = 3306
  allocated_storage = 10
  skip_final_snapshot = true

  apply_immediately = true
}`, rName)
}

func testAccSnapshotInstanceConfig_updateMysqlPort(rName string) string {
	return fmt.Sprintf(`
resource "aws_db_instance" "bar" {
  identifier           = "mydb-rds-%s"
  engine               = "mysql"
  engine_version       = "5.6.35"
  instance_class       = "db.t2.micro"
  name                 = "mydb"
  username             = "foo"
  password             = "barbarbar"
  parameter_group_name = "default.mysql5.6"
  port = 3305
  allocated_storage = 10
  skip_final_snapshot = true

  apply_immediately = true
}`, rName)
}

func testAccAWSDBInstanceConfigWithSubnetGroup(rName string) string {
	return fmt.Sprintf(`
resource "aws_vpc" "foo" {
	cidr_block = "10.1.0.0/16"
	tags {
		Name="testAccAWSDBInstanceConfigWithSubnetGroup"
	}
}

resource "aws_subnet" "foo" {
	cidr_block = "10.1.1.0/24"
	availability_zone = "us-west-2a"
	vpc_id = "${aws_vpc.foo.id}"
	tags {
		Name = "tf-dbsubnet-test-1"
	}
}

resource "aws_subnet" "bar" {
	cidr_block = "10.1.2.0/24"
	availability_zone = "us-west-2b"
	vpc_id = "${aws_vpc.foo.id}"
	tags {
		Name = "tf-dbsubnet-test-2"
	}
}

resource "aws_db_subnet_group" "foo" {
	name = "foo-%s"
	subnet_ids = ["${aws_subnet.foo.id}", "${aws_subnet.bar.id}"]
	tags {
		Name = "tf-dbsubnet-group-test"
	}
}

resource "aws_db_instance" "bar" {
  identifier           = "mydb-rds-%s"
  engine               = "mysql"
  engine_version       = "5.6.35"
  instance_class       = "db.t2.micro"
  name                 = "mydb"
  username             = "foo"
  password             = "barbarbar"
  parameter_group_name = "default.mysql5.6"
  db_subnet_group_name = "${aws_db_subnet_group.foo.name}"
  port = 3305
  allocated_storage = 10
  skip_final_snapshot = true

	backup_retention_period = 0
  apply_immediately = true
}`, rName, rName)
}

func testAccAWSDBInstanceConfigWithSubnetGroupUpdated(rName string) string {
	return fmt.Sprintf(`
resource "aws_vpc" "foo" {
	cidr_block = "10.1.0.0/16"
	tags {
		Name="testAccAWSDBInstanceConfigWithSubnetGroupUpdated"
	}
}

resource "aws_vpc" "bar" {
	cidr_block = "10.10.0.0/16"
	tags {
		Name="testAccAWSDBInstanceConfigWithSubnetGroupUpdated_other"
	}
}

resource "aws_subnet" "foo" {
	cidr_block = "10.1.1.0/24"
	availability_zone = "us-west-2a"
	vpc_id = "${aws_vpc.foo.id}"
	tags {
		Name = "tf-dbsubnet-test-1"
	}
}

resource "aws_subnet" "bar" {
	cidr_block = "10.1.2.0/24"
	availability_zone = "us-west-2b"
	vpc_id = "${aws_vpc.foo.id}"
	tags {
		Name = "tf-dbsubnet-test-2"
	}
}

resource "aws_subnet" "test" {
	cidr_block = "10.10.3.0/24"
	availability_zone = "us-west-2b"
	vpc_id = "${aws_vpc.bar.id}"
	tags {
		Name = "tf-dbsubnet-test-3"
	}
}

resource "aws_subnet" "another_test" {
	cidr_block = "10.10.4.0/24"
	availability_zone = "us-west-2a"
	vpc_id = "${aws_vpc.bar.id}"
	tags {
		Name = "tf-dbsubnet-test-4"
	}
}

resource "aws_db_subnet_group" "foo" {
	name = "foo-%s"
	subnet_ids = ["${aws_subnet.foo.id}", "${aws_subnet.bar.id}"]
	tags {
		Name = "tf-dbsubnet-group-test"
	}
}

resource "aws_db_subnet_group" "bar" {
	name = "bar-%s"
	subnet_ids = ["${aws_subnet.test.id}", "${aws_subnet.another_test.id}"]
	tags {
		Name = "tf-dbsubnet-group-test-updated"
	}
}

resource "aws_db_instance" "bar" {
  identifier           = "mydb-rds-%s"
  engine               = "mysql"
  engine_version       = "5.6.35"
  instance_class       = "db.t2.micro"
  name                 = "mydb"
  username             = "foo"
  password             = "barbarbar"
  parameter_group_name = "default.mysql5.6"
  db_subnet_group_name = "${aws_db_subnet_group.bar.name}"
  port = 3305
  allocated_storage = 10
  skip_final_snapshot = true

	backup_retention_period = 0

  apply_immediately = true
}`, rName, rName, rName)
}

func testAccAWSDBMSSQL_timezone(rInt int) string {
	return fmt.Sprintf(`
resource "aws_vpc" "foo" {
  cidr_block           = "10.1.0.0/16"
  enable_dns_hostnames = true
	tags {
		Name = "tf-rds-mssql-timezone-test"
	}
}

resource "aws_db_subnet_group" "rds_one" {
  name        = "tf_acc_test_%d"
  description = "db subnets for rds_one"

  subnet_ids = ["${aws_subnet.main.id}", "${aws_subnet.other.id}"]
}

resource "aws_subnet" "main" {
  vpc_id            = "${aws_vpc.foo.id}"
  availability_zone = "us-west-2a"
  cidr_block        = "10.1.1.0/24"
}

resource "aws_subnet" "other" {
  vpc_id            = "${aws_vpc.foo.id}"
  availability_zone = "us-west-2b"
  cidr_block        = "10.1.2.0/24"
}

resource "aws_db_instance" "mssql" {
  identifier = "tf-test-mssql-%d"

  db_subnet_group_name = "${aws_db_subnet_group.rds_one.name}"

  instance_class          = "db.t2.micro"
  allocated_storage       = 20
  username                = "somecrazyusername"
  password                = "somecrazypassword"
  engine                  = "sqlserver-ex"
  backup_retention_period = 0
  skip_final_snapshot = true

  #publicly_accessible = true

  vpc_security_group_ids = ["${aws_security_group.rds-mssql.id}"]
}

resource "aws_security_group" "rds-mssql" {
  name = "tf-rds-mssql-test-%d"

  description = "TF Testing"
  vpc_id      = "${aws_vpc.foo.id}"
}

resource "aws_security_group_rule" "rds-mssql-1" {
  type        = "egress"
  from_port   = 0
  to_port     = 0
  protocol    = "-1"
  cidr_blocks = ["0.0.0.0/0"]

  security_group_id = "${aws_security_group.rds-mssql.id}"
}
`, rInt, rInt, rInt)
}

func testAccAWSDBMSSQL_timezone_AKST(rInt int) string {
	return fmt.Sprintf(`
resource "aws_vpc" "foo" {
  cidr_block           = "10.1.0.0/16"
  enable_dns_hostnames = true
	tags {
		Name = "tf-rds-mssql-timezone-test"
	}
}

resource "aws_db_subnet_group" "rds_one" {
  name        = "tf_acc_test_%d"
  description = "db subnets for rds_one"

  subnet_ids = ["${aws_subnet.main.id}", "${aws_subnet.other.id}"]
}

resource "aws_subnet" "main" {
  vpc_id            = "${aws_vpc.foo.id}"
  availability_zone = "us-west-2a"
  cidr_block        = "10.1.1.0/24"
}

resource "aws_subnet" "other" {
  vpc_id            = "${aws_vpc.foo.id}"
  availability_zone = "us-west-2b"
  cidr_block        = "10.1.2.0/24"
}

resource "aws_db_instance" "mssql" {
  identifier = "tf-test-mssql-%d"

  db_subnet_group_name = "${aws_db_subnet_group.rds_one.name}"

  instance_class          = "db.t2.micro"
  allocated_storage       = 20
  username                = "somecrazyusername"
  password                = "somecrazypassword"
  engine                  = "sqlserver-ex"
  backup_retention_period = 0
  skip_final_snapshot = true

  #publicly_accessible = true

  vpc_security_group_ids = ["${aws_security_group.rds-mssql.id}"]
  timezone               = "Alaskan Standard Time"
}

resource "aws_security_group" "rds-mssql" {
  name = "tf-rds-mssql-test-%d"

  description = "TF Testing"
  vpc_id      = "${aws_vpc.foo.id}"
}

resource "aws_security_group_rule" "rds-mssql-1" {
  type        = "egress"
  from_port   = 0
  to_port     = 0
  protocol    = "-1"
  cidr_blocks = ["0.0.0.0/0"]

  security_group_id = "${aws_security_group.rds-mssql.id}"
}
`, rInt, rInt, rInt)
}

var testAccAWSDBInstanceConfigAutoMinorVersion = fmt.Sprintf(`
resource "aws_db_instance" "bar" {
  identifier = "foobarbaz-test-terraform-%d"
	allocated_storage = 10
	engine = "MySQL"
	engine_version = "5.6"
	instance_class = "db.t2.micro"
	name = "baz"
	password = "barbarbarbar"
	username = "foo"
	skip_final_snapshot = true
}
`, acctest.RandInt())

func testAccAWSDBInstanceConfigEc2Classic(rInt int) string {
	return fmt.Sprintf(`
provider "aws" {
  region = "us-east-1"
}

resource "aws_db_instance" "bar" {
  identifier = "foobarbaz-test-terraform-%d"
  allocated_storage = 10
  engine = "mysql"
  engine_version = "5.6"
  instance_class = "db.m3.medium"
  name = "baz"
  password = "barbarbarbar"
  username = "foo"
  publicly_accessible = true
  security_group_names = ["default"]
  parameter_group_name = "default.mysql5.6"
  skip_final_snapshot = true
}
`, rInt)
}

func testAccAWSDBInstanceConfigSuppressInitialState(rInt int) string {
	return fmt.Sprintf(`
resource "aws_db_instance" "bar" {
  identifier = "foobarbaz-test-terraform-%d"
	allocated_storage = 10
	engine = "MySQL"
	instance_class = "db.t2.micro"
	name = "baz"
	password = "barbarbarbar"
	username = "foo"
	skip_final_snapshot = true
}

data "template_file" "test" {
  template = ""
  vars = {
    test_var = "${aws_db_instance.bar.engine_version}"
  }
}
`, rInt)
}
