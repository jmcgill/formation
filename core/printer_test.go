package core_test

import (
	. "formation/core"

	"io/ioutil"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func ContentsOf(filename string) string {
	golden := filepath.Join("testdata", filename)
	contents, err := ioutil.ReadFile(golden)
	Expect(err).ShouldNot(HaveOccurred())
	return string(contents)
}

func BuildNestedResource(key string, value string) InlineResource {
	resource := InlineResource{
		{
			FieldType: SCALAR,
			Key:       key,
			ScalarValue: ScalarValue{
				StringValue: value,
			},
		},
	}
	return resource
}

var _ = Describe("Printer", func() {
	It("should print an empty resource", func() {
		resource := Resource{Name: "test", Type: "simple_resource"}
		printer := Printer{}

		Expect(printer.Print(&resource)).To(Equal(ContentsOf("empty_resource.hcl")))
	})

	It("should print a resource with a scalar field", func() {
		resource := Resource{
			Name: "test",
			Type: "simple_resource",
			Fields: []Field{
				{
					FieldType: SCALAR,
					Key:       "scalar_field",
					ScalarValue: ScalarValue{
						StringValue: "scalar_value",
					},
				},
			},
		}

		printer := Printer{}

		Expect(printer.Print(&resource)).To(Equal(ContentsOf("scalar_field.hcl")))
	})

	It("should print a resource with a map field", func() {
		resource := Resource{
			Name: "test",
			Type: "simple_resource",
			Fields: []Field{
				{
					FieldType: MAP,
					Key:       "map_field",
					MapValue: map[string]ScalarValue{
						"name": {
							StringValue: "Jimmy",
						},
						"age": {
							StringValue: "31",
						},
					},
				},
			},
		}

		printer := Printer{}
		Expect(printer.Print(&resource)).To(Equal(ContentsOf("map_field.hcl")))
	})

	It("should print a resource with a simple list field", func() {
		resource := Resource{
			Name: "test",
			Type: "simple_resource",
			Fields: []Field{
				{
					FieldType: LIST,
					Key:       "simple_list_field",
					ListValue: []Field{
						{
							FieldType: SCALAR,
							ScalarValue: ScalarValue{
								StringValue: "One",
							},
						},
						{
							FieldType: SCALAR,
							ScalarValue: ScalarValue{
								StringValue: "Two",
							},
						},
					},
				},
			},
		}

		printer := Printer{}
		Expect(printer.Print(&resource)).To(Equal(ContentsOf("simple_list_field.hcl")))
	})

	It("should print a resource with a rich list field", func() {
		resource := Resource{
			Name: "test",
			Type: "simple_resource",
			Fields: []Field{
				{
					FieldType: LIST,
					Key:       "rich_list_field",
					ListValue: []Field{
						{
							FieldType:   NESTED,
							NestedValue: BuildNestedResource("nested_field", "One"),
						},
						{
							FieldType:   SCALAR,
							NestedValue: BuildNestedResource("nested_field", "Two"),
						},
					},
				},
			},
		}

		printer := Printer{}
		Expect(printer.Print(&resource)).To(Equal(ContentsOf("rich_list_field.hcl")))
	})
})
