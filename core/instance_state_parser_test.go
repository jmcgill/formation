package core_test

import (
	. "github.com/jmcgill/formation/core"

	"github.com/hashicorp/terraform/terraform"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("InstanceStateParser", func() {
	It("should handle an empty instance state", func() {
		state := terraform.InstanceState{
			Attributes: map[string]string{},
		}

		expectedResource := Resource{
			Name:   "",
			Type:   "",
			Fields: &InlineResource{},
		}

		parser := InstanceStateParser{}
		Expect((*parser.Parse(&state))).To(Equal(expectedResource))
	})

	It("should parse simple fields", func() {
		state := terraform.InstanceState{
			Attributes: map[string]string{
				"simple_field": "value",
			},
		}

		expectedResource := Resource{
			Name: "",
			Type: "",
			Fields: &InlineResource{
				Fields: []*Field{
					{
						FieldType: SCALAR,
						Key:       "simple_field",
						Path:      "simple_field",
						ScalarValue: &ScalarValue{
							StringValue: "value",
						},
					},
				},
			},
		}

		parser := InstanceStateParser{}
		Expect(*parser.Parse(&state)).To(Equal(expectedResource))
	})

	It("should parse a map with multiple keys", func() {
		state := terraform.InstanceState{
			Attributes: map[string]string{
				"map_name.%":         "2",
				"map_name.map_key_1": "map_value_1",
				"map_name.map_key_2": "map_value_2",
			},
		}

		expectedResource := Resource{
			Name: "",
			Type: "",
			Fields: &InlineResource{
				Fields: []*Field{
					{
						FieldType: MAP,
						Key:       "map_name",
						Path:      "map_name",
						NestedValue: &InlineResource{
							Fields: []*Field{
								{
									FieldType: SCALAR,
									Key:       "map_key_1",
									Path:      "map_name.map_key_1",
									ScalarValue: &ScalarValue{
										StringValue: "map_value_1",
									},
								},
								{
									FieldType: SCALAR,
									Key:       "map_key_2",
									Path:      "map_name.map_key_2",
									ScalarValue: &ScalarValue{
										StringValue: "map_value_2",
									},
								},
							},
						},
					},
				},
			},
		}

		parser := InstanceStateParser{}
		Expect(*parser.Parse(&state)).To(Equal(expectedResource))
	})

	It("should parse a list with multiple scalar entries", func() {
		state := terraform.InstanceState{
			Attributes: map[string]string{
				"list_key.#": "2",
				// List prefix is this whole key
				"list_key.1234": "list_value_1",
				"list_key.1235": "list_value_2",
			},
		}

		expectedResource := Resource{
			Name: "",
			Type: "",
			Fields: &InlineResource{
				Fields: []*Field{
					{
						FieldType: LIST,
						Key:       "list_key",
						Path:      "list_key",
						NestedValue: &InlineResource{
							Fields: []*Field{
								{
									FieldType: SCALAR,
									Key:       "",
									Path:      "list_key.1234",
									ScalarValue: &ScalarValue{
										StringValue: "list_value_1",
									},
								},
								{
									FieldType: SCALAR,
									Key:       "",
									Path:      "list_key.1235",
									ScalarValue: &ScalarValue{
										StringValue: "list_value_2",
									},
								},
							},
						},
					},
				},
			},
		}

		parser := InstanceStateParser{}
		Expect(*parser.Parse(&state)).To(Equal(expectedResource))
	})

	It("should parse a list with a nested entry", func() {
		state := terraform.InstanceState{
			Attributes: map[string]string{
				"list_key.#":                        "1",
				"list_key.1234.nested_scalar_key_1": "value_1",
				"list_key.1234.nested_scalar_key_2": "value_2",
			},
		}

		//resource "type" "name" {
		//	list_key {
		//		nested_scalar_key_1 = "value_1"
		//		nested_scalar_key_2 = "value_2"
		//	}
		//}

		expectedResource := Resource{
			Name: "",
			Type: "",
			Fields: &InlineResource{
				Fields: []*Field{
					{
						FieldType: LIST,
						Key:       "list_key",
						Path:      "list_key",
						NestedValue: &InlineResource{
							Fields: []*Field{
								{
									FieldType: NESTED,
									Key:       "",
									NestedValue: &InlineResource{
										Fields: []*Field{
											{
												FieldType: SCALAR,
												Key:       "nested_scalar_key_1",
												Path:      "list_key.1234.nested_scalar_key_1",
												ScalarValue: &ScalarValue{
													StringValue: "value_1",
												},
											},
											{
												FieldType: SCALAR,
												Key:       "nested_scalar_key_2",
												Path:      "list_key.1234.nested_scalar_key_2",
												ScalarValue: &ScalarValue{
													StringValue: "value_2",
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

		parser := InstanceStateParser{}
		x := *parser.Parse(&state)
		Expect(x).To(Equal(expectedResource))
	})
})
