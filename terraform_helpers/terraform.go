package terraform_helpers

import (
	"bytes"
	"fmt"
	"sort"

	"github.com/hashicorp/terraform/terraform"
	"strconv"
	"strings"
)

type InstanceState terraform.InstanceState

//type InstanceState struct {
//	// Attributes are basic information about the resource. Any keys here
//	// are accessible in variable format within Terraform configurations:
//	// ${resourcetype.name.attribute}.
//	Attributes map[string]string
//}

type KeyValue struct {
	Key   string
	Value string
}

func isNumeric(in string) bool {
	if _, err := strconv.Atoi(in); err == nil {
		return true
	}
	return false
}

func asInteger(in string) int {
	v, _ := strconv.Atoi(in)
	return v
}

type KeyValueList []KeyValue

func (s KeyValueList) Len() int {
	return len(s)
}

func (s KeyValueList) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s KeyValueList) Less(i, j int) bool {
	ip := strings.Split(s[i].Key, ".")
	jp := strings.Split(s[j].Key, ".")
	index := 0

	for jp[index] == ip[index] {
		index += 1
	}

	if isNumeric(jp[index]) && isNumeric(ip[index]) {
		return asInteger(ip[index]) < asInteger(jp[index])
	}
	return ip[index] < jp[index]
}

type SortedInstanceState struct {
	Attributes KeyValueList
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

	//attrKeys := make([]string, 0, len(attributes))

	for k, v := range attributes {
		kv := KeyValue{
			Key:   k,
			Value: v,
		}
		r.Attributes = append(r.Attributes, kv)
	}
	sort.Sort(r.Attributes)
	return r
}
