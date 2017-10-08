package terraform

import (
	"bytes"
	"fmt"
	"sort"
)

type InstanceState struct {
	// Attributes are basic information about the resource. Any keys here
	// are accessible in variable format within Terraform configurations:
	// ${resourcetype.name.attribute}.
	Attributes map[string]string
}

type KeyValue struct {
	Key   string
	Value string
}

type SortedInstanceState struct {
	Attributes []KeyValue
}

func (s *InstanceState) String() string {
	var buf bytes.Buffer

	attributes := s.Attributes
	attrKeys := make([]string, 0, len(attributes))
	for ak, _ := range attributes {
		if ak == "id" {
			continue
		}

		attrKeys = append(attrKeys, ak)
	}
	sort.Strings(attrKeys)

	for _, ak := range attrKeys {
		av := attributes[ak]
		buf.WriteString(fmt.Sprintf("%s = %s\n", ak, av))
	}

	return buf.String()
}

func (s *InstanceState) ToSorted() *SortedInstanceState {
	r := new(SortedInstanceState)

	attributes := s.Attributes
	attrKeys := make([]string, 0, len(attributes))
	for ak, _ := range attributes {
		if ak == "id" {
			continue
		}

		attrKeys = append(attrKeys, ak)
	}
	sort.Strings(attrKeys)

	for _, ak := range attrKeys {
		kv := KeyValue{
			Key:   ak,
			Value: attributes[ak],
		}
		r.Attributes = append(r.Attributes, kv)
	}

	return r
}
