---
layout: "aws"
page_title: "AWS: batch"
sidebar_current: "docs-aws-resource-batch-job-definition"
description: |-
  Provides a Batch Job Definition resource.
---

# aws_batch_job_definition

Provides a Batch Job Definition resource.

## Example Usage

```hcl
resource "aws_batch_job_definition" "test" {
	name = "tf_test_batch_job_definition"
	type = "container"
	container_properties = <<CONTAINER_PROPERTIES
{
	"command": ["ls", "-la"],
	"image": "busybox",
	"memory": 1024,
	"vcpus": 1,
	"volumes": [
      {
        "host": {
          "sourcePath": "/tmp"
        },
        "name": "tmp"
      }
    ],
	"environment": [
		{"name": "VARNAME", "value": "VARVAL"}
	],
	"mountPoints": [
		{
          "sourceVolume": "tmp",
          "containerPath": "/tmp",
          "readOnly": false
        }
	],
    "ulimits": [
      {
        "hardLimit": 1024,
        "name": "nofile",
        "softLimit": 1024
      }
    ]
}
CONTAINER_PROPERTIES
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the job definition.
* `container_properties` - (Optional) A valid [container properties](http://docs.aws.amazon.com/batch/latest/APIReference/API_RegisterJobDefinition.html)
    provided as a single valid JSON document. This parameter is required if the `type` parameter is `container`.
* `parameters` - (Optional) Specifies the parameter substitution placeholders to set in the job definition.
* `retry_strategy` - (Optional) Specifies the retry strategy to use for failed jobs that are submitted with this job definition.
    Maximum number of `retry_strategy` is `1`.  Defined below.
* `type` - (Required) The type of job definition.  Must be `container`

## retry_strategy

`retry_strategy` supports the following:

* `attempts` - (Required) The number of times to move a job to the `RUNNABLE` status. You may specify between `1` and `10` attempts.

## Attribute Reference

The following attributes are exported:

* `arn` - The Amazon Resource Name of the job definition.
* `revision` - The revision of the job definition.
