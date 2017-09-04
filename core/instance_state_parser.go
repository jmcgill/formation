package core

import (
	"fmt"
	"formation/terraform"
	"strings"
)

type InstanceStateParser struct {
	currentResource *Resource
}

func (p *InstanceStateParser) Parse(state *terraform.InstanceState) {
	p.currentResource = new(Resource)
	p.currentResource.Fields = make([]Field, 10)

	for k, v := range state.Attributes {
		p.ParseSimpleAttribute(k, v)
	}

	fmt.Println("Hello, world")
}

func (p *InstanceStateParser) ParseAttribute(attribute string, value string) {
	if !strings.ContainsRune(attribute, '.') {
		p.ParseSimpleAttribute(attribute, value)
	}
}

func (p *InstanceStateParser) ParseSimpleAttribute(attribute string, value string) {
	fieldValue := ScalarValue{StringValue: value}
	field := Field{
		FieldType:   SCALAR,
		Key:         attribute,
		ScalarValue: fieldValue,
	}
	p.currentResource.Fields = append(p.currentResource.Fields, field)
}