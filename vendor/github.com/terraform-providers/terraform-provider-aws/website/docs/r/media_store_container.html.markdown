---
layout: "aws"
page_title: "AWS: aws_media_store_container"
sidebar_current: "docs-aws-resource-media-store-container"
description: |-
  Provides a MediaStore Container.
---

# aws_media_store_container

Provides a MediaStore Container.

## Example Usage

```hcl
resource "aws_media_store_container" "example" {
  name = "example"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the container. Must contain alphanumeric characters or underscores.

## Attributes Reference

The following attributes are exported:

* `arn` - The ARN of the container.
* `endpoint` - The DNS endpoint of the container.
