# Formation
*A better Terraform importer*

Terraform is an increasingly popular and robust way to manage infrastructure as code, allowing for changes to infrastructure to be version controlled and peer reviewed.

While it makes sense for new projects to use Terraform, importing existing AWS infrastructure currently requires manually writing out the

Formation is designed to make importing resources from any cloud into Terraform simple. By building on top of the existing Terraform providers, Formation is able to guarantee that it can import all properties supported by Terraform, and can serialize tfstate files with guaranteed compatibility.

## Related Projects

Terraforming (https://github.com/dtan4/terraforming) is a Ruby CLI which uses the AWS Ruby Library to read the details of assets stored in AWS and populate template files.

It supports a small subset of AWS resources, and is missing many properties from basic objects (e.g. tags or cors_headers for S3 buckets). Adding these properties requires effectively replicating much of the work that has already gone into reading AWS resources in Terraform.

## Goals

Terraform should perform two basic actions:


1. Discover and import all resources managed by a particular provider. Importing means generating both a .tfstate file **and** a fully populated .tf file.


2. Create links between imported resources. For example, if an S3 bucket references a particular IAM Role, we should create a link between those two resources to ensure that updates are applied in the right order, and to make it easier to navigate between infrastructure.


  As a concrete example:


    resource "aws_s3_bucket" "my_bucket" {
      iam_role = "arn:27387123/238742272wow293y2cool08231arn"
    }


  should become:


    resource "aws_s3_bucket" "my_bucket" {
      iam_role = "${aws_iam_role.superuser}"
    }

We will address these two problems as orthogonal phases.

# Basic Architecture - Importing

Terraform recently split the logic which reads and mutates resources in a particular cloud provider (e.g. AWS) into a set of providers.

Each Provider exports the `terraform.ResourceProvider` interface, which provides methods to query the current state of, or update, a set of resources.

Within a given Provider, each resource implements the `schema.Resource` type, defining a schema for HCL, the configuration language used to drive Terraform, and methods to Update, Read, Delete and Create that resource. For AWS, these methods are implemented using the AWS Go SDK.

For our purposes, the `Refresh` method is the most interesting. This takes the details of a particular resource - for example, a single S3 bucket - to query as an `InstanceInfo` struct, and calls the Read() method for that resource. This returns an `InstanceState` struct, which **mostly** containing the current state of that resource in AWS. 

There are three important gotchas:


1. The `InstanceState` struct passed to Refresh must already have the ID of the resource, as used by Terraform, populated. If not, Terraform will assume we are creating a new resource. For some resources this ID is the resource name, for others it is an opaque ID e.g. an AWS ARN.


2. Most, but not all, properties of a resource will be read by Terraform and populated in the InstanceState, however some properties are not read unless they are set on the `InstanceState` before it is passed to Refresh. This is true, for example, for `aws_s3_bucket.policy` as policies can either be ended inline or as an external resource.


3. Some fields in `InstanceState` will be set to their default values. It should be possible to detect these fields/values by comparing to the state returned to the resource schema exported by the Provider.


## Parsing InstanceState

The `InstanceState` object contains everything we need to know to create a Terraform object - unfortunately, it is not particularly well structured for spitting out structured resources.

The `InstanceState` represents a deeply nested terraform resource as a flattened list of key paths, with some additional entries to differentiate arrays and maps. Let’s start with a concrete example:


    resource "type" "name" {
      scalar_key = scalar_value
      scalar_array_key = [scalar_array_value_1, scalar_array_value_2]
      map_name {
        map_key_1 = map_value_1
        map_key_2 = map_value_2
      }
      
      # The first element in this array
      expanded_array_key {
        nested_scalar_key = nested_scalar_value
      }
    
      # The second element in this array
      expanded_array_key {
        nested_scalar_key = nested_scalar_value
      }  
    }

Is encoded as:


    # Scalar keys and values are nice and simple
    scalar_key = scalar_value
    
    # Arrays always contain an entry defining how entries are in the array (N) followed by
    # N entries with the values in that array and a unique index for each value.
    scalar_array_key.# = 2
    scalar_array_key.1234 = scalar_array_value_1 
    scalar_array_key.5678 = scalar_array_value_2
    
    # Maps contain an entry definining how many keys are in that map (N) followed by N
    # entries with the keys and values.
    map_name.% = 2
    map_name.map_key_1 = map_value_1
    map_name.map_key_2 = map_value_2
    
    # Arrays can also nest objects. The semantics are the same as an array and scalar
    # combined.
    expanded_array_key.# = 2
    expanded_array_key.6666.nested_scalar_key = nested_scalar_value_1
    expanded_array_key.8888.nested_scalar_key = nested_scalar_value_2

Using a simple recursive parser, we can easily reverse this mapping by walking the keys, and then reconstruct the original nested resource.

Once reconstructed, this resource will be compared to the nested Schema returned by the Provider, and any field with its default value will be removed from the object.

Finally the `InstanceState` can be serialized as a tfstate file using the serializer in `terraform/state.go` and the parsed resource can be serialized as a Terraform resource using some nice gnarly string concatenation. 

## Representing Resources

Resources will be represented using the following structs


    type FieldType int
    
    const (
            SCALAR FieldType = iota
            MAP
            LIST
    )
    
    type ScalarValue struct {
      String_value string
      Integer_value int32
      Boolean_value bool
    }
    
    type Field struct {
      FieldType FieldType  
      Key string
        
      // Only one of these may be filled in
      Scalar_value ScalarValue
      Field_list_value []Field
      Scalar_list_value []ScalarValue
      Map_value map[string]ScalarValue
    }
    
    type Resource struct {
      Name string
      Fields []Field
    }
## Resource Discovery

In order to Import a resource, we must first know its identifier. Discovering these identifiers is different for each resource type, and will be the most significant piece of functionality.

For each resource type in AWS, we will implement an instance of the following interface.


    type ImportedResource {
      Name string
      State *InstanceState
    }
    
    type ResourceParser interface {
      // Discover all resources and return an ImportedResource for each one
      PopulateInstanceStates() []ImportedResource
    }

An implementation of this interface will be responsible for discovering all resources in a particular account, for example by using the AWS API, and creating an ImportedResource for each one with the ID in InstanceState set to the value which terraform expects for that resource type.

# Basic Architecture - Linking

Linking is a slightly more complex problem as it requires an understanding of the possible relationships between resources.

For example, an iam_role is linked using its arn as an identifier, while an S3 bucket is linked using the bucket name as an identifier. It is possible that an S3 bucket name is the same as the ARN for an iam_role, so the identifier to link to is not sufficient in order to work out which link to construct.

In order to construct the correct set of links, we will perform two steps:


1. During the import phase (above) we will create an index for every resource, keyed using a composite key constructed from the resource type, field key and field value, for every field in the root of a resource.


2. We will semantically declare the set of valid links between resource types. There is unfortunately no better way to do this than reading through the Terraform documentation.
## Indexing

Given an opaque identifier in a resource, we want to understand what other resources that identifier might be linking to. We will therefore create an index that allows O(1) lookup of potential resources by identifier, assuming we know the instance type and field name that we expect to link to (more on this later).

For example, the following resource:


    resource "aws_s3_bucket" "my_bucket" {
        name: "my_name",
        arn: "my_arn"
    }

Will produce the following index:


    aws_s3_bucket.name.my_name => my_bucket
    aws_s3_bucket.arn.my_arn => my_bucket


## Resolving Links

In order to determine how to link resources, we need to declare the types of links that can exist between resources. This will be declared as part of the `ResourceParser`  struct in a struct returned by the `GetLinks` method.

This struct will conform to the following interface, which is designed to one day be easily mergeable with the `schema.Schema` struct in Terraform.

`// TODO(jimmy): Linking inside policies`


    type ValueType int
    
    const (
            TypeLink = iota
            TypeList
            TypeMap
    )
    
    type LinkResource struct {
      Schema map[string]*LinkSchema
    }
    
    type Link struct {
      Type string
      Field string
    }
    
    type LinkSchema struct {
      Type ValueType
      Links []Link
      Elem LinkResource
    }


    Schema: map[string]*LinkSchema{
       "iam_role": {
          Type:      schema.TypeLink,
          Links: [
             Type: "aws_iam_role",
             Field: "arn"
          ]
       },
       
       // A link inside a nested resource. This follows the same structure as
       // Resource Schemas in Terraform.
       "users": {
         Type:     schema.TypeList,
         Elem: &LinkResource{
            Schema: map[string]*LinkSchema{
               "iam_role": {
                  Type:     schema.TypeLink,
                  Links: [
                     Type: "aws_iam_role",
                     Field: "arn"
                  ]
               }
            }
         }
      }
    }

    
To resolve a link, we will walk the `LinkSchema` to find the definition of the current field, and then construct a key from `Link.Type` + `Link.Field` + the current value in the imported resource to look up 

It is possible for there to be more than one valid link for a given field type - in this case the declaration will contain multiple `Link` definitions, and the first matching definition will prevail.


## Storing Links

For now, the resolved link will simply be stored as a Terraform link string with the format `${resource_type.resource_name.resource_field`.


# Example Resource Imported

The following is an example of an S3 resource importer


    import (
        "fmt"
        "github.com/aws/aws-sdk-go/aws"
        "github.com/aws/aws-sdk-go/aws/session"
        "github.com/aws/aws-sdk-go/service/s3"
    )
    
    type S3Importer struct {
    }
    
    func (i *S3Importer) PopulateInstanceStates() []ImportedResource {
      var r []ImportedResource
      
      // TODO(jimmy): How best to pass in configuration?
      svc := s3.New(session.New(), &aws.Config{Region: aws.String("us-west-2")})
    
      // Example sending a request using the ListBucketsRequest method.
      req, resp := client.ListBucketsRequest()
      err := req.Send()
      if err {
        fmt.Println("Error listing buckets")
      }
      
      for _, bucket := range resp.Buckets {
        state := InstanceState{}
        state.Id = bucket.Name
    
        // Force policies to be inlined if they exist
        state.Attributes["policy"] = ""
        
        instance := ImportedResource {
          Name: bucket.Name,
          InstanceState: &state
        }
        r = append(r, instance)
      }
      
      return instance
    }
    
    func (i *S3Importer) GetLinks() map[string]*LinkSchema {
      return map[string]*LinkSchema{
       "logging": {
         Type:     schema.TypeList,
         Elem: &LinkResource{
            Schema: map[string]*LinkSchema{
               "target_bucket": {
                  Type:     schema.TypeLink,
                  Links: [
                     Type: "aws_s3_bucket",
                     Field: "id"
                  ]
               }
            }
         }
      },
       "replication_configuration": {
         Type:     schema.TypeList,
         Elem: &LinkResource{
            Schema: map[string]*LinkSchema{
               "role": {
                  Type:     schema.TypeLink,
                  Links: [
                     Type: "aws_iam_role",
                     Field: "arn"
                  ]
               },
               
               "rules": {
                  Type:   schema.List,
                  Elem: &LinkResource {
                    Schema: map[string]*LinkSchema {
                      "destination": {
                        Type: schema.List,
                        Elem: &LinkResouce{
                          Schema: map[string]*LinkSchema {
                            "bucket": {
                              Type:   schema.Link,
                              Links: [
                                Type: "aws_s3_bucket",
                                Field: "arn"
                              ]
                            }
                          }
                        }
                      }
                    }
                  }
               }
            }
         }
      }
    }

Alternate representation:


    "logging[].target_bucket" => "aws_s3_bucket.id"
    "replication_configuration[].role" => "aws_s3_bucket.arn"
    "replication_configuration[].rules[].destination[].bucket" = "aws_s3_bucket.arn"

