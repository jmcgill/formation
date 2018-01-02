# Adding an Importer to Formation
Formation was designed from day one to be built by a team of many. One of the biggest contributions you can make today is to add an importer for an new AWS resource type. Once a new resource is supported, all future users of Formation will be able to pull control of that resource into Terraform.

So let’s get started!


## Background

Terraform defines a resource type for each unique type of infrastructure that can be managed using AWS - for example, EC2 instances are one resource (aws_instance) and Virtual Private Clouds are another (aws_vpc).

Terraform also defines resources for **connections** between types of infrastructure. For example, the aws_iam_group_policy resource creates a connection between an AWS Policy and a particular IAM Group. These resources can be a little bit harder to reason about, as they are modeled differently in the AWS API (policies are embedded within IAM Groups, not linked to them).

A full list of resources can be found at https://www.terraform.io/docs/providers/aws/index.html

Before getting started, follow the instructions to [Install Formation](https://github.com/jmcgill/formation/blob/master/docs/01 - Install.md)

## Choosing a resource

The first step is to pick a type of resource to import. You don’t **need** to have used that resource type before, however it does help a lot if you already have instances in AWS that you can test your importer against.

Our first port of call will be [importers.go](https://github.com/jmcgill/formation/blob/master/aws/importers.go). This contains the most up to date list of which resources can already be imported, and which need work.


    ...
    //"aws_internet_gateway": &AwsInternetGatewayImporter{},
    "aws_vpc": &AwsVpcImporter{},
    ...

From this, we can see that `aws_vpc` already has a valid importer, but `aws_internet_gateway` does not. Let’s uncomment that line, and then dive in and add one!

## Verifying that a resource is importable

Before writing an importer, we should check the Terraform documentation to see if a resource is importable. If a resource **is** importable, there will be an example import command at the bottom of the description of that resource.

If your chosen resource cannot be imported, I suggest putting a pin in it for now. In Part 2 of this series, I’ll explain how to import these types of resources with Formation.

Let’s check the documentation for [aws_internet_gateway](https://www.terraform.io/docs/providers/aws/r/internet_gateway.html)
[](https://www.terraform.io/docs/providers/aws/r/vpc.html)

![](https://d2mxuefqeaa7sj.cloudfront.net/s_BCCB2C5B095C39A59FA68CBA375AEFF10F7BF56D0785A2F544597C845E80B61D_1514906563490_Screen+Shot+2018-01-02+at+10.21.52+AM.png)


Fantastic! This tells us that the resource can be imported, and that the Terraform importer uses the ID of the gateway to identify **which** gateway to import.

This is also a good time to learn about the particular resource if you aren’t familiar with it. In this case, we learn that Internet Gateways are used to provide a Virtual Private Cloud with access to the public internet. Seems important - let’s import it!

## Writing an Importer

Each Formation importer performs three important roles:


1. Find all of the resources of a particular type. In our case, we want to find all Internet Gateways within a particular AWS account and region


2. Assign a unique human readable name and ID to each instance. The ID should be the same ID used by Terraform to import and manage that resource


3. Declare what other resources this resource can reference. References are used by Terraform to make it clear if a resource depends on another resource.

**Finding all Gateways**

Our first step should be to write the code that can find all Internet Gateways within an AWS account and region. Thankfully, the AWS SDK for Go makes this very simple.

You should find an already generated file [aws/aws_internet_gateway.go](https://github.com/jmcgill/formation/blob/master/aws/aws_internet_gateway.go) which contains the scaffolding needed to get started.  Most importantly, this code already contains a `Describe` method which an instance of the AWS SDK service called `svc` which has been pre-initialized.

You should use this client object to interact with the AWS SDK, as it has been configured to be compatible with Terraform’s importers.


    func (*AwsInternetGatewayImporter) Describe(meta interface{}) ([]*core.Instance, error) {
       svc :=  meta.(*AWSClient).ec2conn
      ...

Let’s flesh out the Describe method to find all those gateways! This is pretty simple - we use the DescribeInternetGateways call with a nil argument to get back a slice (array) containing information about each internet gateway in this account.


    func (*AwsInternetGatewayImporter) Describe(meta interface{}) ([]*core.Instance, error) {
       svc :=  meta.(*AWSClient).ec2conn

       // Add code to list resources here
       result, err := svc.DescribeInternetGateways(nil)
       if err != nil {
         return nil, err
       }

        existingInstances := result.InternetGateways

**Tip:** Make sure to delete the `return nil, nil` line at the top of Describe, otherwise your importer will do nothing!

**A note on pagination**

We don’t want to accidentally forget to import some instances, so it’s critical that you check whether that particular API call is paginated (a good signal is if the result includes a field called NextToken). If the call is paginated, you should use the Pages() method to ensure that you fetch all of the data. See aws/aws_iam_role.go for a good example of finding all resources using pagination.

**Mapping to Name and ID**

Great! We’ve found all the Internet Gateways in this account - now we need to come up with a human readable name and a unique ID for each one. The Unique ID should be the same ID that is used to import this type of instance. In our case, this is the Internet Gateway ID.

We store the ID and Name in a `core.Instance`  struct and return an array of these from the Describe method.


    names := make(map[string]int)
    instances := make([]*core.Instance, len(existingInstances))
    for i, existingInstance := range existingInstances {
       gatewayId := existingInstance.InternetGatewayId
       name := NameTagOrDefault(existingInstance.Tags, gatewayId, names)
       instances[i] = &core.Instance{
          Name: name,
          ID:   aws.StringValue(existingInstance.InternetGatewayId),
       }
    }

     return instances, nil

There’s two interesting things worth paying attention to in this code


1. AWS returns pointers to strings from all API calls. The `aws.StringValue` method is some nice syntactic sugar which safely dereferences those strings.


2. The `NameTagOrDefault` method is a helper function to extract the name of this resource from the Name Tag if it exists. The map passed into this helper function is used to ensure that each name is unique, since Tags do not have a uniqueness guarantee. If no Name Tag is present, the (much less human readable) InternetGatewayId is used instead.

**Declare valid references**

Our final step is to let Formation know what other resources this resource may depend on. Unfortunately, the only way to tell is by carefully reading the documentation.

For our Internet Gateway, the examples in the documentation show that it’s standard for Internet Gateways to depend on (and reference) the `id` field of an `aws_vpc` resource (in the example below, an `aws_vpc` called `main`)


    resource "aws_internet_gateway" "gw" {
      vpc_id = "${aws_vpc.main.id}"

      tags {
        Name = "main"
      }
    }

To let Formation know about this, we fill in the Links() method


    // Describes which other resources this resource can reference
    func (*AwsInternetGatewayImporter) Links() map[string]string {
       return map[string]string{
          "vpc_id": "aws_vpc.id",
       }
    }

And we’re done!

## Testing

To test your importer, you can ask Formation to only import instances of that particular resource type:


    go run main.go -resource=aws_internet_gateway

This needs to be run in an environment which has the AWS_ACCESS_KEY_ID and AWS_SECRET_ACCESS_KEY and AWS_REGION environment variables set for an AWS account which contains at least one internet gateway.

This will produce two files:


1. `aws_internet_gateway.tf` will contain a description of each instance that was imported
2. `terraform.tfstate` will be a valid tfstate file structured so that terraform can now manage all of these instances

Finally, we can confirm that everything worked as expected by running `terraform plan` and verifying that there is no drift.

If you see the message `No changes. Infrastructure is up-to-date.` then everything has gone according to plan. Congratulations! You’ve imported your first AWS resource, and made a huge contribution to the Formation project. Thank you!

