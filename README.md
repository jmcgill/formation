# Formation

![Logo](https://raw.githubusercontent.com/jmcgill/formation/master/resources/logo.png)

Want to begin using Terraform, but unsure how to import your existing AWS infrastructure?

Formation is an early stages project to simplify bulk import of resources from AWS into Terraform. It uses the importing
logic that already exists in Terraform to pull down resources, ensuring 1:1 compatibility with all fields that
the Terraform AWS provider supports.

Formation is designed to one day be possible to merge into the upstream Terraform codebase, and makes some (otherwise
odd) design choices in order to preserve this compatibility.

# Current State

Formation is currently a proof of concept. The core logic is correct, but in need of some serious refactoring.

Formation can import a limited subset of resources. See [importers.go](https://github.com/jmcgill/formation/blob/master/aws/importers.go) for a list of resources that are supported.

# Known Weirdness

Formation depends on the fact that the Terraform AWS Provider uses the terraform/helpers interface, and breaks through
the plugin interface offered by Terraform to access fields that would otherwise not be exposed.

Formation also duplicates the code used for configuration/initialization to ensure that our AWS services end up
in the same state as those accessed by the AWS provider.

In an ideal world, this logic would be part of the same package as the AWS Provider, and we would not need this
complexity.

# Contributing

Contributing is _super easy_. For each resource, we need to:

1. Discover/List all resources of that type
2. Extract the unique ID (as used by Terraform) and a human readable Name

See [aws_iam_role.go](https://github.com/jmcgill/formation/blob/master/aws/aws_iam_role.go) for an example of how simple this is!

To contribute, pick any resource that doesn't yet have an imported in [importers.go](https://github.com/jmcgill/formation/blob/master/aws/importers.go) and submit a Pull Request to
implement it.

# Design Document

See [DESIGNDOC.md](https://github.com/jmcgill/formation/blob/master/DESIGNDOC.md)

# Status

Still working on the core framework - so don't expect this to do anything yet!

