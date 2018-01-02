# Formation

![Logo](https://raw.githubusercontent.com/jmcgill/formation/master/resources/logo.png)

Want to begin using Terraform, but unsure how to import your existing AWS infrastructure?

Formation is an early stages project to simplify bulk import of resources from AWS into Terraform. It uses the importing
logic that already exists in Terraform to pull down resources, ensuring 1:1 compatibility with all fields that
the Terraform AWS provider supports.

Formation is designed to one day be possible to merge into the upstream Terraform codebase, and makes some (otherwise
odd) design choices in order to preserve this compatibility.

# Current State

Formation can import a limited subset of resources. See [importers.go](https://github.com/jmcgill/formation/blob/master/aws/importers.go) for a list of resources that are supported.

# Contributing

Contributing is _super easy_ and pull requests are _very welcome_. To contribute, follow the guides below!

1. [Installing Formation](https://github.com/jmcgill/formation/blob/master/docs/01_Install.md)
2. [Adding an Importer](https://github.com/jmcgill/formation/blob/master/docs/02_Importing.md)

# Known Weirdness

Formation depends on the fact that the Terraform AWS Provider uses the terraform/helpers interface, and breaks through
the plugin interface offered by Terraform to access fields that would otherwise not be exposed.

Formation also duplicates the code used for configuration/initialization to ensure that our AWS services end up
in the same state as those accessed by the AWS provider.

In an ideal world, this logic would be part of the same package as the AWS Provider, and we would not need this
complexity.

# Design Document

See [DESIGNDOC.md](https://github.com/jmcgill/formation/blob/master/DESIGNDOC.md)

