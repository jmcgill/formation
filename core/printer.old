package core

import (
	"bytes"
	"fmt"
	"sort"
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

func (p *Printer) printInlineResource(fields InlineResource) {
	p.indent()

	for _, v := range fields {
		p.printField(&v)
	}

	p.unindent()
}

func (p *Printer) printMap(key string, field map[string]ScalarValue) {
	p.write("%s {\n", key)
	p.indent()

	// For consistency and to simplify testing, we guarantee that keys are emitted alphabetically.
	keys := make([]string, len(field))

	i := 0
	for k := range field {
		keys[i] = k
		i++
	}
	sort.Strings(keys)

	for _, v := range keys {
		r := field[v]
		p.write("%s = \"%s\"\n", v, r.StringValue)
	}

	p.unindent()
	p.write("}\n")
}

func (p *Printer) printSimpleList(key string, field []Field) {
	p.write("%s = [\n", key)
	p.indent()

	for _, v := range field {
		p.write("\"%s\",\n", v.ScalarValue.StringValue)
	}

	p.unindent()
	p.write("]\n")
}

func (p *Printer) printRichList(key string, field []Field) {
	for i, v := range field {
		p.write("%s {\n", key)
		p.printInlineResource(v.NestedValue)

		// Place an empty line between successive nested fields, because it looks better.
		if i == (len(field) - 1) {
			p.write("}\n")
		} else {
			p.write("}\n\n")
		}
	}
}

func (p *Printer) printList(key string, field []Field) {
	// Lists can either contain a set of scalar objects, or a set of nested resources.
	if field[0].FieldType == SCALAR {
		p.printSimpleList(key, field)
	} else {
		p.printRichList(key, field)
	}
}

func (p *Printer) printField(field *Field) {
	if field.FieldType == SCALAR {
		p.write("%s = \"%s\"\n", field.Key, field.ScalarValue.StringValue)
	} else if field.FieldType == MAP {
		p.printMap(field.Key, field.MapValue)
	} else if field.FieldType == LIST {
		p.printList(field.Key, field.ListValue)
	}
}

func (p *Printer) Print(resource *Resource) string {
	p.output.Reset()
	p.write("resource \"%s\" \"%s\" {\n", resource.Type, resource.Name)
	p.printInlineResource(resource.Fields)
	p.write("}")
	return p.output.String()
}
