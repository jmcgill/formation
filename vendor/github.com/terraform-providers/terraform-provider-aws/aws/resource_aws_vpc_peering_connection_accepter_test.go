package aws

import (
	"testing"

	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAWSVPCPeeringConnectionAccepter_sameRegion(t *testing.T) {
	var connection ec2.VpcPeeringConnection

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccAwsVPCPeeringConnectionAccepterDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccAwsVPCPeeringConnectionAccepterSameRegion,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSVpcPeeringConnectionExists(
						"aws_vpc_peering_connection_accepter.peer",
						&connection),
					resource.TestCheckResourceAttr(
						"aws_vpc_peering_connection_accepter.peer",
						"accept_status", "active"),
				),
			},
		},
	})
}

func TestAccAWSVPCPeeringConnectionAccepter_differentRegion(t *testing.T) {
	var connection ec2.VpcPeeringConnection

	var providers []*schema.Provider
	providerFactories := map[string]terraform.ResourceProviderFactory{
		"aws": func() (terraform.ResourceProvider, error) {
			p := Provider()
			providers = append(providers, p.(*schema.Provider))
			return p, nil
		},
	}

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccAwsVPCPeeringConnectionAccepterDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccAwsVPCPeeringConnectionAccepterDifferentRegion,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSVpcPeeringConnectionExists(
						"aws_vpc_peering_connection_accepter.peer",
						&connection),
					resource.TestCheckResourceAttr(
						"aws_vpc_peering_connection_accepter.peer",
						"accept_status", "active"),
				),
			},
		},
	})
}

func testAccAwsVPCPeeringConnectionAccepterDestroy(s *terraform.State) error {
	// We don't destroy the underlying VPC Peering Connection.
	return nil
}

const testAccAwsVPCPeeringConnectionAccepterSameRegion = `
resource "aws_vpc" "main" {
	cidr_block = "10.0.0.0/16"
	tags {
		Name = "tf-acc-revoke-vpc-peering-connection-accepter-same-region"
	}
}

resource "aws_vpc" "peer" {
	cidr_block = "10.1.0.0/16"
	tags {
		Name = "tf-acc-revoke-vpc-peering-connection-accepter-same-region"
	}
}

// Requester's side of the connection.
resource "aws_vpc_peering_connection" "peer" {
	vpc_id = "${aws_vpc.main.id}"
	peer_vpc_id = "${aws_vpc.peer.id}"
	auto_accept = false
}

// Accepter's side of the connection.
resource "aws_vpc_peering_connection_accepter" "peer" {
	vpc_peering_connection_id = "${aws_vpc_peering_connection.peer.id}"
	auto_accept = true
}
`

const testAccAwsVPCPeeringConnectionAccepterDifferentRegion = `
provider "aws" {
	alias = "main"
	region = "us-west-2"
}

provider "aws" {
	alias = "peer"
	region = "us-east-1"
}

resource "aws_vpc" "main" {
	provider = "aws.main"
	cidr_block = "10.0.0.0/16"
	tags {
		Name = "tf-acc-revoke-vpc-peering-connection-accepter-different-region"
	}
}

resource "aws_vpc" "peer" {
	provider = "aws.peer"
	cidr_block = "10.1.0.0/16"
	tags {
		Name = "tf-acc-revoke-vpc-peering-connection-accepter-different-region"
	}
}

// Requester's side of the connection.
resource "aws_vpc_peering_connection" "peer" {
	provider = "aws.main"
	vpc_id = "${aws_vpc.main.id}"
	peer_vpc_id = "${aws_vpc.peer.id}"
	peer_region = "us-east-1"
	auto_accept = false
}

// Accepter's side of the connection.
resource "aws_vpc_peering_connection_accepter" "peer" {
	provider = "aws.peer"
	vpc_peering_connection_id = "${aws_vpc_peering_connection.peer.id}"
	auto_accept = true
}
`
