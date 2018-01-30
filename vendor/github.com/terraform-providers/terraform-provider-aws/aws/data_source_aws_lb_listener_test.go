package aws

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDataSourceAWSLBListener_basic(t *testing.T) {
	lbName := fmt.Sprintf("testlistener-basic-%s", acctest.RandStringFromCharSet(13, acctest.CharSetAlphaNum))
	targetGroupName := fmt.Sprintf("testtargetgroup-%s", acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAWSLBListenerConfigBasic(lbName, targetGroupName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.aws_lb_listener.front_end", "load_balancer_arn"),
					resource.TestCheckResourceAttrSet("data.aws_lb_listener.front_end", "arn"),
					resource.TestCheckResourceAttr("data.aws_lb_listener.front_end", "protocol", "HTTP"),
					resource.TestCheckResourceAttr("data.aws_lb_listener.front_end", "port", "80"),
					resource.TestCheckResourceAttr("data.aws_lb_listener.front_end", "default_action.#", "1"),
					resource.TestCheckResourceAttr("data.aws_lb_listener.front_end", "default_action.0.type", "forward"),
					resource.TestCheckResourceAttrSet("data.aws_lb_listener.front_end", "default_action.0.target_group_arn"),
				),
			},
		},
	})
}

func TestAccDataSourceAWSLBListenerBackwardsCompatibility(t *testing.T) {
	lbName := fmt.Sprintf("testlistener-basic-%s", acctest.RandStringFromCharSet(13, acctest.CharSetAlphaNum))
	targetGroupName := fmt.Sprintf("testtargetgroup-%s", acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAWSLBListenerConfigBackwardsCompatibility(lbName, targetGroupName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.aws_alb_listener.front_end", "load_balancer_arn"),
					resource.TestCheckResourceAttrSet("data.aws_alb_listener.front_end", "arn"),
					resource.TestCheckResourceAttr("data.aws_alb_listener.front_end", "protocol", "HTTP"),
					resource.TestCheckResourceAttr("data.aws_alb_listener.front_end", "port", "80"),
					resource.TestCheckResourceAttr("data.aws_alb_listener.front_end", "default_action.#", "1"),
					resource.TestCheckResourceAttr("data.aws_alb_listener.front_end", "default_action.0.type", "forward"),
					resource.TestCheckResourceAttrSet("data.aws_alb_listener.front_end", "default_action.0.target_group_arn"),
				),
			},
		},
	})
}

func TestAccDataSourceAWSLBListener_https(t *testing.T) {
	lbName := fmt.Sprintf("testlistener-https-%s", acctest.RandStringFromCharSet(13, acctest.CharSetAlphaNum))
	targetGroupName := fmt.Sprintf("testtargetgroup-%s", acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProvidersWithTLS,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAWSLBListenerConfigHTTPS(lbName, targetGroupName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.aws_lb_listener.front_end", "load_balancer_arn"),
					resource.TestCheckResourceAttrSet("data.aws_lb_listener.front_end", "arn"),
					resource.TestCheckResourceAttr("data.aws_lb_listener.front_end", "protocol", "HTTPS"),
					resource.TestCheckResourceAttr("data.aws_lb_listener.front_end", "port", "443"),
					resource.TestCheckResourceAttr("data.aws_lb_listener.front_end", "default_action.#", "1"),
					resource.TestCheckResourceAttr("data.aws_lb_listener.front_end", "default_action.0.type", "forward"),
					resource.TestCheckResourceAttrSet("data.aws_lb_listener.front_end", "default_action.0.target_group_arn"),
					resource.TestCheckResourceAttrSet("data.aws_lb_listener.front_end", "certificate_arn"),
					resource.TestCheckResourceAttr("data.aws_lb_listener.front_end", "ssl_policy", "ELBSecurityPolicy-2015-05"),
				),
			},
		},
	})
}

func testAccDataSourceAWSLBListenerConfigBasic(lbName, targetGroupName string) string {
	return fmt.Sprintf(`resource "aws_lb_listener" "front_end" {
   load_balancer_arn = "${aws_lb.alb_test.id}"
   protocol = "HTTP"
   port = "80"

   default_action {
     target_group_arn = "${aws_lb_target_group.test.id}"
     type = "forward"
   }
}

resource "aws_lb" "alb_test" {
  name            = "%s"
  internal        = true
  security_groups = ["${aws_security_group.alb_test.id}"]
  subnets         = ["${aws_subnet.alb_test.*.id}"]

  idle_timeout = 30
  enable_deletion_protection = false

  tags {
    TestName = "TestAccAWSALB_basic"
  }
}

resource "aws_lb_target_group" "test" {
  name = "%s"
  port = 8080
  protocol = "HTTP"
  vpc_id = "${aws_vpc.alb_test.id}"

  health_check {
    path = "/health"
    interval = 60
    port = 8081
    protocol = "HTTP"
    timeout = 3
    healthy_threshold = 3
    unhealthy_threshold = 3
    matcher = "200-299"
  }
}

variable "subnets" {
  default = ["10.0.1.0/24", "10.0.2.0/24"]
  type    = "list"
}

data "aws_availability_zones" "available" {}

resource "aws_vpc" "alb_test" {
  cidr_block = "10.0.0.0/16"

  tags {
    TestName = "TestAccAWSALB_basic"
  }
}

resource "aws_subnet" "alb_test" {
  count                   = 2
  vpc_id                  = "${aws_vpc.alb_test.id}"
  cidr_block              = "${element(var.subnets, count.index)}"
  map_public_ip_on_launch = true
  availability_zone       = "${element(data.aws_availability_zones.available.names, count.index)}"

  tags {
    TestName = "TestAccAWSALB_basic"
  }
}

resource "aws_security_group" "alb_test" {
  name        = "allow_all_alb_test"
  description = "Used for ALB Testing"
  vpc_id      = "${aws_vpc.alb_test.id}"

  ingress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags {
    TestName = "TestAccAWSALB_basic"
  }
}

data "aws_lb_listener" "front_end" {
	arn = "${aws_lb_listener.front_end.arn}"
}`, lbName, targetGroupName)
}

func testAccDataSourceAWSLBListenerConfigBackwardsCompatibility(lbName, targetGroupName string) string {
	return fmt.Sprintf(`resource "aws_alb_listener" "front_end" {
   load_balancer_arn = "${aws_alb.alb_test.id}"
   protocol = "HTTP"
   port = "80"

   default_action {
     target_group_arn = "${aws_alb_target_group.test.id}"
     type = "forward"
   }
}

resource "aws_alb" "alb_test" {
  name            = "%s"
  internal        = true
  security_groups = ["${aws_security_group.alb_test.id}"]
  subnets         = ["${aws_subnet.alb_test.*.id}"]

  idle_timeout = 30
  enable_deletion_protection = false

  tags {
    TestName = "TestAccAWSALB_basic"
  }
}

resource "aws_alb_target_group" "test" {
  name = "%s"
  port = 8080
  protocol = "HTTP"
  vpc_id = "${aws_vpc.alb_test.id}"

  health_check {
    path = "/health"
    interval = 60
    port = 8081
    protocol = "HTTP"
    timeout = 3
    healthy_threshold = 3
    unhealthy_threshold = 3
    matcher = "200-299"
  }
}

variable "subnets" {
  default = ["10.0.1.0/24", "10.0.2.0/24"]
  type    = "list"
}

data "aws_availability_zones" "available" {}

resource "aws_vpc" "alb_test" {
  cidr_block = "10.0.0.0/16"

  tags {
    TestName = "TestAccAWSALB_basic"
  }
}

resource "aws_subnet" "alb_test" {
  count                   = 2
  vpc_id                  = "${aws_vpc.alb_test.id}"
  cidr_block              = "${element(var.subnets, count.index)}"
  map_public_ip_on_launch = true
  availability_zone       = "${element(data.aws_availability_zones.available.names, count.index)}"

  tags {
    TestName = "TestAccAWSALB_basic"
  }
}

resource "aws_security_group" "alb_test" {
  name        = "allow_all_alb_test"
  description = "Used for ALB Testing"
  vpc_id      = "${aws_vpc.alb_test.id}"

  ingress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags {
    TestName = "TestAccAWSALB_basic"
  }
}

data "aws_alb_listener" "front_end" {
	arn = "${aws_alb_listener.front_end.arn}"
}`, lbName, targetGroupName)
}

func testAccDataSourceAWSLBListenerConfigHTTPS(lbName, targetGroupName string) string {
	return fmt.Sprintf(`resource "aws_lb_listener" "front_end" {
   load_balancer_arn = "${aws_lb.alb_test.id}"
   protocol = "HTTPS"
   port = "443"
   ssl_policy = "ELBSecurityPolicy-2015-05"
   certificate_arn = "${aws_iam_server_certificate.test_cert.arn}"

   default_action {
     target_group_arn = "${aws_lb_target_group.test.id}"
     type = "forward"
   }
}

resource "aws_lb" "alb_test" {
  name            = "%s"
  internal        = false
  security_groups = ["${aws_security_group.alb_test.id}"]
  subnets         = ["${aws_subnet.alb_test.*.id}"]

  idle_timeout = 30
  enable_deletion_protection = false

  tags {
    TestName = "TestAccAWSALB_basic"
  }
}

resource "aws_lb_target_group" "test" {
  name = "%s"
  port = 8080
  protocol = "HTTP"
  vpc_id = "${aws_vpc.alb_test.id}"

  health_check {
    path = "/health"
    interval = 60
    port = 8081
    protocol = "HTTP"
    timeout = 3
    healthy_threshold = 3
    unhealthy_threshold = 3
    matcher = "200-299"
  }
}

variable "subnets" {
  default = ["10.0.1.0/24", "10.0.2.0/24"]
  type    = "list"
}

data "aws_availability_zones" "available" {}

resource "aws_vpc" "alb_test" {
  cidr_block = "10.0.0.0/16"

  tags {
    TestName = "TestAccAWSALB_basic"
  }
}

resource "aws_internet_gateway" "gw" {
    vpc_id = "${aws_vpc.alb_test.id}"

    tags {
        TestName = "TestAccAWSALB_basic"
    }
}

resource "aws_subnet" "alb_test" {
  count                   = 2
  vpc_id                  = "${aws_vpc.alb_test.id}"
  cidr_block              = "${element(var.subnets, count.index)}"
  map_public_ip_on_launch = true
  availability_zone       = "${element(data.aws_availability_zones.available.names, count.index)}"

  tags {
    TestName = "TestAccAWSALB_basic"
  }
}

resource "aws_security_group" "alb_test" {
  name        = "allow_all_alb_test"
  description = "Used for ALB Testing"
  vpc_id      = "${aws_vpc.alb_test.id}"

  ingress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags {
    TestName = "TestAccAWSALB_basic"
  }
}

resource "aws_iam_server_certificate" "test_cert" {
  name = "terraform-test-cert-%d"
  certificate_body = "${tls_self_signed_cert.example.cert_pem}"
  private_key      = "${tls_private_key.example.private_key_pem}"
}

resource "tls_private_key" "example" {
  algorithm = "RSA"
}

resource "tls_self_signed_cert" "example" {
  key_algorithm   = "RSA"
  private_key_pem = "${tls_private_key.example.private_key_pem}"

  subject {
    common_name  = "example.com"
    organization = "ACME Examples, Inc"
  }

  validity_period_hours = 12

  allowed_uses = [
    "key_encipherment",
    "digital_signature",
    "server_auth",
  ]
}

data "aws_lb_listener" "front_end" {
	arn = "${aws_lb_listener.front_end.arn}"
}`, lbName, targetGroupName, rand.New(rand.NewSource(time.Now().UnixNano())).Int())
}
