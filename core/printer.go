package core

import (
	"bytes"
	"fmt"
	"strings"
)

const INDENT = 4

type Printer struct {
	output        bytes.Buffer
	currentIndent int
}

func (p *Printer) write(format string, a ...interface{}) {
	p.output.WriteString(strings.Repeat(" ", p.currentIndent))
	p.output.WriteString(fmt.Sprintf(format, a...))
}

func (p *Printer) indent() {
	p.currentIndent += INDENT
}

func (p *Printer) unindent() {
	p.currentIndent -= INDENT
}

func (p *Printer) printInlineResource(resource *InlineResource) {
	if resource == nil {
		return
	}

	p.indent()

	for _, v := range resource.Fields {
		p.printField(v)
	}

	p.unindent()
}

func (p *Printer) printMap(key string, resource *InlineResource) {
	p.write("%s {\n", key)
	p.indent()

	for _, field := range resource.Fields {
		p.printField(field)
	}

	p.unindent()
	p.write("}\n")
}

func (p *Printer) printSimpleList(key string, resource *InlineResource) {
	p.write("%s = [\n", key)
	p.indent()

	for _, v := range resource.Fields {
		p.write("\"%s\",\n", v.ScalarValue.StringValue)
	}

	p.unindent()
	p.write("]\n")
}

func (p *Printer) printRichList(key string, resource *InlineResource) {
	for i, v := range resource.Fields {
		p.write("%s {\n", key)
		p.printInlineResource(v.NestedValue)

		// Place an empty line between successive nested fields, because it looks better.
		if i == (len(resource.Fields) - 1) {
			p.write("}\n")
		} else {
			p.write("}\n\n")
		}
	}
}

func (p *Printer) printList(key string, resource *InlineResource) {
	// Lists can either contain a set of scalar objects, or a set of nested resources.
	if resource.Fields[0].FieldType == SCALAR {
		p.printSimpleList(key, resource)
	} else {
		p.printRichList(key, resource)
	}
}

func (p *Printer) printField(field *Field) {
	if field.FieldType == SCALAR {
		p.write("%s = \"%s\"\n", field.Key, field.ScalarValue.StringValue)
	} else if field.FieldType == MAP {
		p.printMap(field.Key, field.NestedValue)
	} else if field.FieldType == LIST {
		p.printList(field.Key, field.NestedValue)
	}
}

func (p *Printer) Print(resource *Resource) string {
	p.output.Reset()
	p.write("resource \"%s\" \"%s\" {\n", resource.Type, resource.Name)
	p.printInlineResource(resource.Fields)
	p.write("}")
	return p.output.String()
}
