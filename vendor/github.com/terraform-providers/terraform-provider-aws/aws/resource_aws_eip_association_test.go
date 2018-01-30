package aws

import (
	"fmt"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAWSEIPAssociation_basic(t *testing.T) {
	var a ec2.Address

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAWSEIPAssociationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSEIPAssociationConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSEIPExists(
						"aws_eip.bar.0", &a),
					testAccCheckAWSEIPAssociationExists(
						"aws_eip_association.by_allocation_id", &a),
					testAccCheckAWSEIPExists(
						"aws_eip.bar.1", &a),
					testAccCheckAWSEIPAssociationExists(
						"aws_eip_association.by_public_ip", &a),
					testAccCheckAWSEIPExists(
						"aws_eip.bar.2", &a),
					testAccCheckAWSEIPAssociationExists(
						"aws_eip_association.to_eni", &a),
				),
			},
		},
	})
}

func TestAccAWSEIPAssociation_ec2Classic(t *testing.T) {
	var a ec2.Address

	oldvar := os.Getenv("AWS_DEFAULT_REGION")
	os.Setenv("AWS_DEFAULT_REGION", "us-east-1")
	defer os.Setenv("AWS_DEFAULT_REGION", oldvar)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t); testAccEC2ClassicPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAWSEIPAssociationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSEIPAssociationConfig_ec2Classic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSEIPExists("aws_eip.test", &a),
					testAccCheckAWSEIPAssociationExists("aws_eip_association.test", &a),
					resource.TestCheckResourceAttrSet("aws_eip_association.test", "public_ip"),
					resource.TestCheckResourceAttr("aws_eip_association.test", "allocation_id", ""),
					testAccCheckAWSEIPAssociationHasIpBasedId("aws_eip_association.test", &a),
				),
			},
		},
	})
}

func TestAccAWSEIPAssociation_disappears(t *testing.T) {
	var a ec2.Address

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAWSEIPAssociationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSEIPAssociationConfigDisappears,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSEIPExists(
						"aws_eip.bar", &a),
					testAccCheckAWSEIPAssociationExists(
						"aws_eip_association.by_allocation_id", &a),
					testAccCheckEIPAssociationDisappears(&a),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccCheckEIPAssociationDisappears(address *ec2.Address) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := testAccProvider.Meta().(*AWSClient).ec2conn
		opts := &ec2.DisassociateAddressInput{
			AssociationId: address.AssociationId,
		}
		if _, err := conn.DisassociateAddress(opts); err != nil {
			return err
		}
		return nil
	}
}

func testAccCheckAWSEIPAssociationExists(name string, res *ec2.Address) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No EIP Association ID is set")
		}

		conn := testAccProvider.Meta().(*AWSClient).ec2conn
		platforms := testAccProvider.Meta().(*AWSClient).supportedplatforms

		request, err := describeAddressesById(rs.Primary.ID, platforms)
		if err != nil {
			return err
		}

		describe, err := conn.DescribeAddresses(request)
		if err != nil {
			return err
		}

		if len(describe.Addresses) != 1 ||
			(!hasEc2Classic(platforms) && *describe.Addresses[0].AssociationId != *res.AssociationId) {
			return fmt.Errorf("EIP Association not found")
		}

		return nil
	}
}

func testAccCheckAWSEIPAssociationHasIpBasedId(name string, res *ec2.Address) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No EIP Association ID is set")
		}

		if rs.Primary.ID != rs.Primary.Attributes["public_ip"] {
			return fmt.Errorf("Expected EIP Association ID to be equal to Public IP (%q), given: %q",
				rs.Primary.Attributes["public_ip"], rs.Primary.ID)
		}
		return nil
	}
}

func testAccCheckAWSEIPAssociationDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aws_eip_association" {
			continue
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No EIP Association ID is set")
		}

		conn := testAccProvider.Meta().(*AWSClient).ec2conn

		request := &ec2.DescribeAddressesInput{
			Filters: []*ec2.Filter{
				&ec2.Filter{
					Name:   aws.String("association-id"),
					Values: []*string{aws.String(rs.Primary.ID)},
				},
			},
		}
		describe, err := conn.DescribeAddresses(request)
		if err != nil {
			return err
		}

		if len(describe.Addresses) > 0 {
			return fmt.Errorf("EIP Association still exists")
		}
	}
	return nil
}

const testAccAWSEIPAssociationConfig = `
resource "aws_vpc" "main" {
	cidr_block = "192.168.0.0/24"
	tags {
		Name = "testAccAWSEIPAssociationConfig"
	}
}
resource "aws_subnet" "sub" {
	vpc_id = "${aws_vpc.main.id}"
	cidr_block = "192.168.0.0/25"
	availability_zone = "us-west-2a"
}
resource "aws_internet_gateway" "igw" {
	vpc_id = "${aws_vpc.main.id}"
}
resource "aws_instance" "foo" {
	count = 2
	ami = "ami-21f78e11"
	availability_zone = "us-west-2a"
	instance_type = "t1.micro"
	subnet_id = "${aws_subnet.sub.id}"
	private_ip = "192.168.0.${count.index+10}"
}
resource "aws_eip" "bar" {
	count = 3
	vpc = true
}
resource "aws_eip_association" "by_allocation_id" {
	allocation_id = "${aws_eip.bar.0.id}"
	instance_id = "${aws_instance.foo.0.id}"
	depends_on = ["aws_instance.foo"]
}
resource "aws_eip_association" "by_public_ip" {
	public_ip = "${aws_eip.bar.1.public_ip}"
	instance_id = "${aws_instance.foo.1.id}"
	depends_on = ["aws_instance.foo"]
}
resource "aws_eip_association" "to_eni" {
	allocation_id = "${aws_eip.bar.2.id}"
	network_interface_id = "${aws_network_interface.baz.id}"
}
resource "aws_network_interface" "baz" {
	subnet_id = "${aws_subnet.sub.id}"
	private_ips = ["192.168.0.50"]
	depends_on = ["aws_instance.foo"]
	attachment {
		instance = "${aws_instance.foo.0.id}"
		device_index = 1
	}
}
`

const testAccAWSEIPAssociationConfigDisappears = `
resource "aws_vpc" "main" {
	cidr_block = "192.168.0.0/24"
	tags {
		Name = "testAccAWSEIPAssociationConfigDisappears"
	}
}
resource "aws_subnet" "sub" {
	vpc_id = "${aws_vpc.main.id}"
	cidr_block = "192.168.0.0/25"
	availability_zone = "us-west-2a"
}
resource "aws_internet_gateway" "igw" {
	vpc_id = "${aws_vpc.main.id}"
}
resource "aws_instance" "foo" {
	ami = "ami-21f78e11"
	availability_zone = "us-west-2a"
	instance_type = "t1.micro"
	subnet_id = "${aws_subnet.sub.id}"
}
resource "aws_eip" "bar" {
	vpc = true
}
resource "aws_eip_association" "by_allocation_id" {
	allocation_id = "${aws_eip.bar.id}"
	instance_id = "${aws_instance.foo.id}"
}`

const testAccAWSEIPAssociationConfig_ec2Classic = `
provider "aws" {
  region = "us-east-1"
}

resource "aws_eip" "test" {}

data "aws_availability_zones" "available" {}

data "aws_ami" "ubuntu" {
  most_recent = true

  filter {
    name   = "name"
    values = ["ubuntu/images/ebs/ubuntu-trusty-14.04-i386-server-*"]
  }

  filter {
    name   = "virtualization-type"
    values = ["paravirtual"]
  }

  owners = ["099720109477"] # Canonical
}

resource "aws_instance" "test" {
  ami = "${data.aws_ami.ubuntu.id}"
  availability_zone = "${data.aws_availability_zones.available.names[0]}"
  instance_type = "t1.micro"
}

resource "aws_eip_association" "test" {
  public_ip = "${aws_eip.test.public_ip}"
  instance_id = "${aws_instance.test.id}"
}
`
