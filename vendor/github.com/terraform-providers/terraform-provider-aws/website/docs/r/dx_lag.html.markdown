---
layout: "aws"
page_title: "AWS: aws_dx_lag"
sidebar_current: "docs-aws-resource-dx-lag"
description: |-  
  Provides a Direct Connect LAG.
---

# aws_dx_lag

Provides a Direct Connect LAG.

## Example Usage

```hcl
resource "aws_dx_lag" "hoge" {
  name = "tf-dx-lag"
  connections_bandwidth = "1Gbps"
  location = "EqDC2"
  number_of_connections = 2
  force_destroy = true
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the LAG.
* `connections_bandwidth` - (Required) The bandwidth of the individual physical connections bundled by the LAG. Available values: 1Gbps, 10Gbps. Case sensitive.
* `location` - (Required) The AWS Direct Connect location in which the LAG should be allocated. See [DescribeLocations](https://docs.aws.amazon.com/directconnect/latest/APIReference/API_DescribeLocations.html) for the list of AWS Direct Connect locations. Use `locationCode`.
* `number_of_connections` - (Required) The number of physical connections initially provisioned and bundled by the LAG.
* `force_destroy` - (Optional, Default:false) A boolean that indicates all connections associated with the LAG should be deleted so that the LAG can be destroyed without error. These objects are *not* recoverable.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the LAG.
