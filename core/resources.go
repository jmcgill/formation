package core

type FieldType int

const (
	SCALAR FieldType = iota
	MAP
	LIST
	NESTED
)

type ScalarValue struct {
	StringValue  string
	IntegerValue int32
	BooleanValue bool
}

type InlineResource []Field

type Field struct {
	FieldType FieldType
	Key       string

	// Only one of these may be filled in
	ScalarValue ScalarValue
	NestedValue []Field
	ListValue   []Field

	// TODO(jimmy): This should most likely be a map of Fields
	MapValue map[string]ScalarValue
}

type Resource struct {
	Type   string
	Name   string
	Fields InlineResource
}
