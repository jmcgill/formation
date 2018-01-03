package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/jmcgill/formation/core"
	"strconv"
)

type tagNamer struct {
	existingNames map[string]int
}

func NewTagNamer() *tagNamer {
	n := &tagNamer{}
	n.existingNames = make(map[string]int)
	return n
}

func (n *tagNamer) NameOrDefault(tags []*ec2.Tag, otherwise *string) string {
	name := aws.StringValue(otherwise)
	for _, t := range tags {
		if aws.StringValue(t.Key) == "Name" {
			name = aws.StringValue(t.Value)
		}
	}

	if _, ok := n.existingNames[name]; ok {
		n.existingNames[name] += 1
		name = name + "-" + strconv.Itoa(n.existingNames[name])
	} else {
		n.existingNames[name] = 1
	}

	return core.Format(name)
}
