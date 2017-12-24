package core

// ValueType is an enum of the type that can be represented by a schema.
type ValueType int

const (
	TypeInvalid ValueType = iota
	TypeBool
	TypeInt
	TypeFloat
	TypeString
	TypeList
	TypeMap
	TypeSet
	TypeLink
)

type LinkResource struct {
	Schema map[string]*LinkSchema
}

type LinkSchema struct {
	// Type is the type of the value and must be one of the ValueType values.
	//
	// This type not only determines what type is expected/valid in configuring
	// this value, but also what type is returned when ResourceData.Get is
	// called. The types returned by Get are:
	//
	//   TypeBool - bool
	//   TypeInt - int
	//   TypeFloat - float64
	//   TypeString - string
	//   TypeList - []interface{}
	//   TypeMap - map[string]interface{}
	//   TypeSet - *schema.Set
	//
	Type ValueType

	// The following fields are only set for a TypeList, TypeSet, or TypeMap.
	//
	// Elem represents the element type. For a TypeMap, it must be a *Schema
	// with a Type of TypeString, otherwise it may be either a *Schema or a
	// *Resource. If it is *Schema, the element type is just a simple value.
	// If it is *Resource, the element type is a complex structure,
	// potentially with its own lifecycle.
	Elem interface{}

	Resource string
	Field    string
}
