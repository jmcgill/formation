package core

type FieldType int

const (
	SCALAR FieldType = iota
	MAP
	LIST
	NESTED
)

type InlineResource struct {
	Fields []*Field
}

func (r *InlineResource) Append(field *Field) {
	r.Fields = append(r.Fields, field)
}

type ScalarValue struct {
	StringValue  string
	IntegerValue int32
	BooleanValue bool
}

type Field struct {
	FieldType FieldType
	Key       string

	// Only one of these may be filled in
	ScalarValue *ScalarValue
	NestedValue *InlineResource
}

type Resource struct {
	Type   string
	Name   string
	Fields *InlineResource
}
