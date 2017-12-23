package core

import (
	"fmt"
	"formation/terraform_helpers"
	"strconv"
	"strings"
)

type ParentType int

const (
	PARENT_ROOT ParentType = iota
	PARENT_MAP
	PARENT_LIST
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

func (p *InstanceStateParser) Parse(state *terraform_helpers.InstanceState) *Resource {
	resource := new(Resource)
	sortedState := state.ToSorted()

	s := State{
		remainingChildren: 0,
		parent:            new(InlineResource),
		parentType:        PARENT_ROOT,
	}
	p.pushState(&s)
	resource.Fields = s.parent

	for _, a := range sortedState.Attributes {
		fmt.Printf("Attribute is %s = %s\n", a.Key, a.Value)
		p.parseAttribute(a.Key, a.Value)
	}

	// TEMP
	// _ = p.popState()

	// TODO(jimmy): Add an accessor function for current resource
	// resource.Fields = p.state().parent
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
	fmt.Printf("***** Attribute %s - depth %i\n\n\n", attribute, p.state().depth)
	// We cannot establish the prefix for a list item until we first encounter it.
	if p.state().parentType == PARENT_LIST && p.state().prefix == "" {
		fmt.Print("We are in a list with an unestablished prefix!!\n")
		// Determine whether this is a list of scalar objects or a list of nested objects.
		parts := strings.Split(attribute, ".")

		// For scalar lists this is just a unique key for this entry. For nested objects
		// the first entry contains both that unique key _and_ the first nested field.
		fmt.Printf("Depth is %i\n", p.state().depth)
		fmt.Printf("Split parts is %s\n", parts)
		listPrefix := strings.Join(parts[p.state().depth-1:], ".")

		fmt.Printf("The list prefix is %s\n", listPrefix)

		if strings.ContainsRune(listPrefix, '.') {
			fmt.Printf("This is a nested resource, because it contains a dot\n")
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
				parentType:        PARENT_MAP,
				depth:             p.state().depth + 0,
				prefix:            strings.Join(parts[:p.state().depth], ".") + ".",
			}
			p.pushState(s)
		} else {
			fmt.Printf("This is not a nested resource\n")
			// This is a list of scalar objects
			p.state().prefix = strings.Join(parts[:p.state().depth], ".")
		}
	}

	fmt.Printf("Comparing %s to %s\n", attribute, p.state().prefix)
	if !strings.HasPrefix(attribute, p.state().prefix) {
		fmt.Printf("End of prefix matching, adding a new stack item\n")

		// We should also pop if we've reached the end of key
		// TODO(jimmy): Is this always safe?
		if p.stateStack[len(p.stateStack)-2].parentType == PARENT_LIST {
			_ = p.popState()
		}

		// This is a sign that we've either reached the end of a map, or the end of a single
		// item in a list.
		fmt.Printf("Remaining children: %i\n", p.state().remainingChildren)
		p.state().remainingChildren -= 1
		fmt.Printf("Remaining children: %i\n", p.state().remainingChildren)
		if p.state().remainingChildren == 0 {
			fmt.Printf("End of entries to parse\n")
			// Pop an entry off the stack
			_ = p.popState()
		}

		// If this is a list, we need to reset the matched prefix.
		// TODO(jimmy): Factor this out into a single method
		if p.state().parentType == PARENT_LIST {
			fmt.Printf("We are still in a list, let's reset the state\n")
			// This is a bad thing to do, as it clears any previous nesting
			// but it's OK because the next field we see will restore it.
			// Determine whether this is a list of scalar objects or a list of nested objects.
			parts := strings.Split(attribute, ".")

			// For scalar lists this is just a unique key for this entry. For nested objects
			// the first entry contains both that unique key _and_ the first nested field.
			listPrefix := strings.Join(parts[p.state().depth-1:], ".")

			fmt.Printf("The new list prefix is %s\n", listPrefix)

			if strings.ContainsRune(listPrefix, '.') {
				// This is a list of nested resource
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
					parentType:        PARENT_MAP,
					depth:             p.state().depth + 0,
					prefix:            strings.Join(parts[:p.state().depth], ".") + ".",
				}
				p.pushState(s)
			} else {
				// This is a list of scalar objects
				p.state().prefix = strings.Join(parts[:p.state().depth], ".")
			}
			//attribute = strings.TrimPrefix(attribute, p.state().prefix)
		}
	}
	//else {
	attribute = strings.TrimPrefix(attribute, p.state().prefix)
	//}

	fmt.Printf("The current prefix is: %s\n", p.state().prefix)
	fmt.Printf("Looking for dot rune in attribute %s\n", attribute)

	if !strings.ContainsRune(attribute, '.') {
		fmt.Printf("Parsing a simple attribute\n")
		p.parseSimpleAttribute(attribute, value)
	} else {
		parts := strings.Split(attribute, ".")
		fieldName, parts := parts[0], parts[1:]

		// Is this the beginning of a map?
		if parts[0] == "%" {
			children, _ := strconv.Atoi(value)

			// Is this the map declaration?
			field := &Field{
				FieldType:   MAP,
				Key:         fieldName,
				NestedValue: new(InlineResource),
			}
			p.currentResource().Append(field)

			// Push this entry onto the stack.
			s := &State{
				parent:            field.NestedValue,
				remainingChildren: children + 1,
				parentType:        PARENT_MAP,
				depth:             p.state().depth + 1,
				prefix:            p.state().prefix + fieldName + ".",
			}
			p.pushState(s)
		}

		// Is this the beginning of a list?
		if parts[0] == "#" {
			children, _ := strconv.Atoi(value)
			if children != 0 {
				field := &Field{
					FieldType:   LIST,
					Key:         fieldName,
					NestedValue: new(InlineResource),
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

				fmt.Println("Appending a LIST resource")
			}
		}
	}
}

func (p *InstanceStateParser) parseSimpleAttribute(attribute string, value string) {
	fieldValue := ScalarValue{StringValue: value}
	field := &Field{
		FieldType:   SCALAR,
		Key:         attribute,
		ScalarValue: &fieldValue,
	}
	p.currentResource().Append(field)
}
