---
layout: "aws"
page_title: "AWS: aws_elasticache_replication_group"
sidebar_current: "docs-aws-datasource-elasticache-replication-group"
description: |-
  Get information on an ElastiCache Replication Group resource.
---

# Data Source: aws_elasticache_replication_group

Use this data source to get information about an Elasticache Replication Group.

## Example Usage

```hcl
data "aws_elasticache_replication_group" "bar" {
  replication_group_id = "example"
}
```

## Argument Reference

The following arguments are supported:

* `replication_group_id` – (Required) The identifier for the replication group.

## Attributes Reference

The following attributes are exported:

* `replication_group_id` - The identifier for the replication group.
* `replication_group_description` - The description of the replication group.
* `auth_token_enabled` - A flag that enables using an AuthToken (password) when issuing Redis commands.
* `automatic_failover_enabled` - A flag whether a read-only replica will be automatically promoted to read/write primary if the existing primary fails.
* `node_type` – The cluster node type.
* `number_cache_clusters` – The number of cache clusters that the replication group has.
* `snapshot_window` - The daily time range (in UTC) during which ElastiCache begins taking a daily snapshot of your node group (shard).
* `snapshot_retention_limit` - The number of days for which ElastiCache retains automatic cache cluster snapshots before deleting them.
* `port` – The port number on which the configuration endpoint will accept connections.
* `configuration_endpoint_address` - The configuration endpoint address to allow host discovery.
* `primary_endpoint_address` - The endpoint of the primary node in this node group (shard).
