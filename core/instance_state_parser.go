package core

import (
	"github.com/jmcgill/formation/terraform_helpers"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform/terraform"
)

// HERE BE DRAGONS
// Do not read this code - it needs a complete rewrite to be a top down parser and to better differentiate
// between LISTS, SETS and MAPS.

type ParentType int

const (
	PARENT_ROOT ParentType = iota
	PARENT_MAP
	PARENT_LIST
	PARENT_NESTED
)

// We keep a queue of state as we traverse the resource, so that we can walk
// up and down nested resources as required.
type State struct {
	parent            *InlineResource
	parentType        ParentType
	remainingChildren int
	depth             int
	prefix            string
}

type InstanceStateParser struct {
	stateStack []*State
}

func (p *InstanceStateParser) state() *State {
	return p.stateStack[len(p.stateStack)-1]
}

func (p *InstanceStateParser) pushState(state *State) {
	p.stateStack = append(p.stateStack, state)
}

func (p *InstanceStateParser) popState() *State {
	s := p.stateStack[len(p.stateStack)-1]
	p.stateStack = p.stateStack[:len(p.stateStack)-1]
	return s
}

func (p *InstanceStateParser) currentResource() *InlineResource {
	return p.state().parent
}

func (p *InstanceStateParser) Parse(state *terraform.InstanceState) *Resource {
	resource := &Resource{}

	// TODO(jimmy): Wrapping InstanceState like this is no longer valuable - convert to a helper
	// method which returns a SortedInstanceState
	var wrappedState terraform_helpers.InstanceState
	wrappedState = terraform_helpers.InstanceState(*state)
	sortedState := wrappedState.ToSorted()

	s := State{
		remainingChildren: 0,
		parent:            new(InlineResource),
		parentType:        PARENT_ROOT,
	}
	p.pushState(&s)
	resource.Fields = s.parent

	for _, a := range sortedState.Attributes {
		p.parseAttribute(a.Key, a.Value)
	}

	return resource
}

//# Scalar keys and values are nice and simple
//scalar_key = scalar_value
//
//# Arrays always contain an entry defining how entries are in the array (N) followed by
//# N entries with the values in that array and a unique index for each value.
//scalar_array_key.# = 2
//scalar_array_key.1234 = scalar_array_value_1
//scalar_array_key.5678 = scalar_array_value_2
//
//# Maps contain an entry definining how many keys are in that map (N) followed by N
//# entries with the keys and values.
//map_name.% = 2
//map_name.map_key_1 = map_value_1
//map_name.map_key_2 = map_value_2
//
//# Arrays can also nest objects. The semantics are the same as an array and scalar
//# combined.
//expanded_array_key.# = 2
//expanded_array_key.6666.nested_scalar_key = nested_scalar_value_1
//expanded_array_key.8888.nested_scalar_key = nested_scalar_value_2

func (p *InstanceStateParser) parseAttribute(attribute string, value string) {
	originalAttribute := attribute

	// We cannot establish the prefix for a list item until we first encounter it.
	if p.state().parentType == PARENT_LIST && p.state().prefix == "" {
		// Determine whether this is a list of scalar objects or a list of nested objects.
		parts := strings.Split(attribute, ".")

		// For scalar lists this is just a unique key for this entry. For nested objects
		// the first entry contains both that unique key _and_ the first nested field.
		listPrefix := strings.Join(parts[p.state().depth-1:], ".")

		if strings.ContainsRune(listPrefix, '.') {
			// Create a new nested resource
			field := &Field{
				FieldType:   NESTED,
				Key:         "",
				NestedValue: new(InlineResource),
			}
			p.currentResource().Append(field)

			// Push this entry onto the stack.
			s := &State{
				parent:            field.NestedValue,
				remainingChildren: 0,
				parentType:        PARENT_NESTED,
				depth:             p.state().depth + 0,
				prefix:            strings.Join(parts[:p.state().depth], ".") + ".",
			}
			p.pushState(s)
		} else {
			// This is a list of scalar objects
			p.state().prefix = strings.Join(parts[:p.state().depth], ".")
		}
	}

	if !strings.HasPrefix(attribute, p.state().prefix) {
		for !strings.HasPrefix(attribute, p.state().prefix) {
			// We should also pop if we've reached the end of key
			// TODO(jimmy): Is this always safe?
			if p.stateStack[len(p.stateStack)-2].parentType == PARENT_LIST {
				_ = p.popState()
			}

			// This is a sign that we've either reached the end of a map, or the end of a single
			// item in a list.
			p.state().remainingChildren -= 1
			if p.state().remainingChildren <= 0 {
				// Pop an entry off the stack
				_ = p.popState()
			} else {
				break
			}
		}

		// If this is a list, we need to reset the matched prefix.
		// TODO(jimmy): Factor this out into a single method
		if p.state().parentType == PARENT_LIST {
			// This is a bad thing to do, as it clears any previous nesting
			// but it's OK because the next field we see will restore it.
			// Determine whether this is a list of scalar objects or a list of nested objects.
			parts := strings.Split(attribute, ".")

			// For scalar lists this is just a unique key for this entry. For nested objects
			// the first entry contains both that unique key _and_ the first nested field.
			listPrefix := strings.Join(parts[p.state().depth-1:], ".")

			if strings.ContainsRune(listPrefix, '.') {
				// This is a list of nested resource
				// Create a new nested resource
				field := &Field{
					FieldType:   NESTED,
					Key:         "",
					NestedValue: new(InlineResource),
					Path:        originalAttribute,
				}
				p.currentResource().Append(field)

				// Push this entry onto the stack.
				s := &State{
					parent:            field.NestedValue,
					remainingChildren: 0,
					parentType:        PARENT_NESTED,
					depth:             p.state().depth + 0,
					prefix:            strings.Join(parts[:p.state().depth], ".") + ".",
				}
				p.pushState(s)
			} else {
				// This is a list of scalar objects
				p.state().prefix = strings.Join(parts[:p.state().depth], ".")
			}
		}
	}
	attribute = strings.TrimPrefix(attribute, p.state().prefix)

	if !strings.ContainsRune(attribute, '.') || p.state().parentType == PARENT_MAP {
		if p.state().parentType == PARENT_MAP {
			p.state().remainingChildren -= 1
		}
		p.parseSimpleAttribute(attribute, originalAttribute, value)
	} else {
		parts := strings.Split(attribute, ".")
		fieldName, parts := parts[0], parts[1:]

		// Is this the beginning of a map?
		if parts[0] == "%" {
			children, _ := strconv.Atoi(value)
			if children != 0 {
				// Is this the map declaration?
				field := &Field{
					FieldType:   MAP,
					Key:         fieldName,
					NestedValue: new(InlineResource),
					Path:        strings.Replace(originalAttribute, ".%", "", -1),
				}
				p.currentResource().Append(field)

				// Push this entry onto the stack.
				s := &State{
					parent:            field.NestedValue,
					remainingChildren: children,
					parentType:        PARENT_MAP,
					depth:             p.state().depth + 1,
					prefix:            p.state().prefix + fieldName + ".",
				}
				p.pushState(s)
			}
		}

		// Is this the beginning of a list?
		if parts[0] == "#" {
			children, _ := strconv.Atoi(value)
			if children != 0 {
				field := &Field{
					FieldType:   LIST,
					Key:         fieldName,
					NestedValue: new(InlineResource),
					Path:        strings.Replace(originalAttribute, ".#", "", -1),
				}
				p.currentResource().Append(field)

				// Push this entry onto the stack.
				s := &State{
					parent:            field.NestedValue,
					remainingChildren: children,
					parentType:        PARENT_LIST,
					depth:             p.state().depth + 2,
					prefix:            "", //p.state().prefix + fieldName + ".",
				}
				p.pushState(s)
			}
		}
	}
}

func (p *InstanceStateParser) parseSimpleAttribute(attribute string, path string, value string) {
	fieldValue := ScalarValue{StringValue: value}
	field := &Field{
		FieldType:   SCALAR,
		Key:         attribute,
		ScalarValue: &fieldValue,
		Path:        path,
	}

	// Every resource has a computed ID in the root resource
	if path == "id" {
		field.Computed = true
	}

	p.currentResource().Append(field)
}
