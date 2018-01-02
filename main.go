package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"flag"
	"io/ioutil"

	"github.com/jmcgill/formation/aws"
	"github.com/jmcgill/formation/core"
	"os"
	"strings"

	"github.com/hashicorp/terraform/config/configschema"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	aws2 "github.com/terraform-providers/terraform-provider-aws/aws"
)

type UIInput struct {
}

func (u *UIInput) Input(opts *terraform.InputOpts) (string, error) {
	return "us-west-2", nil
}

// Compare a converted resource to the Terraform Schema for that resource type, and
// mark any fields which we know are computed. These are fields which we want to exist in our
// .tfstate files, but not in our .tf files
func MarkComputedFields(r *core.InlineResource, s *configschema.Block) *core.InlineResource {
	for _, f := range r.Fields {
		// This is a scalar object
		if f.FieldType == core.SCALAR && f.Path != "id" {
			f.Computed = s.Attributes[f.Key].Computed
			continue
		}

		// The children of maps are not known ahead of time, so we can't walk further
		// in the schema.
		if f.FieldType == core.MAP {
			continue
		}

		if f.FieldType == core.LIST && f.NestedValue.Fields[0].FieldType == core.NESTED {
			MarkComputedFields(f.NestedValue, &s.BlockTypes[f.Key].Block)
		}

		// Assume that lists of non-objects are never computed
		// This is definitely wrong, but I need a case study to work out what this should look like
		if f.FieldType == core.LIST && f.NestedValue.Fields[0].FieldType != core.NESTED {
			continue
			// MarkComputedFields(f.NestedValue, &s.Attributes[f.Key].Block)
		}

		if f.FieldType == core.NESTED {
			MarkComputedFields(f.NestedValue, s)
		}
	}

	return r
}

func ValueSet(r *core.InlineResource, key string) bool {
	for _, f := range r.Fields {
		if f.Key == key {
			return true
		}
	}
	return false
}

func AppendToInstancePath(prefix string, suffix string) string {
	if prefix == "" {
		return suffix
	}
	return prefix + "." + suffix
}

func DecorateWithDefaultFields(instanceState *terraform.InstanceState, r *core.InlineResource, terraformSchema map[string]*schema.Schema, path string) *core.InlineResource {
	// Are there any default fields at this level in the schema that are not set in our resource?
	for key, v := range terraformSchema {
		if v.Default == nil {
			continue
		}

		if ValueSet(r, key) {
			continue
		}

		var defaultValue string

		// TODO(jimmy): Switch to a TYPE enum - can pull this straight out of the schema object
		isBool := false

		switch v.Default.(type) {
		case bool:
			defaultValue = strconv.FormatBool(v.Default.(bool))
			isBool = true
		case string:
			defaultValue = v.Default.(string)
		}

		instancePath := AppendToInstancePath(path, key)
		instanceState.Attributes[instancePath] = defaultValue

		field := &core.Field{
			FieldType: core.SCALAR,
			Path: instancePath,
			Key:     key,
			ScalarValue: &core.ScalarValue{
				StringValue: defaultValue,
				IsBool:      isBool,
			},
		}
		r.Fields = append(r.Fields, field)
	}

	for _, f := range r.Fields {
		instancePath := AppendToInstancePath(path, f.Key)

		// HACK: Assume all resources are nested
		if f.FieldType == core.LIST && f.NestedValue.Fields[0].FieldType == core.NESTED {
			// Scalar lists will never have a default value, so we don't need to expand the parent key when
			// recursing into a list. Instead we will rely on each nested object to do this.
			DecorateWithDefaultFields(instanceState, f.NestedValue.Fields[0].NestedValue, terraformSchema[f.Key].Elem.(*schema.Resource).Schema, instancePath)
		}
	}

	return r
}

type FieldIndexEntry struct {
	path     string
	resource *core.Resource
}
type FieldIndex map[string][]*FieldIndexEntry

func IndexResource(resource *core.Resource, index FieldIndex) {
	IndexFields(resource, resource.Fields, index)
}

// Index all fields
// TODO(jimmy): Implement a visitor pattern that allows this to be done during the parsing pass.
func IndexFields(resource *core.Resource, r *core.InlineResource, index FieldIndex) {
	for _, f := range r.Fields {
		// This is a scalar object
		if f.FieldType == core.SCALAR {
			// Don't index null values
			if f.ScalarValue.StringValue == "" {
				continue
			}

			e := &FieldIndexEntry{
				path:     resource.Type + "." + f.Path,
				resource: resource,
			}
			index[f.ScalarValue.StringValue] = append(index[f.ScalarValue.StringValue], e)
			continue
		}

		// The children of maps are not known ahead of time, and are never referenced from other resources.
		if f.FieldType == core.MAP {
			continue
		}

		if f.FieldType == core.LIST || f.FieldType == core.NESTED {

			IndexFields(resource, f.NestedValue, index)
		}
	}
}

func  FindLink(index FieldIndex, value string, allowedPath string) (*core.Resource, bool) {
	if fields, ok := index[value]; ok {
		for _, field := range fields {
			if field.path == allowedPath {
				return field.resource, true
			}
		}
	}

	return nil, false
}

func LinkFields(root* core.Resource, r *core.InlineResource, links map[string]string, index FieldIndex) *core.InlineResource {
	RecursivelyLinkFields(root, r, links, index, "")
	return r
}

func RecursivelyLinkFields(root* core.Resource, r *core.InlineResource, links map[string]string, index FieldIndex, path string) {
	for _, f := range r.Fields {
		// This is a scalar object
		if f.FieldType == core.SCALAR {
			var z string
			if path != "" {
				z = path + "." + f.Key
			} else {
				z = f.Key
			}

			// This field _can_ link to another resource
			if allowedPath, ok := links[z]; ok {
				if resource, ok := FindLink(index, f.ScalarValue.StringValue, allowedPath); ok {
					// Avoid self links
					if resource.Name != root.Name {
						// Substitute in the resource name
						resolvedPath := strings.Replace(allowedPath, resource.Type, resource.Type+"."+resource.Name, 1)
						f.Link = resolvedPath
					}
				}
			}
			continue
		}

		// Can a Map ever contain a field which links to another resource? TBD
		if f.FieldType == core.MAP {
			continue
		}

		if f.FieldType == core.LIST {
			var z string
			if path != "" {
				z = path + "." + f.Key
			} else {
				z = f.Key
			}

			RecursivelyLinkFields(root, f.NestedValue, links, index, z)
		}

		if f.FieldType == core.NESTED {
			RecursivelyLinkFields(root, f.NestedValue, links, index, path)
		}
	}
}

type ImportedResource struct {
	resource *core.Resource
	state *terraform.InstanceState
}

func main() {
	tfstate := flag.String("tfstate", "", "Path to an existing tfstate file to merge")
	resourceToImport := flag.String("resource", "", "A specific resource type to import")
	flag.Parse()

	// TODO(jimmy): Bundle these into an object
	allResources := make(map[string][]*ImportedResource)

	// TODO(JIMMY): hide this away in a struct
	index := make(FieldIndex)

	importers := aws.Importers()

	// Restrict to a specific resource type, if requested
	if *resourceToImport != "" {
		importers = map[string]core.Importer{
			*resourceToImport: importers[*resourceToImport],
		}
	}

	// Configure Terraform Plugin
	provider := aws2.Provider()

	// A duplicate Provider which is guaranteed to be configured in the same way as
	// the AWS probvider above. This is needed so that we can access otherwisr private
	// fields that are configured during initialization.
	localProvider := aws.Provider()
	c := terraform.NewResourceConfig(nil)
	localProvider.Input(&UIInput{}, c)

	err := localProvider.Configure(c)
	if err != nil {
		panic("Error configuring internal provider")
	}
	localSchemaProvider := localProvider.(*schema.Provider)

	c = terraform.NewResourceConfig(nil)
	provider.Input(&UIInput{}, c)
	err = provider.Configure(c)
	if err != nil {
		log.Fatalf("Error while configuring provider %s", err)
	}

	// For each importer
	for resourceType, importer := range importers {
		instances, err := importer.Describe(localSchemaProvider.Meta())
		if err != nil {
			panic(err)
		}

		for _, instance := range instances {
			var instancesToImport []*terraform.InstanceState
			importViaTerraform := true

			instanceInfo := &terraform.InstanceInfo{
				// Id is a unique name to represent this instance. This is not related
				// to InstanceState.ID in any way.
				Id: instance.ID,

				// Type is the resource type of this instance
				Type: resourceType,
			}

			if patchyImporter, ok := importer.(core.PatchyImporter); ok {
				instancesToImport, importViaTerraform, err = patchyImporter.Import(instance, localSchemaProvider.Meta())
			}

			if importViaTerraform {
				instancesToImport, err = provider.ImportState(instanceInfo, instance.ID)
				if err != nil {
					log.Printf("[ERROR] Error importing instance: %s", err)
					continue
				}
			}

			// TODO(jimmy): EC2: Add an Import function for when things don't quite go right. Could source UserData there.
			//instancesToImport[0].Attributes["user_data_base64"] = "sentinal"

			for _, instanceToImport := range instancesToImport {
				instanceState, err := provider.Refresh(instanceInfo, instanceToImport)
				if err != nil {
					log.Fatal("Error refreshing Instance State")
				}

				if patchyImporter, ok := importer.(core.PatchyImporter); ok {
					instanceState = patchyImporter.Clean(instanceState, localSchemaProvider.Meta())
				}

				// Convert this resource from Terraform's internal format to a Formation Resource
				parser := core.InstanceStateParser{}
				resource := parser.Parse(instanceState)

				/// Fill in name and type
				resource.Name = instance.Name
				resource.Type = resourceType

				// Get the schema for this resource
				request := &terraform.ProviderSchemaRequest{
					ResourceTypes: []string{resourceType},
				}
				s, _ := provider.GetSchema(request)

				// Mark computed fields - we don't want to output these
				MarkComputedFields(resource.Fields, s.ResourceTypes[resourceType])

				// To get the resource schema we need to poke into the internal implementation of the AWS provider
				schemaProvider := provider.(*schema.Provider)
				DecorateWithDefaultFields(instanceState, resource.Fields, schemaProvider.ResourcesMap[resourceType].Schema, "")

				// Store this resource for later
				allResources[resourceType] = append(allResources[resourceType], &ImportedResource{
					resource: resource,
					state:    instanceState,
				})

				// Index this resource
				IndexFields(resource, resource.Fields, index)
			}
		}
	}

	// At this point, all resources have been index
	for resourceType, resources := range allResources {
		f, err := os.Create(resourceType + ".tf")
		defer f.Close()

		if err != nil {
			log.Fatalf("Error creating file for resource %s\n", resourceType)
		}

		for i, importedResource := range resources {
			resource := importedResource.resource
			LinkFields(resource, resource.Fields, importers[resource.Type].Links(), index)

			printer := core.Printer{}
			printer.PrintToFile(f, resource)

			// Space out resources for readability
			if i != len(resources)-1 {
				fmt.Fprint(f, "\n\n")
			}
		}
	}

	// Write TFState
	state := terraform.State{
		Version: 3,

		// TODO(jimmy): Pull this out of the terraform Context object
		TFVersion: "0.11.1",

		Serial: 1,

		// TODO(jimmy): Work out how this is generated
		Lineage: "Formation",

		Modules: []*terraform.ModuleState{
			{
				Path: []string{
					"root",
				},
			},
		},
	}
	state.Modules[0].Resources = make(map[string]*terraform.ResourceState)

	if *tfstate != "" {
		contents, err := ioutil.ReadFile(*tfstate)
		if err != nil {
			log.Fatal("Error reading existing TFState file")
		}

		err = json.Unmarshal(contents, &state);
		if err != nil {
			log.Fatal("Error unmarshaling JSON")
		}

		// Clear any resources being imported, in case we've changed the way
		// their key is constructed
		if *resourceToImport != "" {
			for key, _ := range state.Modules[0].Resources {
				if strings.HasPrefix(key, *resourceToImport) {
					delete(state.Modules[0].Resources, key)
				}
			}
		}
	}

	for _, resources := range allResources {
		for _, importedResource := range resources {
			resource := importedResource.resource
			r := &terraform.ResourceState{
				Type:    resource.Type,
				Primary: importedResource.state,
				Provider: "provider.aws",
			}
			state.Modules[0].Resources[resource.Type+"."+resource.Name] = r
		}
	}

	f, err := os.Create("terraform.tfstate")
	if err != nil {
		log.Fatal("Failure to create TFState file")
	}
	defer f.Close()

	j, _ := json.MarshalIndent(state, "", "    ")
	f.Write(j)
}
