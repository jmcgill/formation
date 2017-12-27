package core_test

import (
	"fmt"
	. "github.com/jmcgill/formation/core"
	"strings"

	"github.com/hashicorp/terraform/terraform"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func countLeadingTabs(line string) int {
	i := 0
	for _, runeValue := range line {
		if runeValue == '\t' {
			i++
		} else {
			break
		}
	}
	return i
}

func cleanMultiline(x string) string {
	lines := strings.Split(x, "\n")
	output := make([]string, len(lines)-1)

	// Remove the first line
	lines = lines[1:]

	// Remove the indent from all subsequent lines
	indent := countLeadingTabs(lines[0])
	for i, line := range lines {
		if line == "" {
			output[i] = line
		} else {
			output[i] = line[indent:]
		}
	}
	return strings.Join(output, "\n")
}

var _ = Describe("RoundTripInstanceStateToHCL", func() {
	It("should round trip an empty instance state", func() {
		state := terraform.InstanceState{
			Attributes: map[string]string{},
		}

		expected := cleanMultiline(`
		resource "" "" {
		}`)

		parser := InstanceStateParser{}
		printer := Printer{}
		Expect(printer.Print(parser.Parse(&state))).To(Equal(expected))
	})

	It("should round trip simple fields", func() {
		state := terraform.InstanceState{
			Attributes: map[string]string{
				"simple_field": "value",
			},
		}

		expected := cleanMultiline(`
		resource "" "" {
		    simple_field = "value"
		}`)

		parser := InstanceStateParser{}
		printer := Printer{}
		Expect(printer.Print(parser.Parse(&state))).To(Equal(expected))
	})

	// TODO(jimmy): Should MAP just be a nested resource or a list of length 1?
	It("should round trip a map", func() {
		state := terraform.InstanceState{
			Attributes: map[string]string{
				"map_name.%":         "2",
				"map_name.map_key_1": "map_value_1",
				"map_name.map_key_2": "map_value_2",
			},
		}

		expected := cleanMultiline(`
		resource "" "" {
		    map_name {
		        map_key_1 = "map_value_1"
		        map_key_2 = "map_value_2"
		    }
		}`)

		parser := InstanceStateParser{}
		printer := Printer{}
		Expect(printer.Print(parser.Parse(&state))).To(Equal(expected))
	})

	It("should round trip a list with a nested entry and multiple items", func() {
		state := terraform.InstanceState{
			Attributes: map[string]string{
				"list_key.#":         "2",
				"list_key.1234.name": "Jimmy",
				"list_key.1234.age":  "32",
				"list_key.1235.name": "Alice",
				"list_key.1235.age":  "33",
			},
		}

		expected := cleanMultiline(`
		resource "" "" {
		    list_key {
		        age = "32"
		        name = "Jimmy"
		    }

		    list_key {
		        age = "33"
		        name = "Alice"
		    }
		}`)

		parser := InstanceStateParser{}
		printer := Printer{}
		Expect(printer.Print(parser.Parse(&state))).To(Equal(expected))
	})

	It("should round trip a complicated object", func() {
		state := terraform.InstanceState{
			Attributes: map[string]string{
				"alias.#": "1",
				"alias.701075612.evaluate_target_health": "false",
				"alias.701075612.name":                   "s3-website-us-west-2.amazonaws.com",
				"alias.701075612.zone_id":                "Z3BJ6K6RIION7M",
				"failover":                               "",
				"fqdn":                                   "mikeball.me",
				"health_check_id":                        "",
				"id":                                     "Z2OIQETM3FU6D_mikeball.me_A",
				"name":                                   "mikeball.me",
				"records.#":                              "0",
				"set_identifier":                         "",
				"ttl":                                    "0",
				"type":                                   "A",
				"weight":                                 "-1",
				"zone_id":                                "Z2OIQETM3FU6D",
			},
		}

		expected := cleanMultiline(`
		resource "" "" {
		    alias {
		        evaluate_target_health = "false"
		        name = "s3-website-us-west-2.amazonaws.com"
		        zone_id = "Z3BJ6K6RIION7M"
		    }
		    failover = ""
		    fqdn = "mikeball.me"
		    health_check_id = ""
		    name = "mikeball.me"
		    set_identifier = ""
		    ttl = "0"
		    type = "A"
		    weight = "-1"
		    zone_id = "Z2OIQETM3FU6D"
		}`)

		parser := InstanceStateParser{}
		printer := Printer{}
		Expect(printer.Print(parser.Parse(&state))).To(Equal(expected))
	})

	It("should round trip a nested list", func() {
		state := terraform.InstanceState{
			Attributes: map[string]string{
				"alias.#": "1",
				"alias.701075612.evaluate_target_health": "false",
				"alias.701075612.name":                   "s3-website-us-west-2.amazonaws.com",
				"alias.701075612.nested.#":               "1",
				"alias.701075612.nested.1234.name":       "Jimmy",
				"alias.701075612.nested.1234.age":        "31",
			},
		}

		expected := cleanMultiline(`
		resource "" "" {
		    alias {
		        evaluate_target_health = "false"
		        name = "s3-website-us-west-2.amazonaws.com"
		        nested {
		            age = "31"
		            name = "Jimmy"
		        }
		    }
		}`)

		parser := InstanceStateParser{}
		printer := Printer{}
		Expect(printer.Print(parser.Parse(&state))).To(Equal(expected))
	})

	It("should round trip an S3 bucket", func() {
		// S3 objects have deeply nested lists
		state := terraform.InstanceState{
			Attributes: map[string]string{
				"request_payer": "BucketOwner",
				"arn":           "arn:aws:s3:::formation-test-bucket",
				"versioning.0.mfa_delete":                                                                                   "false",
				"bucket":                                                                                                    "formation-test-bucket",
				"bucket_domain_name":                                                                                        "formation-test-bucket.s3.amazonaws.com",
				"server_side_encryption_configuration.#":                                                                    "2",
				"server_side_encryption_configuration.0.rule.#":                                                             "1",
				"server_side_encryption_configuration.0.rule.0.apply_server_side_encryption_by_default.#":                   "1",
				"server_side_encryption_configuration.0.rule.0.apply_server_side_encryption_by_default.0.sse_algorithm":     "AES256",
				"server_side_encryption_configuration.0.rule.0.apply_server_side_encryption_by_default.0.kms_master_key_id": "",
				"server_side_encryption_configuration.0.item":                                                               "One",
				"server_side_encryption_configuration.1.item":                                                               "Two",
				"region":               "us-east-2",
				"versioning.0.enabled": "true",
				"versioning.#":         "1",
				"tags.environment":     "production",
				"tags.%":               "1",
				"logging.#":            "0",
				"website.#":            "0",
				"acceleration_status":  "",
				"id":             "formation-test-bucket",
				"hosted_zone_id": "Z2O1EMRO9K5GLX",
			},
		}

		expected := cleanMultiline(`
		resource "" "" {
		    acceleration_status = ""
		    arn = "arn:aws:s3:::formation-test-bucket"
		    bucket = "formation-test-bucket"
		    bucket_domain_name = "formation-test-bucket.s3.amazonaws.com"
		    hosted_zone_id = "Z2O1EMRO9K5GLX"
		    region = "us-east-2"
		    request_payer = "BucketOwner"
		    server_side_encryption_configuration {
		        item = "One"
		        rule {
		            apply_server_side_encryption_by_default {
		                kms_master_key_id = ""
		                sse_algorithm = "AES256"
		            }
		        }
		    }

		    server_side_encryption_configuration {
		        item = "Two"
		    }
		    tags {
		        environment = "production"
		    }
		    versioning {
		        enabled = "true"
		        mfa_delete = "false"
		    }
		}`)

		parser := InstanceStateParser{}
		printer := Printer{}

		Expect(printer.Print(parser.Parse(&state))).To(Equal(expected))
	})

	It("BUG1: Parses a resource with 0 entries in a map", func() {
		// S3 objects have deeply nested lists
		state := terraform.InstanceState{
			Attributes: map[string]string{
				"bucket_domain_name":                     "formation-test-bucket-two.s3.amazonaws.com",
				"arn":                                    "arn:aws:s3:::formation-test-bucket-two",
				"website.#":                              "0",
				"logging.#":                              "1",
				"region":                                 "us-east-2",
				"hosted_zone_id":                         "Z2O1EMRO9K5GLX",
				"logging.1011077041.target_prefix":       "my-prefix",
				"versioning.0.enabled":                   "false",
				"request_payer":                          "BucketOwner",
				"acceleration_status":                    "",
				"server_side_encryption_configuration.#": "0",
				"versioning.#":                           "1",
				"bucket":                                 "formation-test-bucket-two",
				"tags.%":                                 "0",
				"logging.1011077041.target_bucket": "formation-test-bucket",
				"versioning.0.mfa_delete":          "false",
				"id": "formation-test-bucket-two",
			},
		}

		expected := cleanMultiline(`
		resource "" "" {
		    acceleration_status = ""
		    arn = "arn:aws:s3:::formation-test-bucket-two"
		    bucket = "formation-test-bucket-two"
		    bucket_domain_name = "formation-test-bucket-two.s3.amazonaws.com"
		    hosted_zone_id = "Z2O1EMRO9K5GLX"
		    logging {
		        target_bucket = "formation-test-bucket"
		        target_prefix = "my-prefix"
		    }
		    region = "us-east-2"
		    request_payer = "BucketOwner"
		    versioning {
		        enabled = "false"
		        mfa_delete = "false"
		    }
		}`)

		parser := InstanceStateParser{}
		printer := Printer{}

		fmt.Println(printer.Print(parser.Parse(&state)))
		Expect(printer.Print(parser.Parse(&state))).To(Equal(expected))
	})
})
