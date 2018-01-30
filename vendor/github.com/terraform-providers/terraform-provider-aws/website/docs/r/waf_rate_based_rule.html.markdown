---
layout: "aws"
page_title: "AWS: waf_rate_based_rule"
sidebar_current: "docs-aws-resource-waf-rate-based-rule"
description: |-
  Provides a AWS WAF rule resource.
---

# aws_waf_rate_based_rule

Provides a WAF Rate Based Rule Resource

## Example Usage

```hcl
resource "aws_waf_ipset" "ipset" {
  name = "tfIPSet"

  ip_set_descriptors {
    type  = "IPV4"
    value = "192.0.7.0/24"
  }
}

resource "aws_waf_rate_based_rule" "wafrule" {
  depends_on  = ["aws_waf_ipset.ipset"]
  name        = "tfWAFRule"
  metric_name = "tfWAFRule"

  rate_key = "IP"
  rate_limit = 2000

  predicates {
    data_id = "${aws_waf_ipset.ipset.id}"
    negated = false
    type    = "IPMatch"
  }
}
```

## Argument Reference

The following arguments are supported:

* `metric_name` - (Required) The name or description for the Amazon CloudWatch metric of this rule.
* `name` - (Required) The name or description of the rule.
* `rate_key` - (Required) Valid value is IP.
* `rate_limit` - (Required) The maximum number of requests, which have an identical value in the field specified by the RateKey, allowed in a five-minute period. Minimum value is 2000.
* `predicates` - (Optional) One of ByteMatchSet, IPSet, SizeConstraintSet, SqlInjectionMatchSet, or XssMatchSet objects to include in a rule.

## Nested Blocks

### `predicates`

#### Arguments

* `negated` - (Required) Set this to `false` if you want to allow, block, or count requests
  based on the settings in the specified `ByteMatchSet`, `IPSet`, `SqlInjectionMatchSet`, `XssMatchSet`, or `SizeConstraintSet`.
  For example, if an IPSet includes the IP address `192.0.2.44`, AWS WAF will allow or block requests based on that IP address.
  If set to `true`, AWS WAF will allow, block, or count requests based on all IP addresses _except_ `192.0.2.44`.
* `data_id` - (Required) A unique identifier for a predicate in the rule, such as Byte Match Set ID or IPSet ID.
* `type` - (Required) The type of predicate in a rule, such as `ByteMatchSet` or `IPSet`

## Remarks

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the WAF rule.
