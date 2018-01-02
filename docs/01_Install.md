# Developing Formation
To contribute to Formation, you need to install a few dependencies.

Before starting, ensure that your $GOPATH environment variable is set correctly.

**Install Formation, Terraform and the AWS Provider**


    go get github.com/jmcgill/formation
    go get github.com/hashicorp/terraform
    go get github.com/terraform-providers/terraform-provider-aws

**Install other dependencies**


    go get github.com/onsi/ginkgo
    go get github.com/onsi/gomega
    go get github.com/aws/aws-sdk-go

**Fix bad vendoring in the AWS Provider**

The AWS Provider does not follow vendoring best practices, which causes a build conflict. Deleting the vendored directory fixes this (for now).


    rm -rf $GOPATH/src/github.com/terraform-providers/terraform-provider-aws/vendor/github.com/hashicorp/terraform

**Build Formation**


    cd $GOPATH/src/github.com/jmcgill/formation
    go build

If Formation builds without errors, you’re ready to move on to [Adding an Importer to Formation!](https://github.com/jmcgill/formation/blob/master/docs/02_Importing.md)