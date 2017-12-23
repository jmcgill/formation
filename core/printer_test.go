package core_test

import (
	"fmt"
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

//func BuildNestedResource(key string, value string) InlineResource {
//	resource := InlineResource{
//		{
//			FieldType: SCALAR,
//			Key:       key,
//			ScalarValue: ScalarValue{
//				StringValue: value,
//			},
//		},
//	}
//	return resource
//}

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
			Fields: &InlineResource{
				Fields: []*Field{
					{
						FieldType: SCALAR,
						Key:       "scalar_field",
						ScalarValue: &ScalarValue{
							StringValue: "scalar_value",
						},
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
			Fields: &InlineResource{
				Fields: []*Field{
					{
						FieldType: MAP,
						Key:       "map_field",
						NestedValue: &InlineResource{
							Fields: []*Field{
								{
									FieldType: SCALAR,
									Key:       "age",
									ScalarValue: &ScalarValue{
										StringValue: "31",
									},
								},
								{
									FieldType: SCALAR,
									Key:       "name",
									ScalarValue: &ScalarValue{
										StringValue: "Jimmy",
									},
								},
							},
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
			Fields: &InlineResource{
				Fields: []*Field{
					{
						FieldType: LIST,
						Key:       "simple_list_field",
						NestedValue: &InlineResource{
							Fields: []*Field{
								{
									FieldType: SCALAR,
									Key:       "",
									ScalarValue: &ScalarValue{
										StringValue: "One",
									},
								},
								{
									FieldType: SCALAR,
									Key:       "",
									ScalarValue: &ScalarValue{
										StringValue: "Two",
									},
								},
							},
						},
					},
				},
			},
		}
		printer := Printer{}
		Expect(printer.Print(&resource)).To(Equal(ContentsOf("simple_list_field.hcl")))
	})

	It("should print a resource with a nested list field", func() {
		resource := Resource{
			Name: "test",
			Type: "simple_resource",
			Fields: &InlineResource{
				Fields: []*Field{
					{
						FieldType: LIST,
						Key:       "rich_list_field",
						NestedValue: &InlineResource{
							Fields: []*Field{
								{
									FieldType: NESTED,
									Key:       "",
									NestedValue: &InlineResource{
										Fields: []*Field{
											{
												FieldType: SCALAR,
												Key:       "nested_field",
												ScalarValue: &ScalarValue{
													StringValue: "One",
												},
											},
											{
												FieldType: SCALAR,
												Key:       "nested_field",
												ScalarValue: &ScalarValue{
													StringValue: "Two",
												},
											},
										},
									},
								},
								{
									FieldType: NESTED,
									Key:       "",
									NestedValue: &InlineResource{
										Fields: []*Field{
											{
												FieldType: SCALAR,
												Key:       "nested_field",
												ScalarValue: &ScalarValue{
													StringValue: "Three",
												},
											},
											{
												FieldType: SCALAR,
												Key:       "nested_field",
												ScalarValue: &ScalarValue{
													StringValue: "Four",
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		}

		printer := Printer{}
		fmt.Printf(printer.Print(&resource))
		Expect(printer.Print(&resource)).To(Equal(ContentsOf("rich_list_field.hcl")))
	})
})
