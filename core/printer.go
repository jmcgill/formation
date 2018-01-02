package core

import (
	"bytes"
	//"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
)

const INDENT = 4

type Printer struct {
	output        *io.Writer
	currentIndent int
}

func (p *Printer) write(format string, a ...interface{}) {
	fmt.Fprintf(*p.output, "%s", strings.Repeat(" ", p.currentIndent))
	fmt.Fprintf(*p.output, format, a...)
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
	if field.Computed == true {
		return
	}

	if field.Link != "" {
		p.write("%s = \"${%s}\"\n", field.Key, field.Link)
	} else if field.FieldType == SCALAR {
		if field.ScalarValue.IsBool {
			if field.ScalarValue.StringValue == "true" {
				p.write("%s = true\n", field.Key)
			} else {
				p.write("%s = false\n", field.Key)
			}
		} else {
			if p.printJSON(field) {
				return
			}

			if field.ScalarValue.StringValue == "" {
				return
			}

			// TODO(jimmy): Make this anything except valid key characters
			if strings.Contains(field.Key, "/") {
				p.write("\"%s\" = \"%s\"\n", field.Key, strings.Replace(field.ScalarValue.StringValue, "\"", "\\\"", -1))
			} else {
				p.write("%s = \"%s\"\n", field.Key, strings.Replace(field.ScalarValue.StringValue, "\"", "\\\"", -1))
			}


		}
	} else if field.FieldType == MAP {
		p.printMap(field.Key, field.NestedValue)
	} else if field.FieldType == LIST {
		p.printList(field.Key, field.NestedValue)
	}
}

func (p *Printer) printJSON(field *Field) bool {
	if field.ScalarValue.StringValue == "" {
		return false
	}

	if field.ScalarValue.StringValue[0] != '{' {
		return false
	}

	// ${} references in Policy documents conflict with Terraform. Use &{} instead
	//var d map[string]interface{}
	//err := json.Unmarshal([]byte(field.ScalarValue.StringValue), &d)
	//if err != nil {
	//	fmt.Printf("This doesn't appear to be JSON\n")
	//	return false
	//}
	//
	//s, _ := json.MarshalIndent(d, "", "    ")

	fixed := strings.Replace(string(field.ScalarValue.StringValue), "${", "&{", -1)
	p.write("%s = <<EOF\n%sEOF\n", field.Key, fixed)
	return true
}

func (p *Printer) Print(resource *Resource) string {
	buf := bytes.Buffer{}
	writer := io.Writer(&buf)

	p.output = &writer
	p.write("resource \"%s\" \"%s\" {\n", resource.Type, resource.Name)
	p.printInlineResource(resource.Fields)
	p.write("}")

	return buf.String()
}

func (p *Printer) PrintToFile(file *os.File, resource *Resource) {
	writer := io.Writer(file)
	p.output = &writer

	p.write("resource \"%s\" \"%s\" {\n", resource.Type, resource.Name)
	p.printInlineResource(resource.Fields)
	p.write("}")
}
